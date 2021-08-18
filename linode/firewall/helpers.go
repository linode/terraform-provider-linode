package firewall

import (
	"strings"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func expandFirewallStatus(disabled interface{}) linodego.FirewallStatus {
	return map[bool]linodego.FirewallStatus{
		true:  linodego.FirewallDisabled,
		false: linodego.FirewallEnabled,
	}[disabled.(bool)]
}

func expandFirewallRules(ruleSpecs []interface{}) []linodego.FirewallRule {
	rules := make([]linodego.FirewallRule, len(ruleSpecs))
	for i, ruleSpec := range ruleSpecs {
		ruleSpec := ruleSpec.(map[string]interface{})
		rule := linodego.FirewallRule{}

		rule.Label = ruleSpec["label"].(string)
		rule.Action = ruleSpec["action"].(string)
		rule.Protocol = linodego.NetworkProtocol(strings.ToUpper(ruleSpec["protocol"].(string)))
		rule.Ports = ruleSpec["ports"].(string)

		ipv4 := helper.ExpandStringList(ruleSpec["ipv4"].([]interface{}))
		if len(ipv4) > 0 {
			rule.Addresses.IPv4 = &ipv4
		}
		ipv6 := helper.ExpandStringList(ruleSpec["ipv6"].([]interface{}))
		if len(ipv6) > 0 {
			rule.Addresses.IPv6 = &ipv6
		}
		rules[i] = rule
	}
	return rules
}

func flattenFirewallLinodes(devices []linodego.FirewallDevice) []int {
	linodes := make([]int, 0, len(devices))
	for _, device := range devices {
		if device.Entity.Type == linodego.FirewallDeviceLinode {
			linodes = append(linodes, device.Entity.ID)
		}
	}
	return linodes
}

func flattenFirewallRules(rules []linodego.FirewallRule) []map[string]interface{} {
	specs := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		specs[i] = map[string]interface{}{
			"label":    rule.Label,
			"action":   rule.Action,
			"protocol": rule.Protocol,
			"ports":    rule.Ports,
			"ipv4":     rule.Addresses.IPv4,
			"ipv6":     rule.Addresses.IPv6,
		}
	}
	return specs
}

func flattenFirewallDevices(devices []linodego.FirewallDevice) []map[string]interface{} {
	governedDevices := make([]map[string]interface{}, len(devices))
	for i, device := range devices {
		governedDevices[i] = map[string]interface{}{
			"id":        device.ID,
			"entity_id": device.Entity.ID,
			"type":      device.Entity.Type,
			"label":     device.Entity.Label,
			"url":       device.Entity.URL,
		}
	}
	return governedDevices
}
