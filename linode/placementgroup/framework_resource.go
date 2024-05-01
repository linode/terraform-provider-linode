package placementgroup

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_placement_group",
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
	tflog.Debug(ctx, "Create linode_placement_group")

	var data PlacementGroupModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.PlacementGroupCreateOptions{
		Label:        data.Label.ValueString(),
		Region:       data.Region.ValueString(),
		AffinityType: linodego.PlacementGroupAffinityType(data.AffinityType.ValueString()),
		IsStrict:     data.IsStrict.ValueBool(),
	}

	tflog.Debug(ctx, "client.CreatePlacementGroup(...)", map[string]any{
		"options": createOpts,
	})
	pg, err := client.CreatePlacementGroup(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Placement Group.",
			err.Error(),
		)
		return
	}

	data.FlattenPlacementGroup(pg, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(strconv.Itoa(pg.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_placement_group")

	var data PlacementGroupModel
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

	pg, err := client.GetPlacementGroup(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Placement Group no longer exists.",
				fmt.Sprintf(
					"Removing Linode Placement Group with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Placement Group.",
			err.Error(),
		)
		return
	}

	data.FlattenPlacementGroup(pg, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_placement_group")

	client := r.Meta.Client
	var plan, state PlacementGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	var updateOpts linodego.PlacementGroupUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueString()
	}

	if shouldUpdate {
		id := helper.FrameworkSafeStringToInt(plan.ID.ValueString(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "client.UpdatePlacementGroup(...)", map[string]any{
			"options": updateOpts,
		})
		pg, err := client.UpdatePlacementGroup(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Placement Group (%d).", id),
				err.Error(),
			)
			return
		}
		plan.FlattenPlacementGroup(pg, false)
	}
	plan.CopyFrom(state, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_placement_group")

	client := r.Meta.Client
	var data PlacementGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeletePlacementGroup(...)")
	err := client.DeletePlacementGroup(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete the Placement Group (%s)", data.ID.ValueString()),
			err.Error(),
		)
	}
}

func populateLogAttributes(ctx context.Context, data PlacementGroupModel) context.Context {
	return tflog.SetField(ctx, "id", data.ID)
}
