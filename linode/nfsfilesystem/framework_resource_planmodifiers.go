package nfsfilesystem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModifyPlan marks API-generated attributes unknown whenever an update is planned.
// The framework only does this when the proposed state differs from prior state,
// which misses changes introduced by attribute plan modifiers (e.g. removing
// max_file_count), leaving stale values like stats in the plan.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.Plan.Raw.Equal(req.State.Raw) {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("stats"), types.ObjectUnknown(statsObjectType.AttrTypes))...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated"), timetypes.NewRFC3339Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("mount_target_fqdn"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("snapshot_usage_bytes"), types.Int64Unknown())...)
}

// maxFileCountConfiguredKey is the private state key recording that the user has
// explicitly configured max_file_count.
const maxFileCountConfiguredKey = "max_file_count_configured"

type privateStateSetter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// setMaxFileCountConfigured records in private state whether max_file_count is
// explicitly configured, so that its later removal can be planned as null.
func setMaxFileCountConfigured(ctx context.Context, private privateStateSetter, configValue types.Int64) diag.Diagnostics {
	if !configValue.IsNull() {
		return private.SetKey(ctx, maxFileCountConfiguredKey, []byte("true"))
	}
	return private.SetKey(ctx, maxFileCountConfiguredKey, nil)
}

// maxFileCountPlanModifier keeps the API default for max_file_count when it is
// never configured, and plans null (unlimited) when a previously configured
// value is removed from the configuration.
type maxFileCountPlanModifier struct{}

func (m maxFileCountPlanModifier) Description(context.Context) string {
	return "Keeps the prior state value when unconfigured, or plans null when a configured value is removed."
}

func (m maxFileCountPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m maxFileCountPlanModifier) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Leave the value unknown on create so the API default is used, and
	// respect any explicitly configured value.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}

	configured, diags := req.Private.GetKey(ctx, maxFileCountConfiguredKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if configured != nil {
		resp.PlanValue = types.Int64Null()
		return
	}

	resp.PlanValue = req.StateValue
}
