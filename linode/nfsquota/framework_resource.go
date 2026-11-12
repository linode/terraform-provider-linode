package nfsquota

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_nfs_quota",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "space_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "filesystem_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "id",
				TypeConverter: helper.IDTypeConverterString,
			},
		})
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

	ctx = resourcePopulateLogAttributes(ctx, data)

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	filesystemID := helper.FrameworkSafeInt64ToInt(data.FilesystemID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := data.ToCreateOptions()

	tflog.Debug(ctx, "client.CreateNFSQuota(...)", map[string]any{
		"options": createOpts,
	})
	quota, err := client.CreateNFSQuota(ctx, spaceID, filesystemID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create NFS Quota.",
			err.Error(),
		)
		return
	}

	quota, err = client.WaitForNFSQuotaStatus(
		ctx,
		spaceID,
		filesystemID,
		quota.ID,
		linodego.NFSQuotaStatusActive,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to wait for NFS Quota to become active.",
			err.Error(),
		)
		return
	}

	data.FlattenQuota(quota, true)

	// IDs should always be overridden during creation (see #1085)
	data.ID = types.StringValue(strconv.Itoa(quota.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = resourcePopulateLogAttributes(ctx, data)
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	filesystemID := helper.FrameworkSafeInt64ToInt(data.FilesystemID.ValueInt64(), &resp.Diagnostics)
	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	quota, err := client.GetNFSQuota(ctx, spaceID, filesystemID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NFS Quota does not exist.",
				fmt.Sprintf(
					"Removing NFS Quota with ID %v from state because it no longer exists",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the NFS Quota.",
			err.Error(),
		)
		return
	}

	data.FlattenQuota(quota, false)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = resourcePopulateLogAttributes(ctx, state)

	var updateOpts linodego.NFSQuotaUpdateOptions
	shouldUpdate := false

	if !state.MaxCapacityBytes.Equal(plan.MaxCapacityBytes) {
		updateOpts.MaxCapacityBytes = plan.MaxCapacityBytes.ValueInt64Pointer()
		shouldUpdate = true
	}

	if !state.MaxFileCount.Equal(plan.MaxFileCount) {
		updateOpts.MaxFileCount = plan.MaxFileCount.ValueInt64Pointer()
		shouldUpdate = true
	}

	if !userGroupConfigEqual(state.UserGroupConfig, plan.UserGroupConfig) {
		if plan.UserGroupConfig != nil {
			updateOpts.UserGroupConfig = plan.UserGroupConfig.toOptions()
		}
		shouldUpdate = true
	}

	if shouldUpdate {
		spaceID := helper.FrameworkSafeInt64ToInt(plan.SpaceID.ValueInt64(), &resp.Diagnostics)
		filesystemID := helper.FrameworkSafeInt64ToInt(plan.FilesystemID.ValueInt64(), &resp.Diagnostics)
		id := helper.FrameworkSafeStringToInt(plan.ID.ValueString(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "client.UpdateNFSQuota(...)", map[string]any{
			"options": updateOpts,
		})
		quota, err := client.UpdateNFSQuota(ctx, spaceID, filesystemID, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update NFS Quota (%d).", id),
				err.Error(),
			)
			return
		}

		quota, err = client.WaitForNFSQuotaStatus(
			ctx,
			spaceID,
			filesystemID,
			quota.ID,
			linodego.NFSQuotaStatusActive,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to wait for NFS Quota to become active.",
				err.Error(),
			)
			return
		}

		plan.FlattenQuota(quota, true)
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

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = resourcePopulateLogAttributes(ctx, data)

	// The root quota cannot be deleted; it is removed along with its filesystem.
	if data.IsRoot.ValueBool() {
		tflog.Debug(ctx, "Skipping deletion of root NFS Quota; removing from state only")
		return
	}

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	filesystemID := helper.FrameworkSafeInt64ToInt(data.FilesystemID.ValueInt64(), &resp.Diagnostics)
	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteNFSQuota(...)")
	err := client.DeleteNFSQuota(ctx, spaceID, filesystemID, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Failed to delete the NFS Quota.",
			err.Error(),
		)
		return
	}
}

func resourcePopulateLogAttributes(ctx context.Context, data ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"space_id":      data.SpaceID.ValueInt64(),
		"filesystem_id": data.FilesystemID.ValueInt64(),
		"id":            data.ID.ValueString(),
	})
}
