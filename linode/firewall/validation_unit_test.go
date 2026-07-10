//go:build unit

package firewall

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirewallProtocolValidator(t *testing.T) {
	testCases := []struct {
		name      string
		protocol  string
		wantError bool
	}{
		{name: "zero", protocol: "0"},
		{name: "max", protocol: "255"},
		{name: "all keyword", protocol: "ALL"},
		{name: "tcp keyword", protocol: "TCP"},
		{name: "too large", protocol: "256", wantError: true},
		{name: "leading zero", protocol: "06", wantError: true},
		{name: "invalid keyword", protocol: "any", wantError: true},
		{name: "lowercase all", protocol: "all", wantError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp fwvalidator.StringResponse

			firewallProtocolValidator{}.ValidateString(
				context.Background(),
				fwvalidator.StringRequest{
					Path:        path.Root("protocol"),
					ConfigValue: types.StringValue(tc.protocol),
				},
				&resp,
			)

			assert.Equal(t, tc.wantError, resp.Diagnostics.HasError())
		})
	}
}

func TestExpandRuleOmitsNullPortsFromJSON(t *testing.T) {
	ruleModel := RuleModel{
		Label:       types.StringValue("tf-test-in"),
		Action:      types.StringValue("ACCEPT"),
		Protocol:    types.StringValue("ALL"),
		Ports:       types.StringNull(),
		Description: types.StringNull(),
		IPv4: types.ListValueMust(
			cidrtypes.IPv4PrefixType{},
			[]attr.Value{
				cidrtypes.NewIPv4PrefixValue("0.0.0.0/0"),
			},
		),
		IPv6: types.ListValueMust(
			cidrtypes.IPv6PrefixType{},
			[]attr.Value{},
		),
	}

	var diags diag.Diagnostics
	rule := ExpandRule[linodego.FirewallRuleInbound](context.Background(), ruleModel, &diags)

	require.False(t, diags.HasError())
	require.Empty(t, rule.Ports)

	payload, err := json.Marshal(rule)
	require.NoError(t, err)

	assert.Contains(t, string(payload), `"protocol":"ALL"`)
	assert.NotContains(t, string(payload), `"ports"`)
}
