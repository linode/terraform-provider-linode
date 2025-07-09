//go:build unit

package firewall

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func TestExpandFirewallRules(t *testing.T) {
	testCases := []struct {
		name      string
		ruleSpecs []RuleModel
		expected  []linodego.FirewallRule
	}{
		{
			"Expand Firewall Rule Test 1",
			[]RuleModel{
				{
					Label:    types.StringValue("Rule 1"),
					Action:   types.StringValue("allow"),
					Protocol: types.StringValue("SSH"),
					Ports:    types.StringValue("22"),
					IPv4: types.ListValueMust(
						cidrtypes.IPv4PrefixType{},
						[]attr.Value{
							cidrtypes.NewIPv4PrefixValue("192.168.1.1/24"),
						},
					),
					IPv6: types.ListValueMust(
						cidrtypes.IPv6PrefixType{},
						[]attr.Value{},
					),
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
			var diags diag.Diagnostics
			result := ExpandFirewallRules(context.Background(), tc.ruleSpecs, &diags)
			assert.False(t, diags.HasError())

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
		expected []RuleModel
	}{
		{
			rules: []linodego.FirewallRule{
				rule1, rule2,
			},

			expected: []RuleModel{
				{
					Action: types.StringValue("allow"),
					Label:  types.StringValue("SSH"),
					IPv4: types.ListValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("192.168.0.2"),
						},
					),
					IPv6:     types.ListValueMust(types.StringType, []attr.Value{}),
					Ports:    types.StringValue("22"),
					Protocol: types.StringValue("TCP"),
				},
				{
					Action: types.StringValue("deny"),
					Label:  types.StringValue("Block ICMP"),
					IPv4: types.ListValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("192.168.0.0/24"),
						},
					),
					IPv6: types.ListValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("2001:db8::/64"),
						},
					),
					Ports:    types.StringNull(),
					Protocol: types.StringValue("ICMP"),
				},
			},
		},
	}

	for _, c := range cases {
		out, err := FlattenFirewallRules(context.Background(), c.rules, nil, false)
		if err != nil {
			t.Fatal(err)
		}
		for i, rule := range out {
			if i > len(c.expected) {
				t.Fatal("firewall rules length not matched")
			}
			if !rule.Action.Equal(c.expected[i].Action) ||
				!rule.Label.Equal(c.expected[i].Label) ||
				!rule.IPv4.Equal(c.expected[i].IPv4) ||
				!rule.IPv6.Equal(c.expected[i].IPv6) ||
				!rule.Ports.Equal(c.expected[i].Ports) ||
				!rule.Protocol.Equal(c.expected[i].Protocol) {
				t.Errorf("flatten result mismatches expected values, expected: %v, rule: %v", c.expected[i], rule)
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

	expected := []DeviceModel{
		{
			ID:       types.Int64Value(123),
			EntityID: types.Int64Value(1111),
			Type:     types.StringValue(string(linodego.FirewallDeviceLinode)),
			Label:    types.StringValue("device_entity_1"),
			URL:      types.StringValue("test-firewall.example.com"),
		},
		{
			ID:       types.Int64Value(1234),
			EntityID: types.Int64Value(2222),
			Type:     types.StringValue(string(linodego.FirewallDeviceLinode)),
			Label:    types.StringValue("device_entity_2"),
			URL:      types.StringValue("test-firewall.example-2.com"),
		},
		{
			ID:       types.Int64Value(12345),
			EntityID: types.Int64Value(3333),
			Type:     types.StringValue(string(linodego.FirewallDeviceNodeBalancer)),
			Label:    types.StringValue("device_entity_3"),
			URL:      types.StringValue("test-firewall.example-3.com"),
		},
	}

	result := FlattenFirewallDevices(devices)

	for i, r := range result {
		if !r.ID.Equal(expected[i].ID) ||
			!r.EntityID.Equal(expected[i].EntityID) ||
			!r.Type.Equal(expected[i].Type) ||
			!r.Label.Equal(expected[i].Label) ||
			!r.URL.Equal(expected[i].URL) {
			t.Errorf("flatten result mismatches expected values, expected: %v, result: %v", expected, result)
		}
	}
}
