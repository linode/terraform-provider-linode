package firewalls

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	firewallresource "github.com/linode/terraform-provider-linode/linode/firewall"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

type FirewallRuleModel struct {
	Label    types.String   `tfsdk:"label"`
	Action   types.String   `tfsdk:"action"`
	Protocol types.String   `tfsdk:"protocol"`
	Ports    types.String   `tfsdk:"ports"`
	IPv4     []types.String `tfsdk:"ipv4"`
	IPv6     []types.String `tfsdk:"ipv6"`
}

func (data *FirewallRuleModel) parseRule(rule linodego.FirewallRule) {
	data.Label = types.StringValue(rule.Label)
	data.Action = types.StringValue(rule.Action)
	data.Protocol = types.StringValue(string(rule.Protocol))
	data.Ports = types.StringValue(rule.Ports)

	data.IPv4 = []types.String{}
	data.IPv6 = []types.String{}

	if rule.Addresses.IPv4 != nil {
		data.IPv4 = helper.StringSliceToFramework(*rule.Addresses.IPv4)
	}

	if rule.Addresses.IPv6 != nil {
		data.IPv6 = helper.StringSliceToFramework(*rule.Addresses.IPv6)
	}
}

type FirewallDeviceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	EntityID types.Int64  `tfsdk:"entity_id"`
	Type     types.String `tfsdk:"type"`
	Label    types.String `tfsdk:"label"`
	URL      types.String `tfsdk:"url"`
}

func (data *FirewallDeviceModel) parseDevice(device linodego.FirewallDevice) {
	data.ID = types.Int64Value(int64(device.ID))
	data.EntityID = types.Int64Value(int64(device.Entity.ID))
	data.Type = types.StringValue(string(device.Entity.Type))
	data.Label = types.StringValue(device.Entity.Label)
	data.URL = types.StringValue(device.Entity.URL)
}

// FirewallModel describes the Terraform resource data model to match the
// resource schema.
type FirewallModel struct {
	ID             types.Int64    `tfsdk:"id"`
	Label          types.String   `tfsdk:"label"`
	Tags           []types.String `tfsdk:"tags"`
	Disabled       types.Bool     `tfsdk:"disabled"`
	InboundPolicy  types.String   `tfsdk:"inbound_policy"`
	OutboundPolicy types.String   `tfsdk:"outbound_policy"`
	Linodes        []types.Int64  `tfsdk:"linodes"`
	Status         types.String   `tfsdk:"status"`
	Created        types.String   `tfsdk:"created"`
	Updated        types.String   `tfsdk:"updated"`

	Inbound  []FirewallRuleModel   `tfsdk:"inbound"`
	Outbound []FirewallRuleModel   `tfsdk:"outbound"`
	Devices  []FirewallDeviceModel `tfsdk:"devices"`
}

func (data *FirewallModel) parseFirewall(
	firewall linodego.Firewall,
	rules linodego.FirewallRuleSet,
	devices []linodego.FirewallDevice,
) {
	data.ID = types.Int64Value(int64(firewall.ID))
	data.Label = types.StringValue(firewall.Label)
	data.Tags = helper.StringSliceToFramework(firewall.Tags)
	data.Disabled = types.BoolValue(firewall.Status == linodego.FirewallDisabled)
	data.InboundPolicy = types.StringValue(rules.InboundPolicy)
	data.OutboundPolicy = types.StringValue(rules.OutboundPolicy)
	data.Linodes = helper.IntSliceToFramework(firewallresource.AggregateLinodeIDs(devices))
	data.Status = types.StringValue(string(firewall.Status))

	data.Created = types.StringValue(firewall.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(firewall.Updated.Format(time.RFC3339))

	data.Inbound = make([]FirewallRuleModel, len(rules.Inbound))
	for i, v := range rules.Inbound {
		var rule FirewallRuleModel
		rule.parseRule(v)
		data.Inbound[i] = rule
	}

	data.Outbound = make([]FirewallRuleModel, len(rules.Outbound))
	for i, v := range rules.Outbound {
		var rule FirewallRuleModel
		rule.parseRule(v)
		data.Outbound[i] = rule
	}

	data.Devices = make([]FirewallDeviceModel, len(devices))
	for i, v := range devices {
		var device FirewallDeviceModel
		device.parseDevice(v)
		data.Devices[i] = device
	}
}

type FirewallFilterModel struct {
	ID        types.String                     `tfsdk:"id"`
	Filters   frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order     types.String                     `tfsdk:"order"`
	OrderBy   types.String                     `tfsdk:"order_by"`
	Firewalls []FirewallModel                  `tfsdk:"firewalls"`
}

func (data *FirewallFilterModel) parseFirewalls(
	ctx context.Context,
	client *linodego.Client,
	firewalls []linodego.Firewall,
) diag.Diagnostics {
	var diags diag.Diagnostics

	result := make([]FirewallModel, len(firewalls))

	for i := range firewalls {
		var fwData FirewallModel

		fw := firewalls[i]

		devices, err := client.ListFirewallDevices(ctx, fw.ID, nil)
		if err != nil {
			diags.AddError("Failed to list Firewall devices", err.Error())
			return diags
		}

		rules, err := client.GetFirewallRules(ctx, fw.ID)
		if err != nil {
			diags.AddError("Failed to get Firewall rules", err.Error())
			return diags
		}

		fwData.parseFirewall(fw, *rules, devices)

		result[i] = fwData
	}

	data.Firewalls = result

	return diags
}
