package nb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Ipv4               types.String      `tfsdk:"ipv4"`
	Ipv6               types.String      `tfsdk:"ipv6"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Updated            timetypes.RFC3339 `tfsdk:"updated"`
	Transfer           types.List        `tfsdk:"transfer"`
	Tags               types.Set         `tfsdk:"tags"`
	Firewalls          types.List        `tfsdk:"firewalls"`
}

type nbModelV0 struct {
	ID                 types.String `tfsdk:"id"`
	Label              types.String `tfsdk:"label"`
	Region             types.String `tfsdk:"region"`
	ClientConnThrottle types.Int64  `tfsdk:"client_conn_throttle"`
	Hostname           types.String `tfsdk:"hostname"`
	Ipv4               types.String `tfsdk:"ipv4"`
	Ipv6               types.String `tfsdk:"ipv6"`
	Created            types.String `tfsdk:"created"`
	Updated            types.String `tfsdk:"updated"`
	Tags               types.Set    `tfsdk:"tags"`
	Transfer           types.Map    `tfsdk:"transfer"`
}

func (data *NodeBalancerModel) ParseNonComputedAttrs(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.ID = types.StringValue(strconv.Itoa(nodebalancer.ID))
	data.Label = types.StringPointerValue(nodebalancer.Label)

	tags, diags := types.SetValueFrom(ctx, types.StringType, helper.StringSliceToFramework(nodebalancer.Tags))
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	return nil
}

func (data *NodeBalancerModel) ParseComputedAttrs(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
	firewalls []linodego.Firewall,
) diag.Diagnostics {
	data.ID = types.StringValue(strconv.Itoa(nodebalancer.ID))
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.Ipv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.Ipv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = timetypes.NewRFC3339TimePointerValue(nodebalancer.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(nodebalancer.Updated)

	transfer, diags := ParseTransfer(ctx, nodebalancer.Transfer)
	if diags.HasError() {
		return diags
	}

	data.Transfer = *transfer

	fws, diags := parseNBFirewalls(ctx, firewalls)
	if diags.HasError() {
		return diags
	}

	data.Firewalls = *fws

	return nil
}

func parseNBFirewalls(
	ctx context.Context,
	firewalls []linodego.Firewall,
) (*types.List, diag.Diagnostics) {
	nbFirewalls := make([]attr.Value, len(firewalls))

	for i, fw := range firewalls {
		f := make(map[string]attr.Value)
		f["id"] = types.Int64Value(int64(fw.ID))
		f["label"] = types.StringValue(fw.Label)
		f["inbound_policy"] = types.StringValue(fw.Rules.InboundPolicy)
		f["outbound_policy"] = types.StringValue(fw.Rules.OutboundPolicy)
		f["status"] = types.StringValue(string(fw.Status))
		f["created"] = timetypes.NewRFC3339TimePointerValue(fw.Created).StringValue
		f["updated"] = timetypes.NewRFC3339TimePointerValue(fw.Updated).StringValue

		tags, diags := types.SetValueFrom(ctx, types.StringType, fw.Tags)
		if diags.HasError() {
			return nil, diags
		}
		f["tags"] = tags

		if fw.Rules.Inbound != nil {
			inBound, diags := firewall.ParseFirewallRules(ctx, fw.Rules.Inbound)
			if diags.HasError() {
				return nil, diags
			}
			f["inbound"] = *inBound
		}

		if fw.Rules.Outbound != nil {
			outBound, diags := firewall.ParseFirewallRules(ctx, fw.Rules.Outbound)
			if diags.HasError() {
				return nil, diags
			}
			f["outbound"] = *outBound
		}

		obj, diags := types.ObjectValue(firewallObjType.AttrTypes, f)
		if diags.HasError() {
			return nil, diags
		}

		nbFirewalls[i] = obj
	}

	result, diags := types.ListValue(firewallObjType, nbFirewalls)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}

func ParseTransfer(
	ctx context.Context,
	transfer linodego.NodeBalancerTransfer,
) (*basetypes.ListValue, diag.Diagnostics) {
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

func UpgradeResourceStateValue(val string) (basetypes.Float64Value, diag.Diagnostic) {
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
	Ipv4               types.String      `tfsdk:"ipv4"`
	Ipv6               types.String      `tfsdk:"ipv6"`
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
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.Ipv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.Ipv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = timetypes.NewRFC3339TimePointerValue(nodebalancer.Created)
	data.Updated = timetypes.NewRFC3339TimePointerValue(nodebalancer.Updated)

	transfer, diags := ParseTransfer(ctx, nodebalancer.Transfer)
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
		nbFirewall.parseFirewall(&firewalls[i])
		nbFirewalls[i] = nbFirewall
	}

	data.Firewalls = nbFirewalls

	return nil
}

func (d *NBFirewallModel) parseFirewall(firewall *linodego.Firewall) {
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
