package firewallrulesexpansion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// FirewallRulesExpansionModel describes the data model for the expanded
// firewall rules data source.
type FirewallRulesExpansionModel struct {
	ID             types.String         `tfsdk:"id"`
	FirewallID     types.Int64          `tfsdk:"firewall_id"`
	Inbound        []firewall.RuleModel `tfsdk:"inbound"`
	InboundPolicy  types.String         `tfsdk:"inbound_policy"`
	Outbound       []firewall.RuleModel `tfsdk:"outbound"`
	OutboundPolicy types.String         `tfsdk:"outbound_policy"`
	Version        types.Int64          `tfsdk:"version"`
}

func (data *FirewallRulesExpansionModel) parseExpansion(
	ctx context.Context,
	firewallID int,
	ruleSet linodego.FirewallRuleSet,
	diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, data.FirewallID.String(), false)

	data.Inbound = firewall.FlattenFirewallRules(ctx, ruleSet.Inbound, nil, false, diags)
	if diags.HasError() {
		return
	}

	data.Outbound = firewall.FlattenFirewallRules(ctx, ruleSet.Outbound, nil, false, diags)
	if diags.HasError() {
		return
	}

	data.InboundPolicy = helper.KeepOrUpdateString(data.InboundPolicy, ruleSet.InboundPolicy, false)
	data.OutboundPolicy = helper.KeepOrUpdateString(data.OutboundPolicy, ruleSet.OutboundPolicy, false)
	data.Version = helper.KeepOrUpdateInt64(data.Version, int64(ruleSet.Version), false)
}
