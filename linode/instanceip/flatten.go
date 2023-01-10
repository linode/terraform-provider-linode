package instanceip

import (
	"github.com/linode/linodego"
)

func flattenIPv4(network *linodego.InstanceIPv4Response) map[string]interface{} {
	result := make(map[string]interface{})

	result["private"] = flattenIPs(network.Private)
	result["public"] = flattenIPs(network.Public)
	result["reserved"] = flattenIPs(network.Reserved)
	result["shared"] = flattenIPs(network.Shared)

	return result
}

func flattenIPv6(network *linodego.InstanceIPv6Response) map[string]interface{} {
	result := make(map[string]interface{})

	result["global"] = flattenGlobal(network.Global)
	result["link_local"] = flattenIP(network.LinkLocal)
	result["slaac"] = flattenIP(network.SLAAC)

	return result
}

func flattenIP(network *linodego.InstanceIP) map[string]interface{} {
	result := make(map[string]interface{})

	result["address"] = network.Address
	result["gateway"] = network.Gateway
	result["prefix"] = network.Prefix
	result["rdns"] = network.RDNS
	result["region"] = network.Region
	result["subnet_mask"] = network.SubnetMask
	result["type"] = network.Type

	return result
}

func flattenIPs(network []*linodego.InstanceIP) map[string]interface{} {
	result := make(map[string]interface{})

	for _, net := range network {
		result["address"] = net.Address
		result["gateway"] = net.Gateway
		result["prefix"] = net.Prefix
		result["rdns"] = net.RDNS
		result["region"] = net.Region
		result["subnet_mask"] = net.SubnetMask
		result["type"] = net.Type
	}

	return result
}

func flattenGlobal(network []linodego.IPv6Range) map[string]interface{} {
	result := make(map[string]interface{})

	for _, net := range network {
		result["prefix"] = net.Prefix
		result["range"] = net.Range
		result["region"] = net.Region
		result["route_target"] = net.RouteTarget
	}

	return result
}
