package firewall

import (
	"context"

	"github.com/linode/terraform-provider-linode/linode/helper"

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
	Created        types.String `tfsdk:"created"`
	Updated        types.String `tfsdk:"updated"`
}

func (data *FirewallModel) parseComputedAttributes(
	ctx context.Context,
	firewall *linodego.Firewall,
	rules *linodego.FirewallRuleSet,
	devices []linodego.FirewallDevice,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(firewall.ID))
	data.Status = types.StringValue(string(firewall.Status))
	data.Created = types.StringValue(firewall.Created.Format(helper.TIME_FORMAT))
	data.Updated = types.StringValue(firewall.Updated.Format(helper.TIME_FORMAT))

	linodes, diags := types.SetValueFrom(ctx, types.Int64Type, AggregateLinodeIDs(devices))
	if diags.HasError() {
		return diags
	}
	data.Linodes = linodes

	firewallDevices, diags := parseFirewallDevices(devices)
	if diags.HasError() {
		return diags
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
	tags, diags := types.SetValueFrom(ctx, types.StringType, firewall.Tags)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags
	data.Disabled = types.BoolValue(firewall.Status == linodego.FirewallDisabled)
	data.InboundPolicy = types.StringValue(rules.InboundPolicy)
	data.OutboundPolicy = types.StringValue(rules.OutboundPolicy)
	data.Label = types.StringValue(firewall.Label)

	if rules.Inbound != nil {
		inBound, diags := parseFirewallRules(ctx, rules.Inbound)
		if diags.HasError() {
			return diags
		}
		data.Inbound = *inBound
	}

	if rules.Outbound != nil {
		outBound, diags := parseFirewallRules(ctx, rules.Outbound)
		if diags.HasError() {
			return diags
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
		ipv4, diags := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv4)
		if diags.HasError() {
			return nil, diags
		}
		spec["ipv4"] = ipv4

		ipv6, diags := types.ListValueFrom(ctx, types.StringType, rule.Addresses.IPv6)
		if diags.HasError() {
			return nil, diags
		}
		spec["ipv6"] = ipv6

		obj, diags := types.ObjectValue(datasourceRuleObjectType.AttrTypes, spec)
		if diags.HasError() {
			return nil, diags
		}

		specs[i] = obj
	}

	result, diags := basetypes.NewListValue(
		datasourceRuleObjectType,
		specs,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}

func AggregateLinodeIDs(devices []linodego.FirewallDevice) []int {
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

		obj, diags := types.ObjectValue(deviceObjectType.AttrTypes, governedDevice)
		if diags.HasError() {
			return nil, diags
		}

		governedDevices[i] = obj
	}

	result, diags := basetypes.NewListValue(
		deviceObjectType,
		governedDevices,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}
