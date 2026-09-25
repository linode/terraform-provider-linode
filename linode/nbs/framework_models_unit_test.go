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
	model := &NodeBalancerModel{}
	diags := model.flattenNodeBalancer(t.Context(), &linodego.NodeBalancer{
		Type:                linodego.NBTypeCommon,
		BackendConnectivity: &connectivity,
	})
	require.False(t, diags.HasError())
	require.Equal(t, types.StringValue("common"), model.Type)
	require.Equal(t, types.StringValue("undefined"), model.BackendConnectivity)

	diags = model.flattenNodeBalancer(t.Context(), &linodego.NodeBalancer{})
	require.False(t, diags.HasError())
	require.True(t, model.BackendConnectivity.IsNull())
}
