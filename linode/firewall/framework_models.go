package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// FirewallDataSourceModel describes the Terraform resource data model to match the
// resource schema.
type FirewallDataSourceModel struct {
	ID             types.Int64   `tfsdk:"id"`
	Label          types.String  `tfsdk:"label"`
	Tags           types.Set     `tfsdk:"tags"`
	Disabled       types.Bool    `tfsdk:"disabled"`
	Inbound        []RuleModel   `tfsdk:"inbound"`
	InboundPolicy  types.String  `tfsdk:"inbound_policy"`
	Outbound       []RuleModel   `tfsdk:"outbound"`
	OutboundPolicy types.String  `tfsdk:"outbound_policy"`
	Linodes        types.Set     `tfsdk:"linodes"`
	NodeBalancers  types.Set     `tfsdk:"nodebalancers"`
	Devices        []DeviceModel `tfsdk:"devices"`
	Status         types.String  `tfsdk:"status"`
	Created        types.String  `tfsdk:"created"`
	Updated        types.String  `tfsdk:"updated"`
}

type RuleModel struct {
	Label    types.String `tfsdk:"label"`
	Action   types.String `tfsdk:"action"`
	Ports    types.String `tfsdk:"ports"`
	Protocol types.String `tfsdk:"protocol"`
	IPv4     types.List   `tfsdk:"ipv4"`
	IPv6     types.List   `tfsdk:"ipv6"`
}

type DeviceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	EntityID types.Int64  `tfsdk:"entity_id"`
	Type     types.String `tfsdk:"type"`
	Label    types.String `tfsdk:"label"`
	URL      types.String `tfsdk:"url"`
}

func (data *FirewallDataSourceModel) flattenFirewallForDataSource(
	ctx context.Context,
	firewall *linodego.Firewall,
	devices []linodego.FirewallDevice,
	ruleSet *linodego.FirewallRuleSet,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(firewall.ID))
	data.Status = types.StringValue(string(firewall.Status))
	data.Created = types.StringValue(firewall.Created.Format(helper.TIME_FORMAT))
	data.Updated = types.StringValue(firewall.Updated.Format(helper.TIME_FORMAT))

	linodes, diags := types.SetValueFrom(
		ctx,
		types.Int64Type,
		AggregateEntityIDs(devices, linodego.FirewallDeviceLinode),
	)
	if diags.HasError() {
		return diags
	}
	data.Linodes = linodes

	nodebalancers, diags := types.SetValueFrom(
		ctx,
		types.Int64Type,
		AggregateEntityIDs(devices, linodego.FirewallDeviceNodeBalancer),
	)
	if diags.HasError() {
		return diags
	}
	data.NodeBalancers = nodebalancers

	data.Devices = parseFirewallDevices(devices)

	tags, diags := types.SetValueFrom(ctx, types.StringType, firewall.Tags)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags
	data.Disabled = types.BoolValue(firewall.Status == linodego.FirewallDisabled)
	data.InboundPolicy = types.StringValue(ruleSet.InboundPolicy)
	data.OutboundPolicy = types.StringValue(ruleSet.OutboundPolicy)
	data.Label = types.StringValue(firewall.Label)

	if ruleSet.Inbound != nil {
		inBound, diags := FlattenFirewallRules(ctx, ruleSet.Inbound)
		if diags.HasError() {
			return diags
		}
		data.Inbound = inBound
	}

	if ruleSet.Outbound != nil {
		outBound, diags := FlattenFirewallRules(ctx, ruleSet.Outbound)
		if diags.HasError() {
			return diags
		}
		data.Outbound = outBound
	}

	return nil
}

func FlattenFirewallRules(
	ctx context.Context,
	rules []linodego.FirewallRule,
) ([]RuleModel, diag.Diagnostics) {
	rulesData := make([]RuleModel, len(rules))

	for i, rule := range rules {
		rulesData[i].Action = types.StringValue(rule.Action)
		rulesData[i].Label = types.StringValue(rule.Label)
		rulesData[i].Protocol = types.StringValue(string(rule.Protocol))
		rulesData[i].Ports = types.StringValue(rule.Ports)

		ipv4, diags := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv4)
		if diags.HasError() {
			return nil, diags
		}

		rulesData[i].IPv4 = ipv4

		ipv6, diags := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv6)
		if diags.HasError() {
			return nil, diags
		}

		rulesData[i].IPv6 = ipv6
	}
	return rulesData, nil
}

func AggregateEntityIDs(devices []linodego.FirewallDevice, entityType linodego.FirewallDeviceType) []int {
	results := make([]int, 0, len(devices))
	for _, device := range devices {
		if device.Entity.Type == entityType {
			results = append(results, device.Entity.ID)
		}
	}
	return results
}

func parseFirewallDevices(devices []linodego.FirewallDevice) []DeviceModel {
	governedDevices := make([]DeviceModel, len(devices))
	for i, device := range devices {
		governedDevices[i].ID = types.Int64Value(int64(device.ID))
		governedDevices[i].EntityID = types.Int64Value(int64(device.Entity.ID))
		governedDevices[i].Type = types.StringValue(string(device.Entity.Type))
		governedDevices[i].Label = types.StringValue(device.Entity.Label)
		governedDevices[i].URL = types.StringValue(device.Entity.URL)
	}

	return governedDevices
}
