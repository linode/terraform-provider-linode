package vpc

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_vpc",
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
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var data Model
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: support creating subnets with VPC creation after upgrade to protocol version 6
	vpcCreateOpts := linodego.VPCCreateOptions{
		Label:       data.Label.ValueString(),
		Region:      data.Region.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.IPv6.IsNull() {
		modelIPv6s := make([]ResourceModelIPv6, len(data.IPv6.Elements()))

		resp.Diagnostics.Append(data.IPv6.ElementsAs(ctx, &modelIPv6s, true)...)
		if resp.Diagnostics.HasError() {
			return
		}

		vpcCreateOpts.IPv6 = slices.Collect(
			helper.Map(
				slices.Values(modelIPv6s),
				func(m ResourceModelIPv6) linodego.VPCCreateOptionsIPv6 {
					return linodego.VPCCreateOptionsIPv6{
						Range:           m.Range.ValueStringPointer(),
						AllocationClass: m.AllocationClass.ValueStringPointer(),
					}
				},
			),
		)
	}

	tflog.Debug(ctx, "client.CreateVPC(...)", map[string]any{
		"options": vpcCreateOpts,
	})
	vpc, err := client.CreateVPC(ctx, vpcCreateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create VPC.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenVPC(ctx, vpc, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(strconv.Itoa(vpc.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data Model
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	vpc, err := client.GetVPC(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"VPC no longer exists.",
				fmt.Sprintf(
					"Removing Linode VPC with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the VPC.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenVPC(ctx, vpc, false)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan, state Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	var updateOpts linodego.VPCUpdateOptions
	shouldUpdate := false

	if !state.Description.Equal(plan.Description) {
		shouldUpdate = true
		updateOpts.Description = plan.Description.ValueString()
	}

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if shouldUpdate {
		id := helper.FrameworkSafeStringToInt(plan.ID.ValueString(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "client.UpdateVPC(...)", map[string]any{
			"options": updateOpts,
		})
		vpc, err := client.UpdateVPC(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update VPC (%d).", id),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(plan.FlattenVPC(ctx, vpc, false)...)
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
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	client := r.Meta.Client
	var data Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteVPC(...)")
	err := client.DeleteVPC(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the VPC (%s)", data.ID.ValueString()),
				err.Error(),
			)
		}
		return
	}
}

func populateLogAttributes(ctx context.Context, data Model) context.Context {
	return tflog.SetField(ctx, "id", data.ID)
}
