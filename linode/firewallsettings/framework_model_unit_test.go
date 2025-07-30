//go:build unit

package firewallsettings_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/firewallsettings"
	"github.com/stretchr/testify/assert"
)

func TestFlattenFirewallSettings(t *testing.T) {
	ctx := context.Background()
	defaultFirewallIDsObjectAttrType := firewallsettings.FrameworkResourceSchema.Attributes["default_firewall_ids"].(schema.SingleNestedAttribute).GetType().(types.ObjectType).AttrTypes
	firewallSettings := linodego.FirewallSettings{
		DefaultFirewallIDs: linodego.DefaultFirewallIDs{
			Linode:          linodego.Pointer(123),
			NodeBalancer:    nil,
			PublicInterface: linodego.Pointer(789),
			VPCInterface:    nil,
		},
	}

	expectedModelWhenNotPreservingKnown := firewallsettings.FirewallSettingsModel{
		DefaultFirewallIDs: types.ObjectValueMust(
			defaultFirewallIDsObjectAttrType,
			map[string]attr.Value{
				"linode":           types.Int64Value(123),
				"nodebalancer":     types.Int64Null(),
				"public_interface": types.Int64Value(789),
				"vpc_interface":    types.Int64Null(),
			},
		),
	}

	tests := map[string]struct {
		model         firewallsettings.FirewallSettingsModel
		settings      linodego.FirewallSettings
		expected      firewallsettings.FirewallSettingsModel
		preserveKnown bool
	}{
		"unknown default firewall IDs with preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectUnknown(defaultFirewallIDsObjectAttrType),
			},
			settings:      firewallSettings,
			expected:      expectedModelWhenNotPreservingKnown,
			preserveKnown: true,
		},
		"null default firewall IDs with preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectNull(defaultFirewallIDsObjectAttrType),
			},
			settings: firewallSettings,
			expected: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectNull(defaultFirewallIDsObjectAttrType),
			},
			preserveKnown: true,
		},
		"known default firewall IDs with preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectValueMust(
					defaultFirewallIDsObjectAttrType,
					map[string]attr.Value{
						"linode":           types.Int64Value(123),
						"nodebalancer":     types.Int64Value(456),
						"public_interface": types.Int64Unknown(),
						"vpc_interface":    types.Int64Unknown(),
					},
				),
			},
			settings: firewallSettings,
			expected: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectValueMust(
					defaultFirewallIDsObjectAttrType,
					map[string]attr.Value{
						"linode":           types.Int64Value(123),
						"nodebalancer":     types.Int64Value(456),
						"public_interface": types.Int64Value(789),
						"vpc_interface":    types.Int64Null(),
					},
				),
			},
			preserveKnown: true,
		},
		"unknown default firewall IDs without preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectUnknown(defaultFirewallIDsObjectAttrType),
			},
			settings:      firewallSettings,
			expected:      expectedModelWhenNotPreservingKnown,
			preserveKnown: false,
		},
		"null default firewall IDs without preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectNull(defaultFirewallIDsObjectAttrType),
			},
			settings:      firewallSettings,
			expected:      expectedModelWhenNotPreservingKnown,
			preserveKnown: false,
		},
		"known default firewall IDs without preserving known": {
			model: firewallsettings.FirewallSettingsModel{
				DefaultFirewallIDs: types.ObjectValueMust(
					defaultFirewallIDsObjectAttrType,
					map[string]attr.Value{
						"linode":           types.Int64Value(123),
						"nodebalancer":     types.Int64Value(456),
						"public_interface": types.Int64Unknown(),
						"vpc_interface":    types.Int64Unknown(),
					},
				),
			},
			settings:      firewallSettings,
			expected:      expectedModelWhenNotPreservingKnown,
			preserveKnown: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			diags := &diag.Diagnostics{}
			tt.model.FlattenFirewallSettings(ctx, tt.settings, tt.preserveKnown, diags)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			assert.Equal(t, tt.expected.DefaultFirewallIDs, tt.model.DefaultFirewallIDs,
				"Flattened DefaultFirewallIDs should match expected value")
		})
	}
}
