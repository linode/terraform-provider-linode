//go:build unit

package reservedip

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testReservedIP = linodego.InstanceIP{
	Address:    "198.51.100.5",
	Gateway:    "198.51.100.1",
	SubnetMask: "255.255.255.0",
	Prefix:     24,
	Type:       "ipv4",
	Public:     true,
	RDNS:       "198-51-100-5.ip.linodeusercontent.com",
	LinodeID:   0,
	Region:     "us-east",
	Reserved:   true,
	Tags:       []string{"tf-test", "reserved"},
}

func TestFlattenReservedIP(t *testing.T) {
	ctx := context.Background()
	m := &ResourceModel{}

	diags := m.flatten(ctx, testReservedIP, false)
	require.False(t, diags.HasError(), diags.Errors())

	assert.Equal(t, "198.51.100.5", m.ID.ValueString())
	assert.Equal(t, "198.51.100.5", m.Address.ValueString())
	assert.Equal(t, "198.51.100.1", m.Gateway.ValueString())
	assert.Equal(t, "255.255.255.0", m.SubnetMask.ValueString())
	assert.Equal(t, int64(24), m.Prefix.ValueInt64())
	assert.Equal(t, "ipv4", m.Type.ValueString())
	assert.Equal(t, true, m.Public.ValueBool())
	assert.Equal(t, "198-51-100-5.ip.linodeusercontent.com", m.RDNS.ValueString())
	assert.Equal(t, int64(0), m.LinodeID.ValueInt64())
	assert.Equal(t, "us-east", m.Region.ValueString())
	assert.Equal(t, true, m.Reserved.ValueBool())

	// Tags should be a set with two elements
	assert.False(t, m.Tags.IsNull())
	assert.Equal(t, 2, len(m.Tags.Elements()))

	// vpc_nat_1_1 and assigned_entity are nil → null list
	assert.True(t, m.VPCNAT1To1.IsNull())
	assert.True(t, m.AssignedEntity.IsNull())
}

func TestFlattenReservedIP_NilTags_NullState(t *testing.T) {
	ctx := context.Background()
	ip := testReservedIP
	ip.Tags = nil

	// Model starts with null Tags (fresh create, no prior state)
	m := &ResourceModel{}

	diags := m.flatten(ctx, ip, false)
	require.False(t, diags.HasError(), diags.Errors())

	// Should default to an empty set, not null
	assert.False(t, m.Tags.IsNull(), "expected empty set, got null")
	assert.Equal(t, 0, len(m.Tags.Elements()))
}

func TestFlattenReservedIP_NilTags_PreservesKnownState(t *testing.T) {
	ctx := context.Background()
	ip := testReservedIP
	ip.Tags = nil

	// Model already has known tags from prior state
	existing, diags := types.SetValueFrom(ctx, types.StringType, []string{"existing-tag"})
	require.False(t, diags.HasError())
	m := &ResourceModel{Tags: existing}

	diags = m.flatten(ctx, ip, false)
	require.False(t, diags.HasError(), diags.Errors())

	// Nil API response should not clobber the known state value
	assert.Equal(t, existing, m.Tags)
}

func TestFlattenReservedIP_WithVPCNAT(t *testing.T) {
	ctx := context.Background()
	ip := testReservedIP
	ip.Tags = nil
	ip.VPCNAT1To1 = &linodego.InstanceIPNAT1To1{
		Address:  "10.0.0.5",
		SubnetID: 42,
		VPCID:    7,
	}

	m := &ResourceModel{}
	diags := m.flatten(ctx, ip, false)
	require.False(t, diags.HasError(), diags.Errors())

	assert.False(t, m.VPCNAT1To1.IsNull())
	assert.Equal(t, 1, len(m.VPCNAT1To1.Elements()))
}

func TestFlattenReservedIP_WithAssignedEntity(t *testing.T) {
	ctx := context.Background()
	ip := testReservedIP
	ip.Tags = nil
	ip.AssignedEntity = &linodego.ReservedIPAssignedEntity{
		ID:    12345,
		Label: "my-linode",
		Type:  "linode",
		URL:   "/v4/linode/instances/12345",
	}

	m := &ResourceModel{}
	diags := m.flatten(ctx, ip, false)
	require.False(t, diags.HasError(), diags.Errors())

	assert.False(t, m.AssignedEntity.IsNull())
	assert.Equal(t, 1, len(m.AssignedEntity.Elements()))

	// Verify the entity fields inside the list
	var entities []AssignedEntityModel
	diags = m.AssignedEntity.ElementsAs(ctx, &entities, false)
	require.False(t, diags.HasError())
	require.Len(t, entities, 1)
	assert.Equal(t, int64(12345), entities[0].ID.ValueInt64())
	assert.Equal(t, "my-linode", entities[0].Label.ValueString())
	assert.Equal(t, "linode", entities[0].Type.ValueString())
	assert.Equal(t, "/v4/linode/instances/12345", entities[0].URL.ValueString())
}
