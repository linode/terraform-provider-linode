package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CaseInsensitiveSetPlanModifier ensures that a set does not trigger
// diffs on planned values with different cases.
//
// NOTE: This is not implemented as custom type because custom type semantic equality
// checks do not have granular control over diffs.
func CaseInsensitiveSetPlanModifier() planmodifier.Set {
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
	resultList := make([]string, 0)

	for _, elem := range req.StateValue.Elements() {
		elemStr := elem.(types.String).ValueString()
		oldEntryMap[strings.ToLower(elemStr)] = elemStr
	}

	for _, elem := range req.PlanValue.Elements() {
		elemStr := elem.(types.String).ValueString()
		oldElem, ok := oldEntryMap[strings.ToLower(elemStr)]

		if !ok {
			resultList = append(resultList, elemStr)
			continue
		}

		resultList = append(resultList, oldElem)
	}

	v, diags := types.SetValueFrom(ctx, types.StringType, resultList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = v
}
