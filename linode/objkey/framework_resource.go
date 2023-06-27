package objkey

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_object_storage_key",
				IDType: types.Int64Type,
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

	key, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Object Storage Key",
			err.Error(),
		)
		return
	}

	data.parseKey(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.Meta.Client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := client.GetObjectStorageKey(ctx, int(data.ID.ValueInt64()))
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

	data.parseKey(key)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var updateOpts linodego.ObjectStorageKeyUpdateOptions
	shouldUpdate := false

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if shouldUpdate {
		key, err := r.Meta.Client.UpdateObjectStorageKey(
			ctx,
			int(state.ID.ValueInt64()),
			updateOpts,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Object Storage Key (%d)", state.ID.ValueInt64()),
				err.Error())
			return
		}

		state.parseKey(key)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client
	err := client.DeleteObjectStorageKey(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the Object Storage Key (%d)", data.ID.ValueInt64()),
				err.Error(),
			)
		}
		return
	}

	// a settling cooldown to avoid expired tokens from being returned in listings
	// may be switched to event poller later
	time.Sleep(3 * time.Second)
}
