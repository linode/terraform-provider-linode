//go:build unit

package lke

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandLKEClusterNetworkingCreateOptions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    map[string]any
		expected *linodego.LKEClusterNetworkingCreateOptions
	}{
		{
			name:     "nil block",
			input:    nil,
			expected: &linodego.LKEClusterNetworkingCreateOptions{},
		},
		{
			name: "unset fields are omitted",
			input: map[string]any{
				"routing_mode":                       nil,
				"cluster_cidr_ipv4":                  nil,
				"node_cidr_mask_size_ipv4":           nil,
				"service_cluster_ip_range_ipv4":      nil,
				"max_pods":                           nil,
				"nodebalancer_backend_prefix_length": nil,
			},
			expected: &linodego.LKEClusterNetworkingCreateOptions{},
		},
		{
			name: "explicit zero and empty values reach the API",
			input: map[string]any{
				"routing_mode":                       "vxlan",
				"cluster_cidr_ipv4":                  "",
				"node_cidr_mask_size_ipv4":           0,
				"service_cluster_ip_range_ipv4":      "",
				"max_pods":                           0,
				"nodebalancer_backend_prefix_length": 0,
			},
			expected: &linodego.LKEClusterNetworkingCreateOptions{
				RoutingMode:                     linodego.Pointer(linodego.LKEClusterRoutingModeVXLAN),
				ClusterCIDRIPv4:                 linodego.Pointer(""),
				NodeCIDRMaskSizeIPv4:            linodego.Pointer(0),
				ServiceClusterIPRangeIPv4:       linodego.Pointer(""),
				MaxPods:                         linodego.Pointer(0),
				NodeBalancerBackendPrefixLength: linodego.Pointer(0),
			},
		},
		{
			name: "vxlan with defaults",
			input: map[string]any{
				"routing_mode": "vxlan",
			},
			expected: &linodego.LKEClusterNetworkingCreateOptions{
				RoutingMode: linodego.Pointer(linodego.LKEClusterRoutingModeVXLAN),
			},
		},
		{
			name: "vxlan with custom network configuration",
			input: map[string]any{
				"routing_mode":                       "vxlan",
				"cluster_cidr_ipv4":                  "100.64.0.0/12",
				"node_cidr_mask_size_ipv4":           25,
				"service_cluster_ip_range_ipv4":      "100.96.0.0/22",
				"max_pods":                           64,
				"nodebalancer_backend_prefix_length": 22,
				// Computed-only fields must never be sent
				"cluster_cidr_ipv6":             "fd12:3456:789a:c85::/96",
				"node_cidr_mask_size_ipv6":      112,
				"service_cluster_ip_range_ipv6": "fd12:3456:789a:10ad::/112",
			},
			expected: &linodego.LKEClusterNetworkingCreateOptions{
				RoutingMode:                     linodego.Pointer(linodego.LKEClusterRoutingModeVXLAN),
				ClusterCIDRIPv4:                 linodego.Pointer("100.64.0.0/12"),
				NodeCIDRMaskSizeIPv4:            linodego.Pointer(25),
				ServiceClusterIPRangeIPv4:       linodego.Pointer("100.96.0.0/22"),
				MaxPods:                         linodego.Pointer(64),
				NodeBalancerBackendPrefixLength: linodego.Pointer(22),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := Resource()
			config := map[string]any{}
			if tc.input != nil {
				config["networking"] = []any{tc.input}
			}
			data := schema.TestResourceDataRaw(t, r.Schema, config)
			data.SetId("123")
			state := data.State()
			state.RawConfig = networkingTestConfig(t, r, config)
			data = r.Data(state)

			options := expandLKEClusterNetworkingCreateOptions(context.Background(), data)
			assert.Equal(t, tc.expected, options)
			expectedJSON, err := json.Marshal(tc.expected)
			require.NoError(t, err)
			actualJSON, err := json.Marshal(options)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))
		})
	}
}

func TestCustomDiffValidateNetworking(t *testing.T) {
	for _, tc := range []struct {
		name        string
		priorMode   string
		networking  map[string]any
		expectError bool
	}{
		{
			name:       "existing vxlan retains max pods when routing mode is omitted",
			priorMode:  "vxlan",
			networking: map[string]any{"max_pods": 64},
		},
		{
			name:       "existing vxlan accepts explicit null routing mode",
			priorMode:  "vxlan",
			networking: map[string]any{"routing_mode": nil, "max_pods": 64},
		},
		{
			name:        "new cluster with omitted routing mode defaults to native",
			networking:  map[string]any{"max_pods": 64},
			expectError: true,
		},
		{
			name:        "existing native cluster rejects custom fields",
			priorMode:   "native",
			networking:  map[string]any{"max_pods": 64},
			expectError: true,
		},
		{
			name:        "explicit native overrides prior vxlan",
			priorMode:   "vxlan",
			networking:  map[string]any{"routing_mode": "native", "max_pods": 64},
			expectError: true,
		},
		{
			name:       "new vxlan cluster accepts custom fields",
			networking: map[string]any{"routing_mode": "vxlan", "max_pods": 64},
		},
		{
			name:       "computed native fields are not treated as configured",
			priorMode:  "native",
			networking: map[string]any{},
		},
		{
			name:      "omitted networking block is allowed",
			priorMode: "vxlan",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"tier":       resourceSchema["tier"],
					"networking": resourceSchema["networking"],
				},
				CustomizeDiff: customDiffValidateNetworking,
			}
			config := map[string]any{"tier": TierEnterprise}
			if tc.networking != nil {
				config["networking"] = []any{tc.networking}
			}
			state := &terraform.InstanceState{RawConfig: networkingTestConfig(t, r, config)}
			if tc.priorMode != "" {
				state.ID = "123"
				state.Attributes = map[string]string{
					"tier":                      TierEnterprise,
					"networking.#":              "1",
					"networking.0.routing_mode": tc.priorMode,
					"networking.0.max_pods":     "64",
				}
			}
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if tc.expectError {
				require.ErrorContains(t, err, "networking.0.max_pods: can only be configured")
				return
			}
			require.NoError(t, err)
			if tc.priorMode != "" && diff != nil {
				assert.False(t, diff.RequiresNew(), "omitting computed networking fields must not replace the cluster")
			}
		})
	}
}

func networkingTestConfig(t *testing.T, r *schema.Resource, config map[string]any) cty.Value {
	t.Helper()
	raw, err := json.Marshal(config)
	require.NoError(t, err)
	value, err := ctyjson.Unmarshal(raw, r.CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	return value
}

func TestFlattenLKEClusterNetworking(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    *linodego.LKEClusterNetworking
		expected []map[string]any
	}{
		{
			name:     "no networking",
			input:    nil,
			expected: []map[string]any{},
		},
		{
			name: "ipv4 only",
			input: &linodego.LKEClusterNetworking{
				RoutingMode:                     linodego.LKEClusterRoutingModeNative,
				ClusterCIDRIPv4:                 "10.0.0.0/13",
				NodeCIDRMaskSizeIPv4:            24,
				ServiceClusterIPRangeIPv4:       "10.128.0.0/20",
				MaxPods:                         110,
				NodeBalancerBackendPrefixLength: 21,
			},
			expected: []map[string]any{
				{
					"routing_mode":                       "native",
					"cluster_cidr_ipv4":                  "10.0.0.0/13",
					"node_cidr_mask_size_ipv4":           24,
					"service_cluster_ip_range_ipv4":      "10.128.0.0/20",
					"max_pods":                           110,
					"nodebalancer_backend_prefix_length": 21,
				},
			},
		},
		{
			name: "dual stack",
			input: &linodego.LKEClusterNetworking{
				RoutingMode:                     linodego.LKEClusterRoutingModeVXLAN,
				ClusterCIDRIPv4:                 "100.64.0.0/11",
				NodeCIDRMaskSizeIPv4:            24,
				ServiceClusterIPRangeIPv4:       "100.96.0.0/20",
				MaxPods:                         110,
				NodeBalancerBackendPrefixLength: 21,
				ClusterCIDRIPv6:                 linodego.Pointer("fd12:3456:789a:c85::/96"),
				NodeCIDRMaskSizeIPv6:            linodego.Pointer(112),
				ServiceClusterIPRangeIPv6:       linodego.Pointer("fd12:3456:789a:10ad::/112"),
			},
			expected: []map[string]any{
				{
					"routing_mode":                       "vxlan",
					"cluster_cidr_ipv4":                  "100.64.0.0/11",
					"node_cidr_mask_size_ipv4":           24,
					"service_cluster_ip_range_ipv4":      "100.96.0.0/20",
					"max_pods":                           110,
					"nodebalancer_backend_prefix_length": 21,
					"cluster_cidr_ipv6":                  "fd12:3456:789a:c85::/96",
					"node_cidr_mask_size_ipv6":           112,
					"service_cluster_ip_range_ipv6":      "fd12:3456:789a:10ad::/112",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, flattenLKEClusterNetworking(tc.input))
		})
	}
}
