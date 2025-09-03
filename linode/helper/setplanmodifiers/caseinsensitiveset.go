package setplanmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CaseInsensitiveSet ensures that a set does not trigger
// diffs on planned values with different cases.
//
// NOTE: This is not implemented as custom type because custom type semantic equality
// checks do not have granular control over diffs.
func CaseInsensitiveSet() planmodifier.Set {
	return &caseInsensitiveSetPlanModifier{}
}

type caseInsensitiveSetPlanModifier struct{}

func (d *caseInsensitiveSetPlanModifier) Description(ctx context.Context) string {
	return "Ensures a set does not trigger diffs on planned values with different cases."
}

func (d *caseInsensitiveSetPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *caseInsensitiveSetPlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	oldEntryMap := make(map[string]string)
	resultList := make([]attr.Value, 0) // Use attr.Value to store both known and unknown values

	// Store existing state values in a map, if they are known
	for _, elem := range req.StateValue.Elements() {
		strElem, ok := elem.(types.String)
		if !ok || strElem.IsUnknown() || strElem.IsNull() {
			continue // Skip unknown or null values
		}
		oldEntryMap[strings.ToLower(strElem.ValueString())] = strElem.ValueString()
	}

	// Process planned values, ensuring unknown values are preserved
	for _, elem := range req.PlanValue.Elements() {
		strElem, ok := elem.(types.String)
		if !ok {
			resultList = append(resultList, elem) // Preserve unknown or incompatible values
			continue
		}
		if strElem.IsUnknown() || strElem.IsNull() {
			resultList = append(resultList, strElem) // Keep unknown values unchanged
			continue
		}

		// Normalize case if an old value exists, otherwise keep the new value
		oldElem, ok := oldEntryMap[strings.ToLower(strElem.ValueString())]
		if !ok {
			resultList = append(resultList, strElem)
			continue
		}

		resultList = append(resultList, types.StringValue(oldElem))
	}

	// Convert resultList back into a Set type
	v, diags := types.SetValueFrom(ctx, types.StringType, resultList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = v
}
