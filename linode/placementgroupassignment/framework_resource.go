package placementgroupassignment

import (
	"context"
	"fmt"
	"strings"

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

	pgID, linodeID := plan.GetIDComponents(&resp.Diagnostics)
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
			pgID,
			linodeID,
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

	pgID, linodeID := state.GetIDComponents(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"placement_group_id": pgID,
		"linode_id":          linodeID,
	})

	client := r.Meta.Client

	pg, err := client.GetPlacementGroup(ctx, pgID)
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
			return
		}

		resp.Diagnostics.AddError(
			"Failed to get Placement Group",
			err.Error(),
		)
		return
	}

	// If the Linode doesn't exist in the PG due to external modification,
	// we should mark this assignment for recreation
	if !pgHasID(*pg, linodeID) {
		resp.Diagnostics.AddWarning(
			"Marking Assignment for Recreation",
			fmt.Sprintf(
				"Linode (%d) is no longer assigned to target Placement Group (%d).",
				pg.ID,
				linodeID,
			),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Flatten(*pg, linodeID, false, &resp.Diagnostics)
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

	pgID, linodeID := state.GetIDComponents(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UnassignPlacementGroupLinodes(...)", map[string]any{
		"placement_group_id": pgID,
		"linode_id":          linodeID,
	})

	_, err := client.UnassignPlacementGroupLinodes(ctx, pgID, linodego.PlacementGroupUnAssignOptions{
		Linodes: []int{linodeID},
	})
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf(
					"Attempted to delete PG assignment to Linode %d but the resource was not found",
					linodeID,
				),
				err.Error(),
			)
		} else if lErr, ok := err.(*linodego.Error); ok &&
			lErr.Code == 400 &&
			strings.Contains(lErr.Message, "does not belong to Placement Group") {
			// Nothing to do here, we can assume the Linode isn't attached to the PG
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error deleting PG assignment to Linode %d", linodeID),
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

	// We need to manually set the ID in state
	// because it is not implicitly populated by one of the
	// ID attributes above
	var state PGAssignmentModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pgID, linodeID := state.GetIDComponents(&resp.Diagnostics)
	state.ID = types.StringValue(buildID(pgID, linodeID, &resp.Diagnostics))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func pgHasID(pg linodego.PlacementGroup, id int) bool {
	for _, member := range pg.Members {
		if member.LinodeID == id {
			return true
		}
	}

	return false
}
