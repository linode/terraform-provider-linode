package instancenetworking

import (
	"github.com/linode/linodego"
)

func flattenIPv4(network *linodego.InstanceIPv4Response) []map[string]any {
	result := make(map[string]any)

	result["private"] = flattenIPs(network.Private)
	result["public"] = flattenIPs(network.Public)
	result["reserved"] = flattenIPs(network.Reserved)
	result["shared"] = flattenIPs(network.Shared)

	return []map[string]any{result}
}

func flattenIPv6(network *linodego.InstanceIPv6Response) []map[string]any {
	result := make(map[string]any)

	result["global"] = flattenGlobal(network.Global)
	result["link_local"] = []any{flattenIP(network.LinkLocal)}
	result["slaac"] = []any{flattenIP(network.SLAAC)}

	return []map[string]any{result}
}

func flattenIP(network *linodego.InstanceIP) map[string]any {
	result := make(map[string]any)

	result["address"] = network.Address
	result["gateway"] = network.Gateway
	result["prefix"] = network.Prefix
	result["rdns"] = network.RDNS
	result["region"] = network.Region
	result["subnet_mask"] = network.SubnetMask
	result["type"] = network.Type

	return result
}

func flattenIPs(network []*linodego.InstanceIP) []map[string]any {
	result := make([]map[string]any, len(network))

	for i, net := range network {
		result[i] = flattenIP(net)
	}

	return result
}

func flattenGlobal(network []linodego.IPv6Range) []any {
	result := make(map[string]any)

	for _, net := range network {
		result["prefix"] = net.Prefix
		result["range"] = net.Range
		result["region"] = net.Region
		result["route_target"] = net.RouteTarget
	}

	return []any{result}
}
