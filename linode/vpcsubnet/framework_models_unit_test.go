//go:build unit

package vpcsubnet

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestFlattenSubnet_basic(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	subnet := &linodego.VPCSubnet{
		ID:      100,
		Label:   "test-subnet",
		IPv4:    "10.0.0.0/24",
		VPCType: linodego.VPCTypeRegular,
		Linodes: []linodego.VPCSubnetLinode{},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "100", model.ID.ValueString())
	require.Equal(t, "test-subnet", model.Label.ValueString())
	require.Equal(t, "10.0.0.0/24", model.IPv4.ValueString())
	require.Equal(t, "regular", model.VPCType.ValueString())
	require.False(t, model.Created.IsNull())
	require.False(t, model.Updated.IsNull())
}

func TestFlattenSubnet_rdmaType(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	subnet := &linodego.VPCSubnet{
		ID:      200,
		Label:   "rdma-subnet",
		IPv4:    "192.168.0.0/24",
		VPCType: linodego.VPCTypeRDMA,
		Linodes: []linodego.VPCSubnetLinode{},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "200", model.ID.ValueString())
	require.Equal(t, "rdma-subnet", model.Label.ValueString())
	require.Equal(t, "192.168.0.0/24", model.IPv4.ValueString())
	require.Equal(t, "rdma", model.VPCType.ValueString())
}

func TestFlattenSubnet_withLinodes(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	configID := 101

	subnet := &linodego.VPCSubnet{
		ID:      300,
		Label:   "subnet-with-linodes",
		IPv4:    "10.0.0.0/24",
		VPCType: linodego.VPCTypeRegular,
		Linodes: []linodego.VPCSubnetLinode{
			{
				ID: 1001,
				Interfaces: []linodego.VPCSubnetLinodeInterface{
					{
						ID:       5001,
						Active:   true,
						ConfigID: &configID,
					},
				},
			},
		},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "300", model.ID.ValueString())
	require.Equal(t, "subnet-with-linodes", model.Label.ValueString())
	require.Equal(t, "regular", model.VPCType.ValueString())

	// Verify linodes list
	require.False(t, model.Linodes.IsNull())

	linodeObjs := model.Linodes.Elements()
	require.Len(t, linodeObjs, 1)
}

func TestFlattenSubnet_withDatabases(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	ipv4Range := "10.0.0.4/30"

	subnet := &linodego.VPCSubnet{
		ID:      400,
		Label:   "subnet-with-db",
		IPv4:    "10.0.0.0/24",
		VPCType: linodego.VPCTypeRegular,
		Linodes: []linodego.VPCSubnetLinode{},
		Databases: []linodego.VPCSubnetDatabase{
			{
				ID:         2001,
				IPv4Range:  &ipv4Range,
				IPv6Ranges: []string{"fd00::1/64"},
			},
		},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "400", model.ID.ValueString())
	require.Equal(t, "subnet-with-db", model.Label.ValueString())

	// Verify databases list
	require.False(t, model.Databases.IsNull())

	dbElements := model.Databases.Elements()
	require.Len(t, dbElements, 1)
}

func TestFlattenSubnet_withNodebalancers(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	subnet := &linodego.VPCSubnet{
		ID:      500,
		Label:   "subnet-with-nb",
		IPv4:    "10.0.0.0/24",
		VPCType: linodego.VPCTypeRegular,
		Linodes: []linodego.VPCSubnetLinode{},
		Nodebalancers: []linodego.VPCSubnetNodebalancers{
			{
				ID:        3001,
				Ipv4Range: "10.0.0.8/30",
				Ipv6Ranges: []linodego.VPCSubnetNodebalancersRanges{
					{Range: "fd00::2/64"},
				},
			},
		},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "500", model.ID.ValueString())
	require.Equal(t, "subnet-with-nb", model.Label.ValueString())

	// Verify nodebalancers list
	require.False(t, model.Nodebalancers.IsNull())

	nbElements := model.Nodebalancers.Elements()
	require.Len(t, nbElements, 1)
}

func TestFlattenSubnet_withIPv6(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	subnet := &linodego.VPCSubnet{
		ID:      600,
		Label:   "subnet-with-ipv6",
		IPv4:    "10.0.0.0/24",
		VPCType: linodego.VPCTypeRegular,
		Linodes: []linodego.VPCSubnetLinode{},
		IPv6: []linodego.VPCIPv6Range{
			{Range: "fd00:abcd::/64"},
		},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "600", model.ID.ValueString())
	require.Equal(t, "subnet-with-ipv6", model.Label.ValueString())
	require.Equal(t, "regular", model.VPCType.ValueString())

	// Verify IPv6 is populated
	require.False(t, model.IPv6.IsNull())
	require.False(t, model.IPv6.IsUnknown())

	ipv6Models := make([]ResourceModelIPv6, 0)
	diags = model.IPv6.ElementsAs(context.Background(), &ipv6Models, false)
	require.False(t, diags.HasError(), diags.Errors())
	require.Len(t, ipv6Models, 1)
	require.Equal(t, "fd00:abcd::/64", ipv6Models[0].Range.ValueString())
	require.Equal(t, "fd00:abcd::/64", ipv6Models[0].AllocatedRange.ValueString())
}

func TestFlattenSubnet_DataSource(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	subnet := &linodego.VPCSubnet{
		ID:      700,
		Label:   "ds-subnet",
		IPv4:    "172.16.0.0/24",
		VPCType: linodego.VPCTypeRDMA,
		Linodes: []linodego.VPCSubnetLinode{
			{
				ID:         1001,
				Interfaces: []linodego.VPCSubnetLinodeInterface{},
			},
		},
		Databases: []linodego.VPCSubnetDatabase{
			{
				ID:         2001,
				IPv4Range:  nil,
				IPv6Ranges: []string{},
			},
		},
		Nodebalancers: []linodego.VPCSubnetNodebalancers{
			{
				ID:         3001,
				Ipv4Range:  "172.16.0.4/30",
				Ipv6Ranges: []linodego.VPCSubnetNodebalancersRanges{},
			},
		},
		IPv6:    []linodego.VPCIPv6Range{},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &DataSourceModel{}
	diags := model.FlattenSubnet(context.Background(), subnet, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "700", model.ID.ValueString())
	require.Equal(t, "ds-subnet", model.Label.ValueString())
	require.Equal(t, "172.16.0.0/24", model.IPv4.ValueString())
	require.Equal(t, "rdma", model.VPCType.ValueString())

	// Verify linodes
	require.False(t, model.Linodes.IsNull())
	linodeElements := model.Linodes.Elements()
	require.Len(t, linodeElements, 1)

	// Verify databases
	require.False(t, model.Databases.IsNull())
	dbElements := model.Databases.Elements()
	require.Len(t, dbElements, 1)

	// Verify nodebalancers
	require.False(t, model.Nodebalancers.IsNull())
	nbElements := model.Nodebalancers.Elements()
	require.Len(t, nbElements, 1)
}
