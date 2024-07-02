package nb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/firewall"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalls"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// NodeBalancerModel describes the Terraform resource data model to match the
// resource schema.
type NodeBalancerModel struct {
	ID                 types.String      `tfsdk:"id"`
	Label              types.String      `tfsdk:"label"`
	Region             types.String      `tfsdk:"region"`
	ClientConnThrottle types.Int64       `tfsdk:"client_conn_throttle"`
	FirewallID         types.Int64       `tfsdk:"firewall_id"`
	Hostname           types.String      `tfsdk:"hostname"`
	IPv4               types.String      `tfsdk:"ipv4"`
	IPv6               types.String      `tfsdk:"ipv6"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Updated            timetypes.RFC3339 `tfsdk:"updated"`
	Transfer           types.List        `tfsdk:"transfer"`
	Tags               types.Set         `tfsdk:"tags"`
	Firewalls          types.List        `tfsdk:"firewalls"`
}

type FirewallModel struct {
	ID             types.Int64          `tfsdk:"id"`
	Label          types.String         `tfsdk:"label"`
	Tags           types.Set            `tfsdk:"tags"`
	InboundPolicy  types.String         `tfsdk:"inbound_policy"`
	OutboundPolicy types.String         `tfsdk:"outbound_policy"`
	Status         types.String         `tfsdk:"status"`
	Created        types.String         `tfsdk:"created"`
	Updated        types.String         `tfsdk:"updated"`
	Inbound        []firewall.RuleModel `tfsdk:"inbound"`
	Outbound       []firewall.RuleModel `tfsdk:"outbound"`
}

type nbModelV0 struct {
	ID                 types.String `tfsdk:"id"`
	Label              types.String `tfsdk:"label"`
	Region             types.String `tfsdk:"region"`
	ClientConnThrottle types.Int64  `tfsdk:"client_conn_throttle"`
	Hostname           types.String `tfsdk:"hostname"`
	IPv4               types.String `tfsdk:"ipv4"`
	IPv6               types.String `tfsdk:"ipv6"`
	Created            types.String `tfsdk:"created"`
	Updated            types.String `tfsdk:"updated"`
	Tags               types.Set    `tfsdk:"tags"`
	Transfer           types.Map    `tfsdk:"transfer"`
}

func (data *NodeBalancerModel) FlattenNodeBalancer(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
	firewalls []linodego.Firewall,
	preserveKnown bool,
) diag.Diagnostics {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(nodebalancer.ID), preserveKnown)
	data.Label = helper.KeepOrUpdateStringPointer(data.Label, nodebalancer.Label, preserveKnown)

	tags, diags := types.SetValueFrom(ctx, types.StringType, helper.StringSliceToFramework(nodebalancer.Tags))
	if diags.HasError() {
		return diags
	}
	data.Tags = helper.KeepOrUpdateValue(data.Tags, tags, preserveKnown)

	data.Region = helper.KeepOrUpdateString(data.Region, nodebalancer.Region, preserveKnown)
	data.ClientConnThrottle = helper.KeepOrUpdateInt64(
		data.ClientConnThrottle, int64(nodebalancer.ClientConnThrottle), preserveKnown,
	)
	data.Hostname = helper.KeepOrUpdateStringPointer(data.Hostname, nodebalancer.Hostname, preserveKnown)
	data.IPv4 = helper.KeepOrUpdateStringPointer(data.IPv4, nodebalancer.IPv4, preserveKnown)
	data.IPv6 = helper.KeepOrUpdateStringPointer(data.IPv6, nodebalancer.IPv6, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(nodebalancer.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(nodebalancer.Updated), preserveKnown,
	)

	transfer, diags := FlattenTransfer(ctx, nodebalancer.Transfer)
	if diags.HasError() {
		return diags
	}

	data.Transfer = helper.KeepOrUpdateValue(data.Transfer, *transfer, preserveKnown)

	fws, diags := parseNBFirewalls(ctx, firewalls)
	if diags.HasError() {
		return diags
	}

	data.Firewalls = helper.KeepOrUpdateValue(data.Firewalls, *fws, preserveKnown)

	return nil
}

func (data *NodeBalancerModel) CopyFrom(other NodeBalancerModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.Region = helper.KeepOrUpdateValue(data.Region, other.Region, preserveKnown)
	data.ClientConnThrottle = helper.KeepOrUpdateValue(
		data.ClientConnThrottle, other.ClientConnThrottle, preserveKnown,
	)
	data.FirewallID = helper.KeepOrUpdateValue(data.FirewallID, other.FirewallID, preserveKnown)
	data.Hostname = helper.KeepOrUpdateValue(data.Hostname, other.Hostname, preserveKnown)
	data.IPv4 = helper.KeepOrUpdateValue(data.IPv4, other.IPv4, preserveKnown)
	data.IPv6 = helper.KeepOrUpdateValue(data.IPv6, other.IPv6, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, other.Updated, preserveKnown)
	data.Transfer = helper.KeepOrUpdateValue(data.Transfer, other.Transfer, preserveKnown)
	data.Tags = helper.KeepOrUpdateValue(data.Tags, other.Tags, preserveKnown)
	data.Firewalls = helper.KeepOrUpdateValue(data.Firewalls, other.Firewalls, preserveKnown)
}

func parseNBFirewalls(
	ctx context.Context,
	firewalls []linodego.Firewall,
) (*types.List, diag.Diagnostics) {
	nbFirewalls := make([]FirewallModel, len(firewalls))

	for i, fw := range firewalls {
		nbFirewalls[i].ID = types.Int64Value(int64(fw.ID))
		nbFirewalls[i].Label = types.StringValue(fw.Label)
		nbFirewalls[i].InboundPolicy = types.StringValue(fw.Rules.InboundPolicy)
		nbFirewalls[i].OutboundPolicy = types.StringValue(fw.Rules.OutboundPolicy)
		nbFirewalls[i].Status = types.StringValue(string(fw.Status))
		nbFirewalls[i].Created = timetypes.NewRFC3339TimePointerValue(fw.Created).StringValue
		nbFirewalls[i].Updated = timetypes.NewRFC3339TimePointerValue(fw.Updated).StringValue

		tags, diags := types.SetValueFrom(ctx, types.StringType, fw.Tags)
		if diags.HasError() {
			return nil, diags
		}
		nbFirewalls[i].Tags = tags

		if fw.Rules.Inbound != nil {
			inBound, diags := firewall.FlattenFirewallRules(ctx, fw.Rules.Inbound, nbFirewalls[i].Inbound, false)
			if diags.HasError() {
				return nil, diags
			}
			nbFirewalls[i].Inbound = inBound
		}

		if fw.Rules.Outbound != nil {
			outBound, diags := firewall.FlattenFirewallRules(ctx, fw.Rules.Outbound, nbFirewalls[i].Inbound, false)
			if diags.HasError() {
				return nil, diags
			}
			nbFirewalls[i].Outbound = outBound
		}
	}

	result, diags := types.ListValueFrom(ctx, firewallObjType, nbFirewalls)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}

func FlattenTransfer(
	ctx context.Context,
	transfer linodego.NodeBalancerTransfer,
) (*types.List, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["in"] = helper.Float64PointerValueWithDefault(transfer.In)
	result["out"] = helper.Float64PointerValueWithDefault(transfer.Out)
	result["total"] = helper.Float64PointerValueWithDefault(transfer.Total)

	transferObj, diags := types.ObjectValue(TransferObjectType.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}

	resultList, diags := types.ListValueFrom(ctx, TransferObjectType, []attr.Value{transferObj})

	if diags.HasError() {
		return nil, diags
	}

	return &resultList, nil
}

func UpgradeResourceStateValue(val string) (types.Float64, diag.Diagnostic) {
	if val == "" {
		return types.Float64Value(0), nil
	}
	result, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return types.Float64Null(), diag.NewErrorDiagnostic(
			"Failed to upgrade state.",
			err.Error(),
		)
	}
	return types.Float64Value(result), nil
}

// NodeBalancerDataSourceModel describes the Terraform data model to match the
// data source schema.
type NodeBalancerDataSourceModel struct {
	ID                 types.Int64       `tfsdk:"id"`
	Label              types.String      `tfsdk:"label"`
	Region             types.String      `tfsdk:"region"`
	ClientConnThrottle types.Int64       `tfsdk:"client_conn_throttle"`
	Hostname           types.String      `tfsdk:"hostname"`
	IPv4               types.String      `tfsdk:"ipv4"`
	IPv6               types.String      `tfsdk:"ipv6"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Updated            timetypes.RFC3339 `tfsdk:"updated"`
	Transfer           types.List        `tfsdk:"transfer"`
	Tags               types.Set         `tfsdk:"tags"`
	Firewalls          []NBFirewallModel `tfsdk:"firewalls"`
}

type NBFirewallModel struct {
	ID             types.Int64                   `tfsdk:"id"`
	Label          types.String                  `tfsdk:"label"`
	Tags           []types.String                `tfsdk:"tags"`
	InboundPolicy  types.String                  `tfsdk:"inbound_policy"`
	OutboundPolicy types.String                  `tfsdk:"outbound_policy"`
	Status         types.String                  `tfsdk:"status"`
	Created        timetypes.RFC3339             `tfsdk:"created"`
	Updated        timetypes.RFC3339             `tfsdk:"updated"`
	Inbound        []firewalls.FirewallRuleModel `tfsdk:"inbound"`
	Outbound       []firewalls.FirewallRuleModel `tfsdk:"outbound"`
}

func (data *NodeBalancerDataSourceModel) flattenNodeBalancer(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
	firewalls []linodego.Firewall,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Label = types.StringPointerValue(nodebalancer.Label)
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.IPv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.IPv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = timetypes.NewRFC3339TimePointerValue(nodebalancer.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(nodebalancer.Updated)

	transfer, diags := FlattenTransfer(ctx, nodebalancer.Transfer)
	if diags.HasError() {
		return diags
	}
	data.Transfer = *transfer

	tags, diags := types.SetValueFrom(
		ctx,
		types.StringType,
		helper.StringSliceToFramework(nodebalancer.Tags),
	)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	nbFirewalls := make([]NBFirewallModel, len(firewalls))

	for i := range firewalls {
		var nbFirewall NBFirewallModel
		nbFirewall.FlattenFirewall(&firewalls[i], false)
		nbFirewalls[i] = nbFirewall
	}

	data.Firewalls = nbFirewalls

	return nil
}

func (d *NBFirewallModel) FlattenFirewall(firewall *linodego.Firewall, preserveKnown bool) {
	d.ID = types.Int64Value(int64(firewall.ID))
	d.Label = types.StringValue(firewall.Label)
	d.Tags = helper.StringSliceToFramework(firewall.Tags)
	d.InboundPolicy = types.StringValue(firewall.Rules.InboundPolicy)
	d.OutboundPolicy = types.StringValue(firewall.Rules.OutboundPolicy)
	d.Status = types.StringValue(string(firewall.Status))
	d.Created = timetypes.NewRFC3339TimePointerValue(firewall.Created)
	d.Updated = timetypes.NewRFC3339TimePointerValue(firewall.Updated)

	d.Inbound = make([]firewalls.FirewallRuleModel, len(firewall.Rules.Inbound))
	for i, v := range firewall.Rules.Inbound {
		var rule firewalls.FirewallRuleModel
		rule.ParseRule(v)
		d.Inbound[i] = rule
	}

	d.Outbound = make([]firewalls.FirewallRuleModel, len(firewall.Rules.Outbound))
	for i, v := range firewall.Rules.Outbound {
		var rule firewalls.FirewallRuleModel
		rule.ParseRule(v)
		d.Outbound[i] = rule
	}
}
