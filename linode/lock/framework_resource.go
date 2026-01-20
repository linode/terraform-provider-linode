package lock

import (
	"context"
	"fmt"
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
				Name:          "linode_lock",
				IDType:        types.StringType,
				Schema:        &frameworkResourceSchema,
				IsEarlyAccess: true,
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

	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := plan.GetCreateOptions(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateLock(...)", map[string]any{
		"options": createOpts,
	})

	lock, err := client.CreateLock(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Lock",
			err.Error(),
		)
		return
	}

	plan.FlattenLock(lock, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(lock.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: cleanup when Crossplane fixes it
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	id := helper.StringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "lock_id", id)

	lock, err := client.GetLock(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Lock No Longer Exists",
				fmt.Sprintf(
					"Removing Lock with ID %v from state because it no longer exists",
					state.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Refresh the Lock",
			fmt.Sprintf(
				"Error finding the specified Lock: %s",
				err.Error(),
			),
		)
		return
	}

	state.FlattenLock(lock, false)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	// Locks do not support updates - all fields require replacement
	resp.Diagnostics.AddWarning(
		"Unintended Calling to Update Function",
		"The Update function of 'linode_lock' should never be "+
			"invoked by design. All fields require replacement. "+
			"This function has been implemented for reliability. "+
			"Please consider reporting this as a bug to the provider developers.",
	)

	var state, plan ResourceModel

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
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "lock_id", id)

	tflog.Debug(ctx, "client.DeleteLock(...)", map[string]any{
		"lock_id": id,
	})

	err := client.DeleteLock(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Lock Already Deleted",
				fmt.Sprintf(
					"Attempted to delete Lock %d but it was not found",
					id,
				),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting Lock %d", id),
			err.Error(),
		)
	}
}
