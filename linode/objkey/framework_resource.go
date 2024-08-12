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

var _ resource.ResourceWithValidateConfig = &Resource{}

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
	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := plan.GetCreateOptions(ctx)

	tflog.Debug(ctx, "client.CreateObjectStorageKey(...)", map[string]any{
		"options": createOpts,
	})
	key, err := client.CreateObjectStorageKey(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Object Storage Key",
			err.Error(),
		)
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"key_id": key.ID,
		"label":  key.Label,
	})

	plan.FlattenObjectStorageKey(ctx, key, true, &resp.Diagnostics)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(key.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

	data.FlattenObjectStorageKey(ctx, key, false, &resp.Diagnostics)
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

	updateOpts, shouldUpdate := plan.GetUpdateOptions(ctx, state)
	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateObjectStorageKey(...)", map[string]any{
			"options": updateOpts,
		})
		key, err := r.Meta.Client.UpdateObjectStorageKey(ctx, id, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update Object Storage Key (%d)", id),
				err.Error())
			return
		}

		plan.FlattenObjectStorageKey(ctx, key, true, &resp.Diagnostics)
	}

	plan.CopyFrom(state, true)
	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

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
	tflog.Debug(ctx, "client.DeleteObjectStorageKey(...)")
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
