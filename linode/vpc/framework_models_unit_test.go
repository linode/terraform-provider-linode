//go:build unit

package vpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFlattenVPC_rdmaType(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	vpc := &linodego.VPC{
		ID:          456,
		Label:       "rdma-vpc",
		Description: "An RDMA VPC",
		Region:      "us-ord",
		VPCType:     linodego.VPCTypeRDMA,
		Created:     &createdTime,
		Updated:     &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenVPC(context.Background(), vpc, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "456", model.ID.ValueString())
	require.Equal(t, "rdma-vpc", model.Label.ValueString())
	require.Equal(t, "An RDMA VPC", model.Description.ValueString())
	require.Equal(t, "us-ord", model.Region.ValueString())
	require.Equal(t, "rdma", model.VPCType.ValueString())
}

func TestFlattenVPC_withIPv6(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	vpc := &linodego.VPC{
		ID:          789,
		Label:       "ipv6-vpc",
		Description: "A VPC with IPv6",
		Region:      "us-east",
		VPCType:     linodego.VPCTypeRegular,
		IPv6: []linodego.VPCIPv6Range{
			{Range: "fd00::/52"},
		},
		Created: &createdTime,
		Updated: &updatedTime,
	}

	model := &ResourceModel{}
	diags := model.FlattenVPC(context.Background(), vpc, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "789", model.ID.ValueString())
	require.Equal(t, "ipv6-vpc", model.Label.ValueString())
	require.Equal(t, "regular", model.VPCType.ValueString())

	// Verify IPv6 is populated
	require.False(t, model.IPv6.IsNull())
	require.False(t, model.IPv6.IsUnknown())

	ipv6Models := make([]ResourceModelIPv6, 0)
	diags = model.IPv6.ElementsAs(context.Background(), &ipv6Models, false)
	require.False(t, diags.HasError(), diags.Errors())
	require.Len(t, ipv6Models, 1)
	require.Equal(t, "fd00::/52", ipv6Models[0].Range.ValueString())
	require.Equal(t, "fd00::/52", ipv6Models[0].AllocatedRange.ValueString())
}

func TestFlattenVPC_DataSource(t *testing.T) {
	createdTime := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, 6, 2, 12, 0, 0, 0, time.UTC)

	vpc := &linodego.VPC{
		ID:          123,
		Label:       "test-vpc",
		Description: "A test VPC",
		Region:      "us-east",
		VPCType:     linodego.VPCTypeRDMA,
		IPv6:        []linodego.VPCIPv6Range{},
		Created:     &createdTime,
		Updated:     &updatedTime,
	}

	model := &DataSourceModel{}
	diags := model.FlattenVPC(context.Background(), vpc, false)
	require.False(t, diags.HasError(), diags.Errors())

	require.Equal(t, "123", model.ID.ValueString())
	require.Equal(t, "test-vpc", model.Label.ValueString())
	require.Equal(t, "A test VPC", model.Description.ValueString())
	require.Equal(t, "us-east", model.Region.ValueString())
	require.Equal(t, "rdma", model.VPCType.ValueString())
}
