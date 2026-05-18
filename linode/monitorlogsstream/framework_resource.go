package monitorlogsstream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
				Name:   "linode_monitor_logs_stream",
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

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := data.GetCreateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.CreateLogStream(...)", map[string]any{
		"options": createOpts,
	})
	stream, err := client.CreateLogStream(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create logs stream.",
			err.Error(),
		)
		return
	}

	data.FlattenStream(ctx, stream, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	data.ID = types.StringValue(strconv.Itoa(stream.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	stream, err := client.GetLogStream(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Logs stream no longer exists.",
				fmt.Sprintf(
					"Removing logs stream with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the logs stream.",
			err.Error(),
		)
		return
	}

	data.FlattenStream(ctx, stream, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updateOpts := plan.GetUpdateOptions(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateLogStream(...)", map[string]any{
		"options": updateOpts,
	})
	stream, err := client.UpdateLogStream(ctx, id, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update logs stream.",
			err.Error(),
		)
		return
	}

	plan.FlattenStream(ctx, stream, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteLogStream(...)", map[string]any{
		"id": id,
	})
	if err := client.DeleteLogStream(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete logs stream.",
			err.Error(),
		)
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
