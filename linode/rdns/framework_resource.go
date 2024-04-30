package rdns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	DefaultVolumeCreateTimeout = 15 * time.Minute
	DefaultVolumeUpdateTimeout = 15 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_rdns",
				Schema: &frameworkResourceSchema,
				IDType: types.StringType,
				TimeoutOpts: &timeouts.Opts{
					Update: true,
					Create: true,
				},
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create linode_rdns")

	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, plan)

	createTimeout, diags := plan.Timeouts.Create(ctx, DefaultVolumeCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	address := plan.Address.ValueString()

	tflog.Trace(ctx, "client.GetIPAddress(...)", map[string]any{
		"address": address,
	})

	ip, err := client.GetIPAddress(ctx, address)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get the ip address associated with this RDNS",
			err.Error(),
		)
		return
	}

	defaultRdns := strings.Replace(
		plan.Address.ValueString(),
		".",
		"-",
		-1,
	) + ".ip.linodeusercontent.com"

	if ip.RDNS != defaultRdns {
		resp.Diagnostics.AddWarning(
			"Pre-modified RDNS Address",
			"RDNS was already configured before the creation of this RDNS resource",
		)
	}

	ip, err = updateIPAddress(
		ctx,
		client,
		plan.Address.ValueString(),
		plan.RDNS.ValueStringPointer(),
		plan.WaitForAvailable.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Linode RDNS",
			err.Error(),
		)
		return
	}

	plan.FlattenInstanceIP(ip, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_rdns")

	client := r.Meta.Client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	tflog.Trace(ctx, "client.GetIPAddress(...)")
	ip, err := client.GetIPAddress(ctx, data.ID.ValueString())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"RDNS No Longer Exists",
				fmt.Sprintf(
					"Removing Linode RDNS with IP %v from state because it no longer exists",
					data.ID.ValueString(),
				),
			)
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Failed to read the Linode RDNS", err.Error(),
			)
		}
		return
	}

	data.FlattenInstanceIP(ip, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_rdns")

	var state, plan ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	updateTimeout, diags := plan.Timeouts.Update(ctx, DefaultVolumeUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	client := r.Meta.Client

	var updateOpts linodego.IPAddressUpdateOptions

	resourceUpdated := false

	if !state.RDNS.Equal(plan.RDNS) {
		updateOpts.RDNS = plan.RDNS.ValueStringPointer()
		resourceUpdated = true
	}

	if resourceUpdated {
		ip, err := updateIPAddress(
			ctx,
			client,
			plan.Address.ValueString(),
			plan.RDNS.ValueStringPointer(),
			plan.WaitForAvailable.ValueBool(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to update the Linode RDNS",
				err.Error(),
			)
			return
		}
		plan.FlattenInstanceIP(ip, true)
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_rdns")

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: nil,
	}

	tflog.Debug(ctx, "client.UpdateIPAddress(...)", map[string]any{
		"options": updateOpts,
	})
	_, err := client.UpdateIPAddress(ctx, data.Address.ValueString(), updateOpts)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Target IP for RDNS resetting no longer exists.",
				fmt.Sprintf(
					"The given IP Address (%s) for RDNS resetting no longer exists.",
					data.Address,
				),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete the Linode IP address RDNS",
			fmt.Sprintf(
				"Error deleting the Linode IP address RDNS: %s",
				err.Error(),
			),
		)
	}
}

func populateLogAttributes(ctx context.Context, model ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"address": model.Address.ValueString(),
	})
}
