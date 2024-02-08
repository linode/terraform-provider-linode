package firewall

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"golang.org/x/net/context"
)

// firewallDeviceAssignment is a helper struct intended to be used in conjunction
// with updateFirewallDevices.
type firewallDeviceAssignment struct {
	ID   int
	Type linodego.FirewallDeviceType
}

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

func updateFirewallDevices(
	ctx context.Context,
	d *schema.ResourceData,
	client linodego.Client,
	id int,
	configuredDevices []firewallDeviceAssignment,
) error {
	currentDevices, err := client.ListFirewallDevices(ctx, id, nil)
	if err != nil {
		return err
	}

	// Populate a map to track existing devices by assignment
	deviceMap := make(map[firewallDeviceAssignment]linodego.FirewallDevice)
	for _, device := range currentDevices {
		deviceMap[firewallDeviceAssignment{ID: device.Entity.ID, Type: device.Entity.Type}] = device
	}

	for _, device := range configuredDevices {
		if _, ok := deviceMap[device]; ok {
			// Device exists, drop it from the map so it won't be removed
			delete(deviceMap, device)
			continue
		}

		// Device doesn't exist, create a new one
		createOpts := linodego.FirewallDeviceCreateOptions{
			ID:   device.ID,
			Type: device.Type,
		}

		tflog.Debug(ctx, "client.CreateFirewallDevice(...)", map[string]any{
			"options": createOpts,
		})

		_, err := client.CreateFirewallDevice(ctx, id, createOpts)
		if err != nil {
			return err
		}
	}

	// Clean up remaining devices
	for _, device := range deviceMap {
		tflog.Debug(ctx, "client.DeleteFirewallDevice(...)", map[string]any{
			"device_id": device.ID,
		})

		if err := client.DeleteFirewallDevice(ctx, id, device.ID); err != nil {
			return err
		}
	}

	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"firewall_id": d.Id(),
	})
}
