package iamuser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

type IAMUserDataSourceModel struct {
	Username      types.String   `tfsdk:"username"`
	AccountAccess types.List     `tfsdk:"account_access"`
	EntityAccess  []EntityAccess `tfsdk:"entity_access"`
}

type EntityAccess struct {
	ID    types.Int64  `tfsdk:"id"`
	Type  types.String `tfsdk:"type"`
	Roles types.List   `tfsdk:"roles"`
}

type IAMUserResourceModel struct {
	Username      types.String `tfsdk:"username"`
	AccountAccess types.List   `tfsdk:"account_access"`
	EntityAccess  types.List   `tfsdk:"entity_access"`
}

func (state *IAMUserResourceModel) EntityAccessHasChanges(
	ctx context.Context, plan IAMUserResourceModel, diags *diag.Diagnostics,
) bool {
	oldEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, state.EntityAccess)
	diags.Append(newDiags...)
	newEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, plan.EntityAccess)
	diags.Append(newDiags...)

	if newDiags.HasError() {
		diags.Append(newDiags...)
	}

	return (!oldEntities.Equal(newEntities))
}

func (plan *IAMUserResourceModel) CreateChanges(
	ctx context.Context,
	perms *linodego.UserRolePermissions,
	updateOpts *linodego.UserRolePermissionsUpdateOptions,
	diags *diag.Diagnostics,
) bool {
	shouldUpdate := false
	oldEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, perms.EntityAccess)
	diags.Append(newDiags...)
	newEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, plan.EntityAccess)
	diags.Append(newDiags...)
	accAccess, newDiags := types.ListValueFrom(ctx, types.StringType, perms.AccountAccess)
	diags.Append(newDiags...)

	tflog.Debug(ctx, fmt.Sprintf(" !111111111111!!!!!!!!!!!!! %v ", diags), nil)

	if plan.EntityAccess.IsNull() || plan.EntityAccess.IsUnknown() {
		updateOpts.EntityAccess = perms.EntityAccess
		plan.EntityAccess = oldEntities
	} else if !oldEntities.Equal(newEntities) {
		diags.Append(newEntities.ElementsAs(ctx, &updateOpts.EntityAccess, false)...)
		tflog.Debug(ctx, fmt.Sprintf(" !!!!!!!!!!!!!! %v ", updateOpts), nil)
		tflog.Debug(ctx, fmt.Sprintf(" !!!!!!!!!!!!!! %v ", newEntities), nil)
		tflog.Debug(ctx, fmt.Sprintf(" !!!!!!!!!!!!!! %v ", diags), nil)
		shouldUpdate = true
	}

	if plan.AccountAccess.IsNull() || plan.AccountAccess.IsUnknown() {
		updateOpts.AccountAccess = perms.AccountAccess
		plan.AccountAccess = accAccess
	} else if !plan.AccountAccess.Equal(accAccess) {
		diags.Append(plan.AccountAccess.ElementsAs(ctx, &updateOpts.AccountAccess, false)...)
		shouldUpdate = true
	}

	return shouldUpdate
}

func (plan *IAMUserResourceModel) KeepOrUpdate(
	ctx context.Context,
	perms *linodego.UserRolePermissions,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	accountAccess, newDiags := types.ListValueFrom(ctx, types.StringType, perms.AccountAccess)
	diags.Append(newDiags...)
	entityPerms, newDiags := types.ListValueFrom(ctx, EntityAccessType, perms.EntityAccess)
	diags.Append(newDiags...)
	plan.AccountAccess = accountAccess
	plan.EntityAccess = entityPerms
}

func (data *IAMUserDataSourceModel) ParseIAMUserModel(
	ctx context.Context,
	perms *linodego.UserRolePermissions,
) diag.Diagnostics {
	accountAccess, diags := types.ListValueFrom(ctx, types.StringType, perms.AccountAccess)
	if diags.HasError() {
		return diags
	}
	data.AccountAccess = accountAccess
	entities := make([]EntityAccess, len(perms.EntityAccess))

	for i, r := range perms.EntityAccess {
		entities[i].ID = types.Int64Value(int64(r.ID))
		entities[i].Type = types.StringValue(r.Type)
		roles, diags := types.ListValueFrom(ctx, types.StringType, r.Roles)
		if diags.HasError() {
			return diags
		}
		entities[i].Roles = roles
	}

	data.EntityAccess = entities

	return nil
}
