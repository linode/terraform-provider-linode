package instancenetworking

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	LinodeID types.Int64  `tfsdk:"linode_id"`
	IPV4     types.List   `tfsdk:"ipv4"`
	IPV6     types.List   `tfsdk:"ipv6"`
	ID       types.String `tfsdk:"id"`
}

func (data *DataSourceModel) parseInstanceIPAddressResponse(
	ctx context.Context, ip *linodego.InstanceIPAddressResponse,
) diag.Diagnostics {
	ipv4, diags := flattenIPv4(ctx, ip.IPv4)
	if diags.HasError() {
		return diags
	}

	data.IPV4 = *ipv4

	ipv6, diags := flattenIPv6(ctx, ip.IPv6)
	if diags.HasError() {
		return diags
	}

	data.IPV6 = *ipv6

	id, err := json.Marshal(ip)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

func flattenIPv4(ctx context.Context, network *linodego.InstanceIPv4Response) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	private, diag := flattenIPs(ctx, network.Private)
	if diag.HasError() {
		return nil, diag
	}

	public, diag := flattenIPs(ctx, network.Public)
	if diag.HasError() {
		return nil, diag
	}

	reserved, diag := flattenIPs(ctx, network.Reserved)
	if diag.HasError() {
		return nil, diag
	}

	shared, diag := flattenIPs(ctx, network.Shared)
	if diag.HasError() {
		return nil, diag
	}

	result["private"] = private
	result["public"] = public
	result["reserved"] = reserved
	result["shared"] = shared

	obj, diag := types.ObjectValue(ipv4ObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		ipv4ObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func flattenIPv6(ctx context.Context, network *linodego.InstanceIPv6Response) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	global, diag := flattenGlobal(ctx, network.Global)
	if diag.HasError() {
		return nil, diag
	}

	link_local, diag := flattenIP(ctx, network.LinkLocal)
	if diag.HasError() {
		return nil, diag
	}

	slaac, diag := flattenIP(ctx, network.SLAAC)
	if diag.HasError() {
		return nil, diag
	}

	result["global"] = global
	result["link_local"] = link_local
	result["slaac"] = slaac

	obj, diag := types.ObjectValue(ipv6ObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		ipv6ObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func FlattenIPVPCNAT1To1(data *linodego.InstanceIPNAT1To1) (basetypes.ObjectValue, diag.Diagnostics) {
	if data == nil {
		return types.ObjectNull(VPCNAT1To1Type.AttrTypes), nil
	}

	result := map[string]attr.Value{
		"address":   types.StringValue(data.Address),
		"vpc_id":    types.Int64Value(int64(data.VPCID)),
		"subnet_id": types.Int64Value(int64(data.SubnetID)),
	}

	obj, diag := types.ObjectValue(VPCNAT1To1Type.AttrTypes, result)
	if diag.HasError() {
		return types.ObjectNull(VPCNAT1To1Type.AttrTypes), diag
	}

	return obj, nil
}

func flattenIP(ctx context.Context, network *linodego.InstanceIP) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	vpcNAT1To1, d := FlattenIPVPCNAT1To1(network.VPCNAT1To1)
	if d.HasError() {
		return nil, d
	}

	result["vpc_nat_1_1"] = vpcNAT1To1

	result["address"] = types.StringValue(network.Address)
	result["gateway"] = types.StringValue(network.Gateway)
	result["prefix"] = types.Int64Value(int64(network.Prefix))
	result["rdns"] = types.StringValue(network.RDNS)
	result["region"] = types.StringValue(network.Region)
	result["subnet_mask"] = types.StringValue(network.SubnetMask)
	result["type"] = types.StringValue(string(network.Type))
	result["public"] = types.BoolValue(network.Public)
	result["linode_id"] = types.Int64Value(int64(network.LinodeID))

	obj, d := types.ObjectValue(networkObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}

	return &obj, nil
}

func flattenIPs(ctx context.Context, network []*linodego.InstanceIP) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	resultList := make([]attr.Value, len(network))

	for i, network := range network {
		result, diag := flattenIP(ctx, network)
		if diag.HasError() {
			return nil, diag
		}

		resultList[i] = result
	}

	result, diag := basetypes.NewListValue(
		networkObjectType,
		resultList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}

func flattenIPV6Range(ctx context.Context, network linodego.IPv6Range) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["prefix"] = types.Int64Value(int64(network.Prefix))
	result["range"] = types.StringValue(network.Range)
	result["region"] = types.StringValue(network.Region)
	result["route_target"] = types.StringValue(network.RouteTarget)

	obj, diag := types.ObjectValue(globalObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	return &obj, nil
}

func flattenGlobal(ctx context.Context, network []linodego.IPv6Range) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	resultList := make([]attr.Value, len(network))

	for i, network := range network {
		result, diag := flattenIPV6Range(ctx, network)
		if diag.HasError() {
			return nil, diag
		}

		resultList[i] = result
	}

	result, diag := basetypes.NewListValue(
		globalObjectType,
		resultList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}
