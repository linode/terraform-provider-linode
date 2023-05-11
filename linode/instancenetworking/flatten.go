package instancenetworking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

func flattenIPv4(ctx context.Context, network *linodego.InstanceIPv4Response) (*basetypes.ObjectValue, diag.Diagnostics) {
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

	return &obj, nil
}

func flattenIPv6(ctx context.Context, network *linodego.InstanceIPv6Response) (*basetypes.ObjectValue, diag.Diagnostics) {
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

	return &obj, nil
}

func flattenIP(ctx context.Context, network *linodego.InstanceIP) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["address"] = types.StringValue(network.Address)
	result["gateway"] = types.StringValue(network.Gateway)
	result["prefix"] = types.Int64Value(int64(network.Prefix))
	result["rdns"] = types.StringValue(network.RDNS)
	result["region"] = types.StringValue(network.Region)
	result["subnet_mask"] = types.StringValue(network.SubnetMask)
	result["type"] = types.StringValue(string(network.Type))

	obj, diag := types.ObjectValue(networkObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	return &obj, nil
}

func flattenIPs(ctx context.Context, network []*linodego.InstanceIP) (*basetypes.ListValue, diag.Diagnostics) {
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

func flattenIPV6Range(ctx context.Context, network *linodego.IPv6Range) (*basetypes.ObjectValue, diag.Diagnostics) {
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

func flattenGlobal(ctx context.Context, network []linodego.IPv6Range) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(network))

	for i, network := range network {
		result, diag := flattenIPV6Range(ctx, &network)
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
