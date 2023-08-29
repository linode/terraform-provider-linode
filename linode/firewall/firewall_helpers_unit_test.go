//go:build unit

package firewall

import (
	"reflect"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
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
// Functions under test: expandFirewallStatus, expandFirewallRules, flattenFirewallLinodes, flattenFirewallRules, flattenFirewallDevices

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
			result := expandFirewallStatus(tc.disabled)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestExpandFirewallRules(t *testing.T) {
	testCases := []struct {
		name      string
		ruleSpecs []interface{}
		expected  []linodego.FirewallRule
	}{
		{
			"Expand Firewall Rule Test 1",
			[]interface{}{
				map[string]interface{}{
					"label":    "Rule 1",
					"action":   "allow",
					"protocol": "SSH",
					"ports":    "22",
					"ipv4":     []interface{}{"192.168.1.1/24"},
					"ipv6":     []interface{}{},
				},
			},
			[]linodego.FirewallRule{
				{
					Action:      "allow",
					Label:       "Rule 1",
					Description: "Allow SSH connections",
					Ports:       "22",
					Protocol:    "SSH",
					Addresses: linodego.NetworkAddresses{
						IPv4: &[]string{"192.168.1.1/24"},
						IPv6: &[]string{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandFirewallRules(tc.ruleSpecs)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d rules, but got %d", len(tc.expected), len(result))
			}
			for i, expectedRule := range tc.expected {
				assert.Equal(t, result[i].Label, expectedRule.Label)
				assert.Equal(t, result[i].Action, expectedRule.Action)
				assert.Equal(t, result[i].Protocol, expectedRule.Protocol)
				assert.Equal(t, result[i].Addresses.IPv4, expectedRule.Addresses.IPv4)
			}
		})
	}
}

func TestFlattenFirewallRules(t *testing.T) {
	rule1 := linodego.FirewallRule{
		Action:      "allow",
		Label:       "SSH",
		Description: "Allow SSH connections",
		Ports:       "22",
		Protocol:    "TCP",
		Addresses: linodego.NetworkAddresses{
			IPv4: &[]string{"192.168.0.2"},
			IPv6: &[]string{},
		},
	}

	rule2 := linodego.FirewallRule{
		Action:      "deny",
		Label:       "Block ICMP",
		Description: "Block ICMP traffic",
		Ports:       "",
		Protocol:    "ICMP",
		Addresses: linodego.NetworkAddresses{
			IPv4: &[]string{"192.168.0.0/24"},
			IPv6: &[]string{"2001:db8::/64"},
		},
	}

	cases := []struct {
		rules    []linodego.FirewallRule
		expected []map[string]interface{}
	}{
		{
			rules: []linodego.FirewallRule{
				rule1, rule2,
			},

			expected: []map[string]interface{}{
				{
					"action":   "allow",
					"label":    "SSH",
					"ipv4":     &[]string{"192.168.0.2"},
					"ipv6":     &[]string{},
					"ports":    "22",
					"protocol": "TCP",
				},
				{
					"action":   "deny",
					"label":    "Block ICMP",
					"ipv4":     &[]string{"192.168.0.0/24"},
					"ipv6":     &[]string{"2001:db8::/64"},
					"ports":    "",
					"protocol": "ICMP",
				},
			},
		},
	}

	for _, c := range cases {
		out := flattenFirewallRules(c.rules)

		for i, rule := range out {
			if i < len(c.expected) {
				compareRule(rule, c.expected[i])
			} else {
				break
			}
		}
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

	devices := []linodego.FirewallDevice{
		{
			ID:     123,
			Entity: deviceEntity1,
		},
		{
			ID:     1234,
			Entity: deviceEntity2,
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
