package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// DomainRecordTargetUseStateIfSematicEquals set the state value to plan
// when the plan value will be change to the state value by the API.
// e.g. "sometarget" ==> "sometarget.example.com" when the domain is "example.com"
func DomainRecordTargetUseStateIfSematicEquals() planmodifier.String {
	return &domainRecordTargetUseStateIfSematicEqualsPlanModifier{}
}

type domainRecordTargetUseStateIfSematicEqualsPlanModifier struct{}

func (d *domainRecordTargetUseStateIfSematicEqualsPlanModifier) Description(ctx context.Context) string {
	return "Ensures a set does not trigger diffs on planned values with different cases."
}

func (d *domainRecordTargetUseStateIfSematicEqualsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *domainRecordTargetUseStateIfSematicEqualsPlanModifier) PlanModifyString(
	ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	// bring state into plan if the plan is semantic equivalent to the state
	if !req.StateValue.IsNull() {
		provisioned := req.StateValue.ValueString()
		declared := req.PlanValue.ValueString()
		if len(strings.Split(declared, ".")) == 1 && strings.Contains(provisioned, declared) {
			resp.PlanValue = req.StateValue
		}
	}
}
