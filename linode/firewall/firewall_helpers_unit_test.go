//go:build unit

package firewall

import (
	"reflect"
	"testing"

	"github.com/linode/linodego"
)

// Assertion helpers
func compareRule(rule map[string]interface{}, expected map[string]interface{}) bool {
	return rule["action"] == expected["action"] &&
		rule["label"] == expected["label"] &&
		rule["ports"] == expected["ports"] &&
		rule["protocol"] == expected["protocol"] &&
		reflect.DeepEqual(*rule["ipv4"].(*[]string), *expected["ipv4"].(*[]string)) &&
		reflect.DeepEqual(*rule["ipv6"].(*[]string), *expected["ipv6"].(*[]string))
}

// Unit tests for private functions in helper
// Functions under test: expandFirewallStatus, expandFirewallRules, flattenFirewallDeviceIDs, flattenFirewallRules, flattenFirewallDevices

func TestExpandFirewallStatus(t *testing.T) {
	testCases := []struct {
		name     string
		disabled interface{}
		expected linodego.FirewallStatus
	}{
		{"Firewall Enabled Test", false, linodego.FirewallEnabled},
		{"Firewall Disabled Test", true, linodego.FirewallDisabled},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandFirewallStatus(tc.disabled.(bool))
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestFlattenFirewallDevices(t *testing.T) {
	deviceEntity1 := linodego.FirewallDeviceEntity{
		ID:    1111,
		Type:  linodego.FirewallDeviceLinode,
		Label: "device_entity_1",
		URL:   "test-firewall.example.com",
	}

	deviceEntity2 := linodego.FirewallDeviceEntity{
		ID:    2222,
		Type:  linodego.FirewallDeviceLinode,
		Label: "device_entity_2",
		URL:   "test-firewall.example-2.com",
	}

	deviceEntity3 := linodego.FirewallDeviceEntity{
		ID:    3333,
		Type:  linodego.FirewallDeviceNodeBalancer,
		Label: "device_entity_3",
		URL:   "test-firewall.example-3.com",
	}

	devices := []linodego.FirewallDevice{
		{
			ID:     123,
			Entity: deviceEntity1,
		},
		{
			ID:     1234,
			Entity: deviceEntity2,
		},
		{
			ID:     12345,
			Entity: deviceEntity3,
		},
	}

	expected := []map[string]interface{}{
		{
			"id":        123,
			"entity_id": 1111,
			"type":      linodego.FirewallDeviceLinode,
			"label":     "device_entity_1",
			"url":       "test-firewall.example.com",
		},
		{
			"id":        1234,
			"entity_id": 2222,
			"type":      linodego.FirewallDeviceLinode,
			"label":     "device_entity_2",
			"url":       "test-firewall.example-2.com",
		},
		{
			"id":        12345,
			"entity_id": 3333,
			"type":      linodego.FirewallDeviceNodeBalancer,
			"label":     "device_entity_3",
			"url":       "test-firewall.example-3.com",
		},
	}

	result := flattenFirewallDevices(devices)

	for i, r := range result {
		for key, value := range expected[i] {
			if r[key] != value {
				t.Errorf("Mismatched value for key '%s' at index %d. Expected: %v, Got: %v", key, i, value, r[key])
			}
		}
	}
}
