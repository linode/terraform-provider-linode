package iamuser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type IAMUserDataSourceModel struct {
	Username      types.String        `tfsdk:"username"`
	AccountAccess types.List          `tfsdk:"account_access"`
	EntityAccess  []EntityAccessModel `tfsdk:"entity_access"`
}

type EntityAccessModel struct {
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

	if plan.EntityAccess.IsNull() || plan.EntityAccess.IsUnknown() {
		updateOpts.EntityAccess = perms.EntityAccess
		plan.EntityAccess = oldEntities
	} else if !oldEntities.Equal(newEntities) {
		var entities []EntityAccessModel

		diags.Append(newEntities.ElementsAs(ctx, &entities, false)...)

		updateOpts.EntityAccess = helper.MapSlice[EntityAccessModel, linodego.UserAccess](
			entities,
			func(e EntityAccessModel) linodego.UserAccess {
				return e.ToLinodego(ctx, diags)
			},
		)

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

func (plan *IAMUserResourceModel) UpdateChanges(
	ctx context.Context,
	state *IAMUserResourceModel,
	updateOpts *linodego.UserRolePermissionsUpdateOptions,
	diags *diag.Diagnostics,
) bool {
	shouldUpdate := false
	oldEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, state.EntityAccess)
	diags.Append(newDiags...)
	newEntities, newDiags := types.ListValueFrom(ctx, EntityAccessType, plan.EntityAccess)
	diags.Append(newDiags...)

	whichEntities := oldEntities
	if !oldEntities.Equal(newEntities) {
		whichEntities = newEntities
		shouldUpdate = true
	}
	var entities []EntityAccessModel

	diags.Append(whichEntities.ElementsAs(ctx, &entities, false)...)

	updateOpts.EntityAccess = helper.MapSlice[EntityAccessModel, linodego.UserAccess](
		entities,
		func(e EntityAccessModel) linodego.UserAccess {
			return e.ToLinodego(ctx, diags)
		},
	)

	if plan.AccountAccess.IsNull() || plan.AccountAccess.IsUnknown() {
		diags.Append(state.AccountAccess.ElementsAs(ctx, &updateOpts.AccountAccess, false)...)
		plan.AccountAccess = state.AccountAccess
	} else if !plan.AccountAccess.Equal(state.AccountAccess) {
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

	plan.AccountAccess = accountAccess
	entityPerms, newDiags := types.ListValueFrom(ctx, EntityAccessType, helper.MapSlice(perms.EntityAccess, func(e linodego.UserAccess) EntityAccessModel {
		return flattenEntityAccess(ctx, e, diags)
	}))
	diags.Append(newDiags...)

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
	entities := make([]EntityAccessModel, len(perms.EntityAccess))

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

func (m *EntityAccessModel) ToLinodego(ctx context.Context, diags *diag.Diagnostics) linodego.UserAccess {
	var roles []string
	diags.Append(m.Roles.ElementsAs(ctx, &roles, false)...)

	return linodego.UserAccess{
		ID:    helper.FrameworkSafeInt64ToInt(m.ID.ValueInt64(), diags),
		Type:  m.Type.ValueString(),
		Roles: roles,
	}
}

func flattenEntityAccess(ctx context.Context, access linodego.UserAccess, diags *diag.Diagnostics) EntityAccessModel {
	roles, d := types.ListValueFrom(ctx, types.StringType, access.Roles)
	diags.Append(d...)

	return EntityAccessModel{
		ID:    types.Int64Value(int64(access.ID)),
		Type:  types.StringValue(access.Type),
		Roles: roles,
	}
}
