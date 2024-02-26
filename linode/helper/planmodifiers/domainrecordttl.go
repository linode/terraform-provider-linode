package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// DomainRecordTTLUseStateIfPlanCanBeRoundedToState can modify the plan to round the
// TTL value to nearest validated value. A warning will be inserted if the configured
// TTL is not a validated value.
func DomainRecordTTLUseStateIfPlanCanBeRoundedToState() planmodifier.Int64 {
	return &domainRecordTTLRoundingPlanModifier{}
}

type domainRecordTTLRoundingPlanModifier struct{}

func (d *domainRecordTTLRoundingPlanModifier) Description(ctx context.Context) string {
	return "Ensures a set does not trigger diffs on planned values with different cases."
}

func (d *domainRecordTTLRoundingPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *domainRecordTTLRoundingPlanModifier) PlanModifyInt64(
	ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response,
) {
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		configuredTTL := req.ConfigValue.ValueInt64()
		configuredTTLAccepted := false
		for _, acceptedTTL := range helper.AcceptedTTLs {
			if int64(acceptedTTL) == configuredTTL {
				configuredTTLAccepted = true
			}
		}
		if !configuredTTLAccepted {
			resp.Diagnostics.AddWarning(
				"Invalid TTL Value",
				"An invalid value of TTL is rounded to the nearest valid value, "+
					"but it will no longer be accepted in the next major version of "+
					fmt.Sprintf(
						"Linode Terraform provider. The validated values are: %v",
						helper.AcceptedTTLs,
					),
			)
		}
	}

	rounder := func(n int64) int64 {
		if n == 0 {
			return 0
		}

		for _, value := range helper.AcceptedTTLs {
			if n <= int64(value) {
				return int64(value)
			}
		}
		return int64(helper.AcceptedTTLs[len(helper.AcceptedTTLs)-1])
	}

	// bring state into plan if plan can be rounded to state
	if !req.StateValue.IsNull() {
		if rounder(req.PlanValue.ValueInt64()) == req.StateValue.ValueInt64() {
			resp.PlanValue = req.StateValue
		}
	}
}
