package firewallrulesets

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/firewallruleset"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type RuleSetFilterModel struct {
	ID       types.String                          `tfsdk:"id"`
	Filters  frameworkfilter.FiltersModelType       `tfsdk:"filter"`
	RuleSets []firewallruleset.RuleSetBaseModel     `tfsdk:"rulesets"`
}

func (data *RuleSetFilterModel) parseRuleSets(
	ctx context.Context,
	rulesets []linodego.RuleSet,
	diags *diag.Diagnostics,
) {
	data.RuleSets = make([]firewallruleset.RuleSetBaseModel, len(rulesets))
	for i, rs := range rulesets {
		data.RuleSets[i].FlattenRuleSet(ctx, rs, diags, false)
	}
}
