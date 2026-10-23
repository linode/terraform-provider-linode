//go:build unit

package lke

import (
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
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
				"routing_mode":                       "",
				"cluster_cidr_ipv4":                  "",
				"node_cidr_mask_size_ipv4":           0,
				"service_cluster_ip_range_ipv4":      "",
				"max_pods":                           0,
				"nodebalancer_backend_prefix_length": 0,
			},
			expected: &linodego.LKEClusterNetworkingCreateOptions{},
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
			assert.Equal(t, tc.expected, expandLKEClusterNetworkingCreateOptions(tc.input))
		})
	}
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
