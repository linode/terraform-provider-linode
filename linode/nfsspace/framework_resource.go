package nfsspace

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

// deletePollInterval is how often the NFS Space deletion status is polled.
const deletePollInterval = 5 * time.Second

// deleteTimeout is the maximum amount of time to wait for an NFS Space to be
// fully deleted before returning.
const deleteTimeout = 15 * time.Minute

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:          "linode_nfs_space",
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

	createOpts := linodego.NFSSpaceCreateOptions{
		Label:       data.Label.ValueString(),
		Description: data.Description.ValueStringPointer(),
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createOpts.Tags = &tags
	}

	tflog.Debug(ctx, "client.CreateNFSSpace(...)", map[string]any{
		"options": createOpts,
	})
	space, err := client.CreateNFSSpace(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create NFS Space.",
			err.Error(),
		)
		return
	}

	space, err = client.WaitForNFSSpaceStatus(
		ctx,
		space.ID,
		linodego.NFSSpaceStatusActive,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to wait for NFS Space %d to become active.", space.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Flatten(ctx, space, true)...)
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

	space, err := client.GetNFSSpace(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NFS Space no longer exists.",
				fmt.Sprintf(
					"Removing NFS Space with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the NFS Space.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.Flatten(ctx, space, false)...)
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
	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.NFSSpaceUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		shouldUpdate = true
		updateOpts.Label = plan.Label.ValueStringPointer()
	}

	if !state.Description.Equal(plan.Description) {
		shouldUpdate = true
		updateOpts.Description = plan.Description.ValueStringPointer()
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
		tflog.Debug(ctx, "client.UpdateNFSSpace(...)", map[string]any{
			"options": updateOpts,
		})
		space, err := client.UpdateNFSSpace(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update NFS Space (%d).", id),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(plan.Flatten(ctx, space, false)...)
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
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.DeleteNFSSpace(...)")
	if err := deleteNFSSpaceWithRetry(ctx, client, id); err != nil {
		if !linodego.IsNotFound(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the NFS Space (%d).", id),
				err.Error(),
			)
		}
		return
	}

	// NFS Space deletion is asynchronous; poll until the API reports the
	// space no longer exists or deleting.
	if err := waitForNFSSpaceDeletion(ctx, client, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to wait for the NFS Space (%d) to be deleted.", id),
			err.Error(),
		)
		return
	}
}

// deleteNFSSpaceWithRetry sends the delete request, retrying while the API
// reports a 409 (the space still has filesystems). A child filesystem that is
// still in a "deleting" status can keep the space associated for a short
// window after the filesystem GET already returns 404, so we retry until the
// association is released or the timeout is reached.
func deleteNFSSpaceWithRetry(ctx context.Context, client *linodego.Client, id int) error {
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	ticker := time.NewTicker(deletePollInterval)
	defer ticker.Stop()

	for {
		err := client.DeleteNFSSpace(ctx, id)
		if err == nil || !linodego.ErrHasStatus(err, http.StatusConflict) {
			return err
		}

		tflog.Debug(ctx, "NFS Space deletion returned 409, retrying", map[string]any{
			"error": err.Error(),
		})

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting to delete NFS Space %d: %w", id, ctx.Err())
		}
	}
}

// waitForNFSSpaceDeletion polls the API until the given NFS Space no longer
// exists or the deletion timeout is reached. Deletion is asynchronous: the
// space lingers in a "deleting" status before it is fully removed, so we must
// keep polling until the API reports it as gone (404) or no longer accessible
// (403).
func waitForNFSSpaceDeletion(ctx context.Context, client *linodego.Client, id int) error {
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	ticker := time.NewTicker(deletePollInterval)
	defer ticker.Stop()

	for {
		_, err := client.GetNFSSpace(ctx, id)
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
			return fmt.Errorf("timed out waiting for NFS Space %d to be deleted: %w", id, ctx.Err())
		}
	}
}

func populateLogAttributes(ctx context.Context, id string) context.Context {
	return tflog.SetField(ctx, "id", id)
}
