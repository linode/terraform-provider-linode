package iamuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type Resource struct {
	helper.BaseResource
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_iam_user",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan IAMUserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := plan.Username.ValueString()
	ctx = populateLogAttributes(ctx, username)

	client := r.Meta.Client

	perms, err := client.GetUserRolePermissions(ctx, username)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get User Roles", err.Error())
		return
	}

	updateOpts := perms.GetUpdateOptions()
	shouldUpdate := plan.CreateChanges(ctx, perms, &updateOpts, &resp.Diagnostics)

	plan.KeepOrUpdate(ctx, perms, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateUserRolePermissions(...)", map[string]any{
			"username": username,
			"options":  updateOpts,
		})

		newPerms, err := client.UpdateUserRolePermissions(ctx, username, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Update User", err.Error())
			return
		}

		plan.KeepOrUpdate(ctx, newPerms, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client

	var state IAMUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.Username, resp) {
		return
	}

	user := state.Username.ValueString()

	ctx = populateLogAttributes(ctx, user)

	perms, err := client.GetUserRolePermissions(ctx, user)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"User not found",
				fmt.Sprintf("Removing %s perms from the state", user),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to Get User Perms", err.Error())
		return
	}

	state.KeepOrUpdate(ctx, perms, true, &resp.Diagnostics)
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
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state IAMUserResourceModel
	diags := resp.Diagnostics

	diags.Append(req.State.Get(ctx, &state)...)
	diags.Append(req.Plan.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}

	username := plan.Username.ValueString()
	ctx = populateLogAttributes(ctx, username)

	client := r.Meta.Client

	updateOpts := linodego.UserRolePermissionsUpdateOptions{}
	shouldUpdate := false

	if !state.AccountAccess.Equal(plan.AccountAccess) {
		diags.Append(plan.AccountAccess.ElementsAs(ctx, &updateOpts.AccountAccess, false)...)
		if diags.HasError() {
			return
		}
		shouldUpdate = true
	}

	if !state.EntityAccessHasChanges(ctx, plan, &diags) {
		diags.Append(plan.EntityAccess.ElementsAs(ctx, &updateOpts.EntityAccess, false)...)
		if diags.HasError() {
			return
		}
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateUserRolePermissions(...)", map[string]any{
			"username": username,
			"options":  updateOpts,
		})

		perms, err := client.UpdateUserRolePermissions(ctx, username, updateOpts)
		if err != nil {
			diags.AddError("Failed to Update User", err.Error())
			return
		}
		plan.KeepOrUpdate(ctx, perms, true, &diags)
		if diags.HasError() {
			return
		}
	}

	diags.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete has no impact "+r.Config.Name)
	var state IAMUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// TODO IMPORT STATE MAYBE? need to import the resource instead of creating
func populateLogAttributes(ctx context.Context, username string) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"username": username,
	})
}
