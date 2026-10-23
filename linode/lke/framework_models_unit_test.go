//go:build unit

package lke

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseNetworking(t *testing.T) {
	t.Run("no networking", func(t *testing.T) {
		assert.Equal(t, []LKENetworking{}, parseNetworking(nil))
	})

	t.Run("ipv4 only", func(t *testing.T) {
		result := parseNetworking(&linodego.LKEClusterNetworking{
			RoutingMode:                     linodego.LKEClusterRoutingModeVXLAN,
			ClusterCIDRIPv4:                 "100.64.0.0/11",
			NodeCIDRMaskSizeIPv4:            24,
			ServiceClusterIPRangeIPv4:       "100.96.0.0/20",
			MaxPods:                         110,
			NodeBalancerBackendPrefixLength: 21,
		})

		assert.Equal(t, []LKENetworking{
			{
				RoutingMode:                     types.StringValue("vxlan"),
				ClusterCIDRIPv4:                 types.StringValue("100.64.0.0/11"),
				NodeCIDRMaskSizeIPv4:            types.Int64Value(24),
				ServiceClusterIPRangeIPv4:       types.StringValue("100.96.0.0/20"),
				MaxPods:                         types.Int64Value(110),
				NodeBalancerBackendPrefixLength: types.Int64Value(21),
				ClusterCIDRIPv6:                 types.StringNull(),
				NodeCIDRMaskSizeIPv6:            types.Int64Null(),
				ServiceClusterIPRangeIPv6:       types.StringNull(),
			},
		}, result)
	})

	t.Run("dual stack", func(t *testing.T) {
		result := parseNetworking(&linodego.LKEClusterNetworking{
			RoutingMode:                     linodego.LKEClusterRoutingModeVXLAN,
			ClusterCIDRIPv4:                 "100.64.0.0/11",
			NodeCIDRMaskSizeIPv4:            24,
			ServiceClusterIPRangeIPv4:       "100.96.0.0/20",
			MaxPods:                         110,
			NodeBalancerBackendPrefixLength: 21,
			ClusterCIDRIPv6:                 linodego.Pointer("fd12:3456:789a:c85::/96"),
			NodeCIDRMaskSizeIPv6:            linodego.Pointer(112),
			ServiceClusterIPRangeIPv6:       linodego.Pointer("fd12:3456:789a:10ad::/112"),
		})

		assert.Len(t, result, 1)
		assert.Equal(t, types.StringValue("fd12:3456:789a:c85::/96"), result[0].ClusterCIDRIPv6)
		assert.Equal(t, types.Int64Value(112), result[0].NodeCIDRMaskSizeIPv6)
		assert.Equal(t, types.StringValue("fd12:3456:789a:10ad::/112"), result[0].ServiceClusterIPRangeIPv6)
	})
}
