package databaseaccesscontrols

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Name:   "linode_database_access_controls",
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

	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID := helper.FrameworkSafeInt64ToInt(plan.DatabaseID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	dbType := plan.DatabaseType.ValueString()

	ctx = populateLogAttributes(ctx, dbID, dbType)

	allowList := expandAllowList(ctx, plan.AllowList, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := updateDBAllowListByEngine(ctx, client, dbType, dbID, allowList); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update allow_list for database %d", dbID),
			err.Error(),
		)
		return
	}

	plan.Flatten(dbID, dbType, allowList, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(formatID(dbID, dbType))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	dbID, dbType, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse database ID", err.Error())
		return
	}

	ctx = populateLogAttributes(ctx, dbID, dbType)

	allowList, err := getDBAllowListByEngine(ctx, client, dbType, dbID)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Database No Longer Exists",
				fmt.Sprintf(
					"Removing allow_list %q from state because it no longer exists",
					state.ID.ValueString(),
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get allow list for database %d", dbID),
			err.Error(),
		)
		return
	}

	state.Flatten(dbID, dbType, allowList, false, &resp.Diagnostics)
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

	var plan, state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID := helper.FrameworkSafeInt64ToInt(plan.DatabaseID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	dbType := plan.DatabaseType.ValueString()

	ctx = populateLogAttributes(ctx, dbID, dbType)

	if !plan.AllowList.Equal(state.AllowList) {
		allowList := expandAllowList(ctx, plan.AllowList, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := updateDBAllowListByEngine(ctx, client, dbType, dbID, allowList); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update allow_list for database %d", dbID),
				err.Error(),
			)
			return
		}
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
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID, dbType, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse database ID", err.Error())
		return
	}

	ctx = populateLogAttributes(ctx, dbID, dbType)

	if err := updateDBAllowListByEngine(ctx, client, dbType, dbID, []string{}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update allow_list for database %d", dbID),
			err.Error(),
		)
		return
	}
}

// expandAllowList converts the given Terraform set into a slice of
// allow list entries.
func expandAllowList(
	ctx context.Context,
	allowList types.Set,
	diags *diag.Diagnostics,
) []string {
	if allowList.IsNull() || allowList.IsUnknown() {
		return []string{}
	}

	result := make([]string, 0, len(allowList.Elements()))

	diags.Append(allowList.ElementsAs(ctx, &result, false)...)

	return result
}

func populateLogAttributes(ctx context.Context, dbID int, dbType string) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"database_id":   dbID,
		"database_type": dbType,
	})
}
