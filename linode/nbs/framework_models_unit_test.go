//go:build unit

package nbs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
)

func TestFlattenNodeBalancerConnectivity(t *testing.T) {
	connectivity := linodego.NBBackendConnectivityUndefined
	backendIPv6Prefix := "2600:3c22:1:20:0:3039::/96"
	model := &NodeBalancerModel{}
	diags := model.flattenNodeBalancer(t.Context(), &linodego.NodeBalancer{
		Type:                linodego.NBTypeCommon,
		BackendConnectivity: &connectivity,
		BackendIPv6Prefix:   &backendIPv6Prefix,
	})
	require.False(t, diags.HasError())
	require.Equal(t, types.StringValue("common"), model.Type)
	require.Equal(t, types.StringValue("undefined"), model.BackendConnectivity)
	require.Equal(t, types.StringValue(backendIPv6Prefix), model.BackendIPv6Prefix)

	diags = model.flattenNodeBalancer(t.Context(), &linodego.NodeBalancer{})
	require.False(t, diags.HasError())
	require.True(t, model.BackendConnectivity.IsNull())
	require.True(t, model.BackendIPv6Prefix.IsNull())
}
