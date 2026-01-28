package volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// defaultEnabledOnCreate is a String plan modifier that sets the planned value to
// "enabled" during create when the attribute is omitted (null/unknown).
// It does nothing on updates (when a prior state value exists).
//
// This keeps the plan consistent with Create() behavior without affecting updates,
// where UseStateForUnknown preserves the existing state when the field is omitted.
type defaultEnabledOnCreate struct{}

func DefaultEnabledOnCreate() planmodifier.String { return defaultEnabledOnCreate{} }

func (m defaultEnabledOnCreate) Description(ctx context.Context) string {
	return "Defaults to \"enabled\" on create when encryption is omitted"
}

func (m defaultEnabledOnCreate) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultEnabledOnCreate) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If there is a prior state value, this is an update; do nothing.
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		return
	}
	// Create path: if the user omitted encryption (null/unknown), set planned value to "enabled".
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = types.StringValue("enabled")
	}
}
