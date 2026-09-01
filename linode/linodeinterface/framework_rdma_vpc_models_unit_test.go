//go:build unit

package linodeinterface

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestRDMAVPCAttrModelGetUpdateOptions(t *testing.T) {
	plan := createRDMAVPCAttrModel(t, 352276, "10.0.0.29")
	state := createRDMAVPCAttrModel(t, 352275, "10.0.0.22")

	var diags diag.Diagnostics
	opts, shouldUpdate := plan.GetUpdateOptions(context.Background(), state, &diags)

	require.False(t, diags.HasError())
	require.True(t, shouldUpdate)
	require.NotNil(t, opts.SubnetID)
	require.Equal(t, 352276, *opts.SubnetID)
	require.NotNil(t, opts.IPv4)
	require.Len(t, opts.IPv4.Addresses, 1)
	require.Equal(t, "10.0.0.29", opts.IPv4.Addresses[0].Address)
	require.NotNil(t, opts.IPv4.Addresses[0].Primary)
	require.True(t, *opts.IPv4.Addresses[0].Primary)
}

func TestRDMAVPCIPv4AttrModelGetUpdateOptions(t *testing.T) {
	plan := createRDMAVPCIPv4AttrModel(t, "10.0.0.29")
	state := createRDMAVPCIPv4AttrModel(t, "10.0.0.22")

	var diags diag.Diagnostics
	opts, shouldUpdate := plan.GetUpdateOptions(context.Background(), state, &diags)

	require.False(t, diags.HasError())
	require.True(t, shouldUpdate)
	require.Len(t, opts.Addresses, 1)
	require.Equal(t, "10.0.0.29", opts.Addresses[0].Address)
	require.NotNil(t, opts.Addresses[0].Primary)
	require.True(t, *opts.Addresses[0].Primary)
}

func createRDMAVPCAttrModel(t *testing.T, subnetID int64, address string) *RDMAVPCAttrModel {
	t.Helper()

	return &RDMAVPCAttrModel{
		SubnetID: types.Int64Value(subnetID),
		IPv4:     createRDMAVPCIPv4ObjectValue(t, address),
	}
}

func createRDMAVPCIPv4AttrModel(t *testing.T, address string) *RDMAVPCIPv4AttrModel {
	t.Helper()

	ctx := context.Background()
	listValue, listDiags := types.ListValueFrom(ctx, configuredRDMAVPCInterfaceIPv4Address.Type(), []RDMAVPCIPv4AddressAttrModel{
		{Address: types.StringValue(address), Primary: types.BoolValue(true)},
	})
	require.False(t, listDiags.HasError())

	return &RDMAVPCIPv4AttrModel{Addresses: listValue}
}

func createRDMAVPCIPv4ObjectValue(t *testing.T, address string) types.Object {
	t.Helper()

	ctx := context.Background()
	listValue, listDiags := types.ListValueFrom(ctx, configuredRDMAVPCInterfaceIPv4Address.Type(), []RDMAVPCIPv4AddressAttrModel{
		{Address: types.StringValue(address), Primary: types.BoolValue(true)},
	})
	require.False(t, listDiags.HasError())

	attrTypes := resourceRDMAVPCIPv4Attribute.GetType().(basetypes.ObjectType).AttrTypes
	objectValue, objectDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"addresses": listValue,
	})
	require.False(t, objectDiags.HasError())

	return objectValue
}
