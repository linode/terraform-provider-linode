package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

// FirewallModel describes the Terraform resource data model to match the
// resource schema.
type FirewallModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Label          types.String `tfsdk:"label"`
	Tags           types.Set    `tfsdk:"tags"`
	Disabled       types.Bool   `tfsdk:"disabled"`
	Inbound        types.List   `tfsdk:"inbound"`
	InboundPolicy  types.String `tfsdk:"inbound_policy"`
	Outbound       types.List   `tfsdk:"outbound"`
	OutboundPolicy types.String `tfsdk:"outbound_policy"`
	Linodes        types.Set    `tfsdk:"linodes"`
	Devices        types.List   `tfsdk:"devices"`
	Status         types.String `tfsdk:"status"`
}

func (data *FirewallModel) parseComputedAttributes(
	ctx context.Context,
	firewall *linodego.Firewall,
	rules *linodego.FirewallRuleSet,
	devices []linodego.FirewallDevice,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(firewall.ID))
	data.Disabled = types.BoolValue(firewall.Status == linodego.FirewallDisabled)
	data.Status = types.StringValue(string(firewall.Status))

	linodes, diag := types.SetValueFrom(ctx, types.Int64Type, parseFirewallLinodes(devices))
	if diag.HasError() {
		return diag
	}
	data.Linodes = linodes

	firewallDevices, diag := parseFirewallDevices(devices)
	if diag.HasError() {
		return diag
	}
	data.Devices = *firewallDevices

	return nil
}

func (data *FirewallModel) parseNonComputedAttributes(
	ctx context.Context,
	firewall *linodego.Firewall,
	rules *linodego.FirewallRuleSet,
	devices []linodego.FirewallDevice,
) diag.Diagnostics {
	tags, diag := types.SetValueFrom(ctx, types.StringType, firewall.Tags)
	if diag.HasError() {
		return diag
	}
	data.Tags = tags
	data.InboundPolicy = types.StringValue(rules.InboundPolicy)
	data.OutboundPolicy = types.StringValue(rules.OutboundPolicy)
	data.Label = types.StringValue(firewall.Label)

	if rules.Inbound != nil {
		inBound, diag := parseFirewallRules(ctx, rules.Inbound)
		if diag.HasError() {
			return diag
		}
		data.Inbound = *inBound
	}

	if rules.Outbound != nil {
		outBound, diag := parseFirewallRules(ctx, rules.Outbound)
		if diag.HasError() {
			return diag
		}
		data.Outbound = *outBound
	}

	return nil
}

func parseFirewallRules(
	ctx context.Context,
	rules []linodego.FirewallRule,
) (*basetypes.ListValue, diag.Diagnostics) {
	specs := make([]attr.Value, len(rules))

	for i, rule := range rules {
		spec := make(map[string]attr.Value)
		spec["action"] = types.StringValue(rule.Action)
		spec["label"] = types.StringValue(rule.Label)
		spec["protocol"] = types.StringValue(string(rule.Protocol))
		spec["ports"] = types.StringValue(rule.Ports)
		ipv4, diag := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv4)
		if diag.HasError() {
			return nil, diag
		}
		spec["ipv4"] = ipv4

		ipv6, diag := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv6)
		if diag.HasError() {
			return nil, diag
		}
		spec["ipv6"] = ipv6

		obj, diag := types.ObjectValue(datasourceRuleObjectType.AttrTypes, spec)
		if diag.HasError() {
			return nil, diag
		}

		specs[i] = obj
	}

	result, diag := basetypes.NewListValue(
		datasourceRuleObjectType,
		specs,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}

func parseFirewallLinodes(devices []linodego.FirewallDevice) []int {
	linodes := make([]int, 0, len(devices))
	for _, device := range devices {
		if device.Entity.Type == linodego.FirewallDeviceLinode {
			linodes = append(linodes, device.Entity.ID)
		}
	}
	return linodes
}

func parseFirewallDevices(devices []linodego.FirewallDevice) (*basetypes.ListValue, diag.Diagnostics) {
	governedDevices := make([]attr.Value, len(devices))
	for i, device := range devices {
		governedDevice := make(map[string]attr.Value)
		governedDevice["id"] = types.Int64Value(int64(device.ID))
		governedDevice["entity_id"] = types.Int64Value(int64(device.Entity.ID))
		governedDevice["type"] = types.StringValue(string(device.Entity.Type))
		governedDevice["label"] = types.StringValue(device.Entity.Label)
		governedDevice["url"] = types.StringValue(device.Entity.URL)

		obj, diag := types.ObjectValue(deviceObjectType.AttrTypes, governedDevice)
		if diag.HasError() {
			return nil, diag
		}

		governedDevices[i] = obj
	}

	result, diag := basetypes.NewListValue(
		deviceObjectType,
		governedDevices,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}
