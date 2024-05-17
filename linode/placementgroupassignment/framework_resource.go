package placementgroupassignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_placement_group_assignment",
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
	tflog.Debug(ctx, "Create linode_placement_group_assignment")

	var plan PGAssignmentModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pgID := helper.FrameworkSafeInt64ToInt(
		plan.PlacementGroupID.ValueInt64(),
		&resp.Diagnostics,
	)

	linodeID := helper.FrameworkSafeInt64ToInt(
		plan.LinodeID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	assignOpts := linodego.PlacementGroupAssignOptions{
		Linodes:       []int{linodeID},
		CompliantOnly: plan.CompliantOnly.ValueBoolPointer(),
	}

	tflog.Debug(ctx, "client.AssignPlacementGroupLinodes(...)", map[string]any{
		"placement_group_id": pgID,
		"options":            assignOpts,
	})
	pg, err := client.AssignPlacementGroupLinodes(ctx, pgID, assignOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error assigning Linode to Placement Group",
			err.Error(),
		)
		return
	}

	plan.Flatten(*pg, linodeID, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(
		buildID(
			idFormat{
				PGID:     pgID,
				LinodeID: linodeID,
			},
			&resp.Diagnostics,
		),
	)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_placement_group_assignment")

	var state PGAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	idData := parseID(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"pg_id":     idData.PGID,
		"linode_id": idData.LinodeID,
	})

	client := r.Meta.Client

	pg, err := client.GetPlacementGroup(ctx, idData.PGID)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Placement Group No Longer Exists",
				fmt.Sprintf(
					"Removing Placement Group assignment %s from state because the "+
						"target Placement Group no longer exists",
					state.ID.String(),
				),
			)
			resp.State.RemoveResource(ctx)
		}

		resp.Diagnostics.AddError(
			"Failed to get Placement Group",
			err.Error(),
		)
		return
	}

	state.Flatten(*pg, idData.LinodeID, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_placement_group_assignment")
	resp.Diagnostics.AddWarning(
		"Unintended Calling to Update Function",
		"The Update function of 'linode_placement_group_assignment' should never be "+
			"invoked by design. This function has been redundantly implemented "+
			"for improved reliability. Please consider reporting this as a bug "+
			"to the provider developers.",
	)

	var state, plan PGAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_placement_group_assignment")
	var state PGAssignmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	client := r.Meta.Client

	idData := parseID(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UnassignPlacementGroupLinodes(...)", map[string]any{
		"placement_group_id": idData.PGID,
		"linode_id":          idData.LinodeID,
	})

	_, err := client.UnassignPlacementGroupLinodes(ctx, idData.PGID, linodego.PlacementGroupUnAssignOptions{
		Linodes: []int{idData.LinodeID},
	})
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf(
					"Attempted to delete PG assignment to Linode %d but the resource was not found",
					idData.LinodeID,
				),
				err.Error(),
			)
		} else if linodego.ErrHasStatus(err, 400) {
			// Nothing to do here, we can assume the Linode isn't attached to the PG
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error deleting PG assignment to Linode %d", idData.LinodeID),
				err.Error(),
			)
		}
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import linode_placement_group_assignment")

	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "placement_group_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "linode_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
		},
	)

	// Manually populate the ID attribute
	var state PGAssignmentModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(
		buildID(
			idFormat{
				PGID:     helper.FrameworkSafeInt64ToInt(state.PlacementGroupID.ValueInt64(), &resp.Diagnostics),
				LinodeID: helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics),
			},
			&resp.Diagnostics,
		),
	)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
