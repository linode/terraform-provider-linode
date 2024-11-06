package instancereservedipassignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// For internal only: please refer to KB page for more information

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_reserved_ip_assignment",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
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
	tflog.Debug(ctx, "Create linode_reserved_ip_assignment")
	var plan InstanceIPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, &plan)

	linodeID := helper.FrameworkSafeInt64ToInt(
		plan.LinodeID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	isPublic := plan.Public.ValueBool()

	client := r.Meta.Client
	var ip *linodego.InstanceIP
	var err error

	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		// Assign a reserved IP
		createOpts := linodego.InstanceReserveIPOptions{
			Type:    "ipv4",
			Public:  isPublic,
			Address: plan.Address.ValueString(),
		}
		ip, err = client.AssignInstanceReservedIP(ctx, linodeID, createOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to assign reserved IP to instance (%d)", linodeID),
				err.Error(),
			)
			return
		}
	}

	if !plan.RDNS.IsNull() && !plan.RDNS.IsUnknown() {
		rdns := plan.RDNS.ValueString()

		options := linodego.IPAddressUpdateOptions{
			RDNS: &rdns,
		}

		tflog.Debug(ctx, "client.UpdateIPAddress(...)", map[string]any{
			"options": options,
		})

		if _, err := client.UpdateIPAddress(ctx, ip.Address, options); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf(
					"failed to set RDNS for instance (%d) ip (%s)",
					linodeID, ip.Address,
				),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(plan.flattenInstanceIP(ctx, *ip, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(ip.Address)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ApplyImmediately.ValueBool() {
		tflog.Debug(ctx, "Attempting apply_immediately")

		instance, err := client.GetInstance(ctx, linodeID)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Get Linode Instance", err.Error())
			return
		}

		if instance.Status == linodego.InstanceRunning {
			tflog.Info(ctx, "detected instance in running status, rebooting instance")
			ctx, cancel := context.WithTimeout(ctx, time.Duration(600)*time.Second)
			resp.Diagnostics.Append(helper.FrameworkRebootInstance(ctx, linodeID, client, 0)...)
			cancel()
		} else {
			tflog.Info(ctx, "Detected instance not in running status, can't perform a reboot.")
		}
	}
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_instance_reserved_ip")
	var state InstanceIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: cleanup when Crossplane fixes it
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	ctx = populateLogAttributes(ctx, &state)

	client := r.Meta.Client
	address := state.ID.ValueString() // Use ID as address

	ip, err := client.GetIPAddress(ctx, address)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Reserved IP No Longer Exists",
				fmt.Sprintf("Removing reserved IP %s from state because it no longer exists", state.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Refresh Reserved IP",
			fmt.Sprintf("Error finding the specified Reserved IP: %s", err.Error()),
		)
		return
	}

	if ip == nil {
		resp.Diagnostics.AddError("nil Pointer", "received nil pointer for reserved IP")
		return
	}

	// Populate state with fetched data
	resp.Diagnostics.Append(state.flattenInstanceIP(ctx, *ip, false)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan InstanceIPModel
	var state InstanceIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.RDNS.Equal(state.RDNS) {
		rdns := plan.RDNS.ValueStringPointer()
		updateOptions := linodego.IPAddressUpdateOptions{
			RDNS: rdns,
		}

		client := r.Meta.Client
		address := plan.Address.ValueString()
		linodeID := plan.LinodeID.ValueInt64()

		tflog.Debug(ctx, "client.UpdateIPAddress(...)", map[string]any{
			"options": updateOptions,
		})

		ip, err := client.UpdateIPAddress(ctx, address, updateOptions)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Update RDNS",
				fmt.Sprintf(
					"failed to update RDNS for instance (%d) ip (%s): %s",
					linodeID, address, err,
				),
			)
			return
		}
		if ip == nil {
			resp.Diagnostics.AddError(
				"Failed to Get Updated IP",
				fmt.Sprintf(
					"ip is a nil pointer after update operation for instance (%d) ip (%s): %s",
					linodeID, address, err,
				),
			)
			return
		}
		resp.Diagnostics.Append(plan.flattenInstanceIP(ctx, *ip, true)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	plan.CopyFrom(ctx, state, true)

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
	var state InstanceIPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	address := state.Address.ValueString()
	linodeID := helper.FrameworkSafeInt64ToInt(
		state.LinodeID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteInstanceIPAddress(...)")
	if err := client.DeleteInstanceIPAddress(ctx, linodeID, address); err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				"Failed to Delete IP",
				fmt.Sprintf(
					"failed to delete instance (%d) ip (%s): %s",
					linodeID, address, err.Error(),
				),
			)
		}
	}
}

func populateLogAttributes(ctx context.Context, data *InstanceIPModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": data.LinodeID.ValueInt64(),
		"address":   data.ID.ValueString(),
	})
}
