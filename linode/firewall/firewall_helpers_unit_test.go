//go:build unit

package firewall

import (
	"context"
	"reflect"
	"testing"

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
						types.StringType,
						[]attr.Value{
							types.StringValue("192.168.1.1/24"),
						},
					),
					IPv6: types.ListValueMust(
						types.StringType,
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
		var diags diag.Diagnostics
		out := FlattenFirewallRules(context.Background(), c.rules, nil, false, &diags)
		if diags.HasError() {
			t.Fatal(diags.Errors())
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

func TestSeparateRulesetRefs_MixedRules(t *testing.T) {
	rules := []linodego.FirewallRule{
		{RuleSet: 4010},
		{
			Action:   "ACCEPT",
			Label:    "allow-ssh",
			Ports:    "22",
			Protocol: "TCP",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
				IPv6: &[]string{"::/0"},
			},
		},
		{RuleSet: 4011},
		{
			Action:   "ACCEPT",
			Label:    "allow-http",
			Ports:    "80",
			Protocol: "TCP",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"10.0.0.0/8"},
				IPv6: &[]string{},
			},
		},
	}

	rulesetIDs, inlineRules := separateRulesetRefs(rules)

	assert.Equal(t, []int64{4010, 4011}, rulesetIDs)
	assert.Len(t, inlineRules, 2)
	assert.Equal(t, "allow-ssh", inlineRules[0].Label)
	assert.Equal(t, "allow-http", inlineRules[1].Label)
}

func TestSeparateRulesetRefs_NoRulesets(t *testing.T) {
	rules := []linodego.FirewallRule{
		{
			Action:   "ACCEPT",
			Label:    "allow-all",
			Protocol: "TCP",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
				IPv6: &[]string{"::/0"},
			},
		},
	}

	rulesetIDs, inlineRules := separateRulesetRefs(rules)

	assert.Nil(t, rulesetIDs)
	assert.Len(t, inlineRules, 1)
	assert.Equal(t, "allow-all", inlineRules[0].Label)
}

func TestSeparateRulesetRefs_OnlyRulesets(t *testing.T) {
	rules := []linodego.FirewallRule{
		{RuleSet: 100},
		{RuleSet: 200},
	}

	rulesetIDs, inlineRules := separateRulesetRefs(rules)

	assert.Equal(t, []int64{100, 200}, rulesetIDs)
	assert.Nil(t, inlineRules)
}

func TestSeparateRulesetRefs_Empty(t *testing.T) {
	rulesetIDs, inlineRules := separateRulesetRefs(nil)

	assert.Nil(t, rulesetIDs)
	assert.Nil(t, inlineRules)
}

func TestExpandFirewallRuleSet_WithRulesets(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	data := &FirewallResourceModel{
		InboundRuleSet: types.ListValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(4010),
		}),
		OutboundRuleSet: types.ListValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(4011),
		}),
		Inbound: []RuleModel{
			{
				Label:    types.StringValue("allow-ssh"),
				Action:   types.StringValue("ACCEPT"),
				Protocol: types.StringValue("TCP"),
				Ports:    types.StringValue("22"),
				IPv4: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("0.0.0.0/0"),
				}),
				IPv6: types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
		Outbound: []RuleModel{
			{
				Label:    types.StringValue("outbound-tcp"),
				Action:   types.StringValue("ACCEPT"),
				Protocol: types.StringValue("TCP"),
				Ports:    types.StringValue("1-65535"),
				IPv4: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("pl::subnets:2010"),
				}),
				IPv6: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("pl::subnets:2010"),
				}),
			},
		},
		InboundPolicy:  types.StringValue("DROP"),
		OutboundPolicy: types.StringValue("DROP"),
	}

	result := data.ExpandFirewallRuleSet(ctx, &diags)
	assert.False(t, diags.HasError())

	// Inbound: 1 ruleset ref + 1 inline rule
	assert.Len(t, result.Inbound, 2)
	assert.Equal(t, 4010, result.Inbound[0].RuleSet)
	assert.Equal(t, "", result.Inbound[0].Label)
	assert.Equal(t, "allow-ssh", result.Inbound[1].Label)
	assert.Equal(t, 0, result.Inbound[1].RuleSet)

	// Outbound: 1 ruleset ref + 1 inline rule
	assert.Len(t, result.Outbound, 2)
	assert.Equal(t, 4011, result.Outbound[0].RuleSet)
	assert.Equal(t, "outbound-tcp", result.Outbound[1].Label)

	assert.Equal(t, "DROP", result.InboundPolicy)
	assert.Equal(t, "DROP", result.OutboundPolicy)
}

func TestExpandFirewallRuleSet_NoRulesets(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	data := &FirewallResourceModel{
		InboundRuleSet:  types.ListNull(types.Int64Type),
		OutboundRuleSet: types.ListNull(types.Int64Type),
		Inbound: []RuleModel{
			{
				Label:    types.StringValue("allow-ssh"),
				Action:   types.StringValue("ACCEPT"),
				Protocol: types.StringValue("TCP"),
				Ports:    types.StringValue("22"),
				IPv4: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("0.0.0.0/0"),
				}),
				IPv6: types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
		Outbound:       nil,
		InboundPolicy:  types.StringValue("ACCEPT"),
		OutboundPolicy: types.StringValue("DROP"),
	}

	result := data.ExpandFirewallRuleSet(ctx, &diags)
	assert.False(t, diags.HasError())

	// Only the inline rule, no ruleset refs
	assert.Len(t, result.Inbound, 1)
	assert.Equal(t, 0, result.Inbound[0].RuleSet)
	assert.Equal(t, "allow-ssh", result.Inbound[0].Label)
	assert.Nil(t, result.Outbound)
}

func TestFlattenFirewallRules_PrefixListStrings(t *testing.T) {
	// Validate that prefix list tokens survive flatten (not rejected by CIDR validation)
	rules := []linodego.FirewallRule{
		{
			Action:   "ACCEPT",
			Label:    "outbound-subnet-tcp",
			Ports:    "1-65535",
			Protocol: "TCP",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"pl::subnets:2010"},
				IPv6: &[]string{"pl::subnets:2010"},
			},
		},
		{
			Action:   "ACCEPT",
			Label:    "outbound-registry",
			Ports:    "443",
			Protocol: "TCP",
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"pl:system:ps:managed:container:registry"},
				IPv6: &[]string{"pl:system:ps:managed:container:registry"},
			},
		},
	}

	var diags diag.Diagnostics
	result := FlattenFirewallRules(context.Background(), rules, nil, false, &diags)
	assert.False(t, diags.HasError())
	assert.Len(t, result, 2)

	// First rule: subnet prefix list
	assert.Equal(t, "outbound-subnet-tcp", result[0].Label.ValueString())
	var ipv4Vals []string
	diags.Append(result[0].IPv4.ElementsAs(context.Background(), &ipv4Vals, false)...)
	assert.False(t, diags.HasError())
	assert.Equal(t, []string{"pl::subnets:2010"}, ipv4Vals)

	// Second rule: registry prefix list
	assert.Equal(t, "outbound-registry", result[1].Label.ValueString())
	var registryVals []string
	diags.Append(result[1].IPv4.ElementsAs(context.Background(), &registryVals, false)...)
	assert.False(t, diags.HasError())
	assert.Equal(t, []string{"pl:system:ps:managed:container:registry"}, registryVals)
}
