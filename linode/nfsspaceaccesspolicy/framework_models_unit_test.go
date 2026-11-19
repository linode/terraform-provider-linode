//go:build unit

package nfsspaceaccesspolicy

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenNFSSpaceAccessPolicy(t *testing.T) {
	ctx := context.Background()
	cert := "-----BEGIN CERTIFICATE-----"
	created := time.Date(2026, time.May, 1, 14, 0, 0, 0, time.UTC)

	policy := &linodego.NFSSpaceAccessPolicy{
		SpaceID:    12345,
		Label:      "prod-space-policy",
		Enabled:    true,
		MTLSCACert: &cert,
		MTLSMode:   linodego.NFSMTLSModeOptional,
		Status:     linodego.NFSAccessPolicyStatusActive,
		Created:    &created,
		VPCACL: []linodego.NFSSpaceAccessPolicyVPC{
			{
				ID:    123,
				Label: linodego.Pointer("prod-vpc"),
				URL:   linodego.Pointer("/vpcs/123"),
				Subnets: []linodego.NFSSpaceAccessPolicyVPCSubnet{
					{ID: 456, Label: linodego.Pointer("app"), Range: linodego.Pointer("2400:8900:0:1::/64")},
				},
			},
			{
				ID:    345,
				Label: linodego.Pointer("staging-vpc"),
				Range: linodego.Pointer("2400:8900:1::/48"),
			},
		},
	}

	var model ResourceModel
	var diags diag.Diagnostics
	model.Flatten(ctx, 12345, policy, false, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("12345"), model.ID)
	assert.Equal(t, types.Int64Value(12345), model.SpaceID)
	assert.Equal(t, types.StringValue("prod-space-policy"), model.Label)
	assert.Equal(t, types.BoolValue(true), model.Enabled)
	assert.Equal(t, types.StringValue(cert), model.MTLSCACert)
	assert.Equal(t, types.StringValue("optional"), model.MTLSMode)
	assert.Equal(t, types.StringValue("active"), model.Status)
	assert.True(t, model.Updated.IsNull())

	var vpcs []VPCModel
	require.False(t, model.VPCs.ElementsAs(ctx, &vpcs, false).HasError())
	require.Len(t, vpcs, 2)

	vpcsByID := map[int64]VPCModel{}
	for _, v := range vpcs {
		vpcsByID[v.ID.ValueInt64()] = v
	}
	assert.Equal(t, 1, len(vpcsByID[123].Subnets.Elements()))
	assert.True(t, vpcsByID[345].Subnets.IsNull())

	var acl []VPCACLModel
	require.False(t, model.VPCACL.ElementsAs(ctx, &acl, false).HasError())
	require.Len(t, acl, 2)
	assert.Equal(t, types.StringValue("prod-vpc"), acl[0].Label)
	assert.True(t, acl[0].Range.IsNull())
	assert.Equal(t, 1, len(acl[0].Subnets.Elements()))
	assert.Equal(t, types.StringValue("2400:8900:1::/48"), acl[1].Range)
}

func TestNFSSpaceAccessPolicyGetUpdateOptions(t *testing.T) {
	ctx := context.Background()

	vpcSet := func(id int64, subnets ...int64) types.Set {
		subnetValue := types.SetNull(types.Int64Type)
		if len(subnets) > 0 {
			subnetValue, _ = types.SetValueFrom(ctx, types.Int64Type, subnets)
		}
		set, _ := types.SetValueFrom(ctx, vpcObjectType, []VPCModel{
			{ID: types.Int64Value(id), Subnets: subnetValue},
		})
		return set
	}

	plan := ResourceModel{
		Enabled:    types.BoolValue(true),
		VPCs:       vpcSet(123, 456, 789),
		MTLSCACert: types.StringUnknown(),
		MTLSMode:   types.StringValue("optional"),
	}

	t.Run("create", func(t *testing.T) {
		var diags diag.Diagnostics
		opts, shouldUpdate := plan.GetUpdateOptions(ctx, nil, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.True(t, shouldUpdate)
		assert.Equal(t, linodego.Pointer(true), opts.Enabled)
		assert.Nil(t, opts.MTLSCACert)
		assert.Equal(t, linodego.Pointer(linodego.NFSMTLSModeOptional), opts.MTLSMode)
		require.NotNil(t, opts.VPCs)
		require.Len(t, *opts.VPCs, 1)
		assert.Equal(t, 123, (*opts.VPCs)[0].ID)
		assert.ElementsMatch(t, []int{456, 789}, (*opts.VPCs)[0].Subnets)
	})

	t.Run("no changes", func(t *testing.T) {
		state := plan
		var diags diag.Diagnostics
		_, shouldUpdate := plan.GetUpdateOptions(ctx, &state, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
		assert.False(t, shouldUpdate)
	})

	t.Run("vpc change only", func(t *testing.T) {
		state := plan
		state.VPCs = vpcSet(123)
		var diags diag.Diagnostics
		opts, shouldUpdate := plan.GetUpdateOptions(ctx, &state, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.True(t, shouldUpdate)
		assert.Nil(t, opts.Enabled)
		assert.Nil(t, opts.MTLSMode)
		require.NotNil(t, opts.VPCs)
	})
}

func TestFlattenNFSSpaceAccessPolicyPreserveKnown(t *testing.T) {
	ctx := context.Background()

	policy := &linodego.NFSSpaceAccessPolicy{
		Label:    "api-label",
		Enabled:  true,
		MTLSMode: linodego.NFSMTLSModeDisabled,
		Status:   linodego.NFSAccessPolicyStatusActive,
	}

	model := ResourceModel{
		SpaceID:  types.Int64Value(12345),
		ID:       types.StringUnknown(),
		Label:    types.StringUnknown(),
		Enabled:  types.BoolValue(false),
		VPCs:     types.SetUnknown(vpcObjectType),
		MTLSMode: types.StringUnknown(),
		VPCACL:   types.ListUnknown(vpcACLObjectType),
	}

	var diags diag.Diagnostics
	model.Flatten(ctx, 12345, policy, true, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("12345"), model.ID)
	assert.Equal(t, types.StringValue("api-label"), model.Label)
	assert.Equal(t, types.BoolValue(false), model.Enabled)
	assert.Equal(t, types.StringValue("disabled"), model.MTLSMode)
	assert.False(t, model.VPCs.IsUnknown())
	assert.Empty(t, model.VPCs.Elements())
	assert.False(t, model.VPCACL.IsUnknown())
}
