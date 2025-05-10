package firewalltemplate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/firewall"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// FirewallTemplateDataSourceModel describes the Terraform
// resource data model to match the resource schema.
type FirewallTemplateBaseModel struct {
	Slug           types.String         `tfsdk:"slug"`
	Inbound        []firewall.RuleModel `tfsdk:"inbound"`
	InboundPolicy  types.String         `tfsdk:"inbound_policy"`
	Outbound       []firewall.RuleModel `tfsdk:"outbound"`
	OutboundPolicy types.String         `tfsdk:"outbound_policy"`
}

// FirewallTemplateDataSourceModel describes the Terraform
// resource data model to match the resource schema.
type FirewallTemplateDataSourceModel struct {
	FirewallTemplateBaseModel
	ID types.String `tfsdk:"id"`
}

func (data *FirewallTemplateBaseModel) FlattenFirewallTemplate(
	ctx context.Context,
	template linodego.FirewallTemplate,
	diags *diag.Diagnostics,
	preserveKnown bool,
) {
	data.Slug = helper.KeepOrUpdateString(data.Slug, template.Slug, preserveKnown)

	data.Inbound = firewall.FlattenFirewallRules(ctx, template.Rules.Inbound, nil, false, diags)
	if diags.HasError() {
		return
	}

	data.Outbound = firewall.FlattenFirewallRules(ctx, template.Rules.Outbound, nil, false, diags)
	if diags.HasError() {
		return
	}

	data.InboundPolicy = helper.KeepOrUpdateString(data.InboundPolicy, template.Rules.InboundPolicy, preserveKnown)
	data.OutboundPolicy = helper.KeepOrUpdateString(data.OutboundPolicy, template.Rules.OutboundPolicy, preserveKnown)
}

func (data *FirewallTemplateDataSourceModel) parseFirewallTemplate(
	ctx context.Context,
	template linodego.FirewallTemplate,
	diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, template.Slug, false)
	data.FirewallTemplateBaseModel.FlattenFirewallTemplate(ctx, template, diags, false)
}
