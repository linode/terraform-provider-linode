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
