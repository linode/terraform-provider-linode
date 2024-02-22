package objkey

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
				Name:   "linode_object_storage_key",
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
	tflog.Debug(ctx, "Create linode_object_storage_key")
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label: data.Label.ValueString(),
	}

	if data.BucketAccess != nil {
		accessSlice := make(
			[]linodego.ObjectStorageKeyBucketAccess,
			len(data.BucketAccess),
		)

		for i, v := range data.BucketAccess {
			accessSlice[i] = v.toLinodeObject()
		}

		createOpts.BucketAccess = &accessSlice
	}

	tflog.Info(ctx, "Creating the object storage key", map[string]any{
		"createOpts": createOpts,
	})
	key, err := client.CreateObjectStorageKey(ctx, createOpts)

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"key_id": key.ID,
		"label":  key.Label,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Object Storage Key",
			err.Error(),
		)
		return
	}

	data.FlattenObjectStorageKey(key, true)

	// IDs should always be overridden during creation (see #1085)
	data.ID = types.StringValue(strconv.Itoa(key.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_object_storage_key")
	client := r.Meta.Client
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"key_id": data.ID.ValueString(),
		"label":  data.Label.ValueString(),
	})

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching the object storage key")

	key, err := client.GetObjectStorageKey(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Object Storage Key",
				fmt.Sprintf(
					"Removing Object Storage Key with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Object Storage Key",
			err.Error(),
		)
		return
	}

	data.FlattenObjectStorageKey(key, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_object_storage_key")
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"key_id": state.ID.ValueString(),
		"label":  state.Label.ValueString(),
	})

	id := helper.StringToInt(plan.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateOpts linodego.ObjectStorageKeyUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Info(ctx, "Updating object storage key", map[string]any{
			"updateOpts": updateOpts,
		})
		key, err := r.Meta.Client.UpdateObjectStorageKey(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Object Storage Key (%d)", id),
				err.Error())
			return
		}

		plan.FlattenObjectStorageKey(key, true)
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_object_storage_key")
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"key_id": data.ID.ValueString(),
		"label":  data.Label.ValueString(),
	})

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	tflog.Info(ctx, "Deleting the object storage key")
	err := client.DeleteObjectStorageKey(ctx, id)
	if err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the Object Storage Key (%d)", id),
				err.Error(),
			)
		}
		return
	}

	// a settling cooldown to avoid expired tokens from being returned in listings
	// may be switched to event poller later
	time.Sleep(3 * time.Second)
}
