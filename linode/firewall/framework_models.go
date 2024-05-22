package firewall

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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

// FirewallResourceModel describes the Terraform resource data model to match the
// resource schema.
type FirewallResourceModel struct {
	ID             types.String      `tfsdk:"id"`
	Label          types.String      `tfsdk:"label"`
	Tags           types.Set         `tfsdk:"tags"`
	Disabled       types.Bool        `tfsdk:"disabled"`
	Inbound        []RuleModel       `tfsdk:"inbound"`
	InboundPolicy  types.String      `tfsdk:"inbound_policy"`
	Outbound       []RuleModel       `tfsdk:"outbound"`
	OutboundPolicy types.String      `tfsdk:"outbound_policy"`
	Linodes        types.Set         `tfsdk:"linodes"`
	NodeBalancers  types.Set         `tfsdk:"nodebalancers"`
	Devices        types.List        `tfsdk:"devices"`
	Status         types.String      `tfsdk:"status"`
	Created        timetypes.RFC3339 `tfsdk:"created"`
	Updated        timetypes.RFC3339 `tfsdk:"updated"`
}

type RuleModel struct {
	Label    types.String `tfsdk:"label"`
	Action   types.String `tfsdk:"action"`
	Ports    types.String `tfsdk:"ports"`
	Protocol types.String `tfsdk:"protocol"`
	IPv4     types.List   `tfsdk:"ipv4"`
	IPv6     types.List   `tfsdk:"ipv6"`
	// TODO: add Description
}

type DeviceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	EntityID types.Int64  `tfsdk:"entity_id"`
	Type     types.String `tfsdk:"type"`
	Label    types.String `tfsdk:"label"`
	URL      types.String `tfsdk:"url"`
}

func (data *RuleModel) Expands(ctx context.Context, diags *diag.Diagnostics) (rule linodego.FirewallRule) {
	rule.Action = data.Action.ValueString()
	rule.Label = data.Label.ValueString()
	rule.Ports = data.Ports.ValueString()
	rule.Protocol = linodego.NetworkProtocol(data.Protocol.ValueString())

	newDiags := data.IPv4.ElementsAs(ctx, &rule.Addresses.IPv4, false)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	newDiags = data.IPv6.ElementsAs(ctx, &rule.Addresses.IPv6, false)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	return
}

func ExpandFirewallRules(ctx context.Context, rulesList []RuleModel, diags *diag.Diagnostics) []linodego.FirewallRule {
	result := make([]linodego.FirewallRule, len(rulesList))
	for i, v := range rulesList {
		result[i] = v.Expands(ctx, diags)
	}

	return result
}

func isDisabled(firewall linodego.Firewall) bool {
	return firewall.Status == linodego.FirewallDisabled
}

func (data *FirewallResourceModel) ExpandFirewallRuleSet(
	ctx context.Context, diags *diag.Diagnostics,
) (rules linodego.FirewallRuleSet) {
	rules.Inbound = ExpandFirewallRules(ctx, data.Inbound, diags)
	if diags.HasError() {
		return
	}

	rules.Outbound = ExpandFirewallRules(ctx, data.Outbound, diags)
	if diags.HasError() {
		return
	}

	rules.InboundPolicy = data.InboundPolicy.ValueString()
	rules.OutboundPolicy = data.OutboundPolicy.ValueString()
	return
}

func (data *FirewallResourceModel) getCreateOptions(
	ctx context.Context, diags *diag.Diagnostics,
) (createOpts linodego.FirewallCreateOptions) {
	createOpts.Label = data.Label.ValueString()

	newDiags := data.Tags.ElementsAs(ctx, &createOpts.Tags, false)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	createOpts.Devices.Linodes = helper.ExpandFwInt64Set(data.Linodes, diags)
	if diags.HasError() {
		return
	}

	createOpts.Devices.NodeBalancers = helper.ExpandFwInt64Set(data.NodeBalancers, diags)
	if diags.HasError() {
		return
	}

	createOpts.Rules = data.ExpandFirewallRuleSet(ctx, diags)
	if diags.HasError() {
		return
	}

	createOpts.Rules.InboundPolicy = data.InboundPolicy.ValueString()
	createOpts.Rules.OutboundPolicy = data.OutboundPolicy.ValueString()

	return
}

func (plan *FirewallResourceModel) getUpdateOptions(
	ctx context.Context, state FirewallResourceModel, diags *diag.Diagnostics,
) (updateOpts linodego.FirewallUpdateOptions, shouldUpdate bool) {
	if !plan.Label.Equal(state.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}
	if !plan.Tags.Equal(state.Tags) {
		newDiags := plan.Tags.ElementsAs(ctx, &updateOpts.Tags, false)
		diags.Append(newDiags...)
		if diags.HasError() {
			return
		}
		shouldUpdate = true
	}
	if !plan.Disabled.Equal(state.Disabled) {
		updateOpts.Status = expandFirewallStatus(plan.Disabled.ValueBool())
		shouldUpdate = true
	}
	return
}

func (data *FirewallResourceModel) flattenDevices(
	ctx context.Context,
	devices []linodego.FirewallDevice,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	data.Linodes = helper.KeepOrUpdateIntSet(data.Linodes, AggregateEntityIDs(devices, linodego.FirewallDeviceLinode), preserveKnown, diags)
	if diags.HasError() {
		return
	}

	data.NodeBalancers = helper.KeepOrUpdateIntSet(data.NodeBalancers, AggregateEntityIDs(devices, linodego.FirewallDeviceNodeBalancer), preserveKnown, diags)
	if diags.HasError() {
		return
	}

	deviceModels := FlattenFirewallDevices(devices)
	devicesList, newDiags := types.ListValueFrom(ctx, deviceObjectType, deviceModels)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Devices = helper.KeepOrUpdateValue(data.Devices, devicesList, preserveKnown)
}

func (data *FirewallResourceModel) flattenRules(
	ctx context.Context,
	ruleSet *linodego.FirewallRuleSet,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	inboundRules, newDiags := FlattenFirewallRules(ctx, ruleSet.Inbound, data.Inbound, preserveKnown)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Inbound = inboundRules

	outboundRules, newDiags := FlattenFirewallRules(ctx, ruleSet.Outbound, data.Outbound, preserveKnown)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}
	data.Outbound = outboundRules
}

func (data *FirewallResourceModel) flattenFirewallForResource(
	firewall *linodego.Firewall,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(firewall.ID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, firewall.Label, preserveKnown)
	data.Tags = helper.KeepOrUpdateStringSet(data.Tags, firewall.Tags, preserveKnown, diags)
	if diags.HasError() {
		return
	}

	data.Disabled = helper.KeepOrUpdateBool(data.Disabled, isDisabled(*firewall), preserveKnown)

	data.InboundPolicy = helper.KeepOrUpdateString(data.InboundPolicy, firewall.Rules.InboundPolicy, preserveKnown)
	data.OutboundPolicy = helper.KeepOrUpdateString(data.OutboundPolicy, firewall.Rules.OutboundPolicy, preserveKnown)

	data.Status = helper.KeepOrUpdateString(data.Status, string(firewall.Status), preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, timetypes.NewRFC3339TimePointerValue(firewall.Created), preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, timetypes.NewRFC3339TimePointerValue(firewall.Updated), preserveKnown)
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

	data.Devices = FlattenFirewallDevices(devices)

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
		inBound, diags := FlattenFirewallRules(ctx, ruleSet.Inbound, nil, false)
		if diags.HasError() {
			return diags
		}
		data.Inbound = inBound
	}

	if ruleSet.Outbound != nil {
		outBound, diags := FlattenFirewallRules(ctx, ruleSet.Outbound, nil, false)
		if diags.HasError() {
			return diags
		}
		data.Outbound = outBound
	}

	return nil
}

func (data *FirewallResourceModel) CopyFrom(
	other FirewallResourceModel,
	preserveKnown bool,
) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.Tags = helper.KeepOrUpdateValue(data.Tags, other.Tags, preserveKnown)
	data.Disabled = helper.KeepOrUpdateValue(data.Disabled, other.Disabled, preserveKnown)
	data.InboundPolicy = helper.KeepOrUpdateValue(data.InboundPolicy, other.InboundPolicy, preserveKnown)
	data.OutboundPolicy = helper.KeepOrUpdateValue(data.OutboundPolicy, other.OutboundPolicy, preserveKnown)
	data.Linodes = helper.KeepOrUpdateValue(data.Linodes, other.Linodes, preserveKnown)
	data.NodeBalancers = helper.KeepOrUpdateValue(data.NodeBalancers, other.NodeBalancers, preserveKnown)
	data.Devices = helper.KeepOrUpdateValue(data.Devices, other.Devices, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, other.Updated, preserveKnown)

	if !preserveKnown {
		data.Inbound = other.Inbound
		data.Outbound = other.Outbound
	}
}

func (state *FirewallResourceModel) RulesHaveChanges(
	ctx context.Context, plan FirewallResourceModel, diags *diag.Diagnostics,
) bool {
	oldInbound, newDiags := types.ListValueFrom(ctx, RuleObjectType, state.Inbound)
	diags.Append(newDiags...)

	oldOutbound, newDiags := types.ListValueFrom(ctx, RuleObjectType, state.Outbound)
	diags.Append(newDiags...)

	newOutbound, newDiags := types.ListValueFrom(ctx, RuleObjectType, plan.Inbound)
	diags.Append(newDiags...)

	newInbound, newDiags := types.ListValueFrom(ctx, RuleObjectType, plan.Outbound)
	diags.Append(newDiags...)

	if newDiags.HasError() {
		diags.Append(newDiags...)
	}

	return !oldInbound.Equal(newInbound) || oldOutbound.Equal(newOutbound)
}

func (state *FirewallResourceModel) LinodesOrNodeBalancersHaveChanges(
	ctx context.Context, plan FirewallResourceModel,
) bool {
	return !state.Linodes.Equal(plan.Linodes) || state.NodeBalancers.Equal(plan.NodeBalancers)
}

func FlattenFirewallRules(
	ctx context.Context,
	rules []linodego.FirewallRule,
	knownRules []RuleModel,
	preserveKnown bool,
) ([]RuleModel, diag.Diagnostics) {
	if preserveKnown && knownRules == nil {
		return make([]RuleModel, 0), nil
	}

	if !preserveKnown {
		knownRules = make([]RuleModel, len(rules))
	}

	for i := range knownRules {
		if i > len(rules) {
			break
		}

		knownRules[i].Action = helper.KeepOrUpdateString(knownRules[i].Action, rules[i].Action, preserveKnown)
		knownRules[i].Label = helper.KeepOrUpdateString(knownRules[i].Label, rules[i].Label, preserveKnown)
		knownRules[i].Protocol = helper.KeepOrUpdateString(knownRules[i].Protocol, string(rules[i].Protocol), preserveKnown)
		knownRules[i].Ports = helper.KeepOrUpdateString(knownRules[i].Ports, rules[i].Ports, preserveKnown)

		ipv4, diags := types.ListValueFrom(ctx, types.StringType, rules[i].Addresses.IPv4)
		if diags.HasError() {
			return nil, diags
		}

		knownRules[i].IPv4 = helper.KeepOrUpdateValue(knownRules[i].IPv4, ipv4, preserveKnown)

		ipv6, diags := types.ListValueFrom(ctx, types.StringType, rules[i].Addresses.IPv6)
		if diags.HasError() {
			return nil, diags
		}

		knownRules[i].IPv6 = helper.KeepOrUpdateValue(knownRules[i].IPv6, ipv6, preserveKnown)
	}
	return knownRules, nil
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

func FlattenFirewallDevices(devices []linodego.FirewallDevice) []DeviceModel {
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
