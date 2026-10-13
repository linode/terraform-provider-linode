package nfsfilesystem

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

const deletePollInterval = 5 * time.Second

// deleteTimeout is the maximum amount of time to wait for an NFS Filesystem to
// be fully deleted before returning.
const deleteTimeout = 15 * time.Minute

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:          "linode_nfs_filesystem",
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

	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.NFSFilesystemCreateOptions{
		Label:            data.Label.ValueString(),
		Region:           data.Region.ValueString(),
		MaxCapacityBytes: data.MaxCapacityBytes.ValueInt64(),
	}

	if !data.ProtocolVersions.IsNull() && !data.ProtocolVersions.IsUnknown() {
		var protocols []string
		resp.Diagnostics.Append(data.ProtocolVersions.ElementsAs(ctx, &protocols, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createOpts.ProtocolVersions = linodego.Pointer(helper.MapSlice(protocols, func(p string) linodego.NFSProtocolVersion {
			return linodego.NFSProtocolVersion(p)
		}))
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createOpts.Tags = &tags
	}

	tflog.Debug(ctx, "client.CreateNFSFilesystem(...)", map[string]any{
		"options": createOpts,
	})
	fs, err := client.CreateNFSFilesystem(ctx, spaceID, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create NFS Filesystem.",
			err.Error(),
		)
		return
	}

	fs, err = client.WaitForNFSFilesystemStatus(
		ctx,
		spaceID,
		fs.ID,
		linodego.NFSFilesystemStatusActive,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to wait for NFS Filesystem %d to become active.", fs.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Flatten(ctx, fs, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	ctx = populateLogAttributes(ctx, data.ID.ValueString())

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	fs, err := client.GetNFSFilesystem(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NFS Filesystem no longer exists.",
				fmt.Sprintf(
					"Removing NFS Filesystem with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the NFS Filesystem.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Flatten(ctx, fs, false)...)
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

	ctx = populateLogAttributes(ctx, state.ID.ValueString())

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	spaceID := helper.FrameworkSafeInt64ToInt(state.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.NFSFilesystemUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueStringPointer()
	}

	if !state.MaxCapacityBytes.Equal(plan.MaxCapacityBytes) {
		shouldUpdate = true
		updateOpts.MaxCapacityBytes = plan.MaxCapacityBytes.ValueInt64Pointer()
	}

	if !state.Tags.Equal(plan.Tags) {
		shouldUpdate = true
		var tags []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateOpts.Tags = &tags
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateNFSFilesystem(...)", map[string]any{
			"options": updateOpts,
		})
		fs, err := client.UpdateNFSFilesystem(ctx, spaceID, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update NFS Filesystem (%d).", id),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(plan.Flatten(ctx, fs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.CopyFrom(state.BaseModel, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	client := r.Meta.Client
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data.ID.ValueString())

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteNFSFilesystem(...)")
	err := client.DeleteNFSFilesystem(ctx, spaceID, id)
	if err != nil {
		if !linodego.IsNotFound(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the NFS Filesystem (%d).", id),
				err.Error(),
			)
		}
		return
	}

	// NFS Filesystem deletion is asynchronous; poll until the API reports the
	// filesystem no longer exists or deleting.
	if err := waitForNFSFilesystemDeletion(ctx, client, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to wait for the NFS Filesystem (%d) to be deleted.", id),
			err.Error(),
		)
		return
	}
}

// waitForNFSFilesystemDeletion polls the API until the given NFS Filesystem no
// longer exists or the deletion timeout is reached. Deletion is asynchronous:
// the filesystem lingers in a "deleting" status before it is fully removed, so
// we must keep polling until the API reports it as gone (404) or no longer
// accessible (403). Returning early while the filesystem still exists would
// cause the parent NFS Space deletion to fail with a 409.
func waitForNFSFilesystemDeletion(ctx context.Context, client *linodego.Client, id int) error {
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	ticker := time.NewTicker(deletePollInterval)
	defer ticker.Stop()

	for {
		_, err := client.GetNFSFilesystem(ctx, id)
		if err != nil {
			if linodego.IsNotFound(err) || linodego.ErrHasStatus(err, http.StatusForbidden) {
				return nil
			}
			return err
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for NFS Filesystem %d to be deleted: %w", id, ctx.Err())
		}
	}
}

func populateLogAttributes(ctx context.Context, id string) context.Context {
	return tflog.SetField(ctx, "id", id)
}
