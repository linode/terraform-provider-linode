//go:build unit

package nfsfilesystemaccesspolicy

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

func TestFlattenNFSFilesystemAccessPolicy(t *testing.T) {
	ctx := context.Background()
	created := time.Date(2026, time.May, 1, 14, 0, 0, 0, time.UTC)

	policy := &linodego.NFSFilesystemAccessPolicy{
		FilesystemID: 67890,
		Label:        "training-data-policy",
		Enabled:      true,
		Inherit:      true,
		SquashPolicy: linodego.NFSSquashPolicyRootSquash,
		Protocols:    []linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4},
		Status:       linodego.NFSAccessPolicyStatusActive,
		Created:      &created,
		LinodeACL: []linodego.NFSFilesystemAccessPolicyLinode{
			{
				ID:    12345678,
				Label: linodego.Pointer("ml-worker-01"),
				URL:   linodego.Pointer("/linode/instances/12345678"),
				IPv6:  linodego.Pointer("2400:8900::f03c:91ff:fe01:2345"),
			},
		},
	}

	var model ResourceModel
	var diags diag.Diagnostics
	model.Flatten(ctx, 12345, 67890, policy, false, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("12345,67890"), model.ID)
	assert.Equal(t, types.Int64Value(12345), model.SpaceID)
	assert.Equal(t, types.Int64Value(67890), model.FilesystemID)
	assert.Equal(t, types.StringValue("training-data-policy"), model.Label)
	assert.Equal(t, types.BoolValue(true), model.Enabled)
	assert.Equal(t, types.BoolValue(true), model.Inherit)
	assert.Equal(t, types.StringValue("root_squash"), model.SquashPolicy)
	assert.Equal(t, types.StringValue("active"), model.Status)
	assert.True(t, model.Updated.IsNull())

	var linodeIDs []int64
	require.False(t, model.LinodeIDs.ElementsAs(ctx, &linodeIDs, false).HasError())
	assert.Equal(t, []int64{12345678}, linodeIDs)

	var protocols []string
	require.False(t, model.Protocols.ElementsAs(ctx, &protocols, false).HasError())
	assert.Equal(t, []string{"nfsv4"}, protocols)

	var acl []LinodeACLModel
	require.False(t, model.LinodeACL.ElementsAs(ctx, &acl, false).HasError())
	require.Len(t, acl, 1)
	assert.Equal(t, types.StringValue("ml-worker-01"), acl[0].Label)
	assert.Equal(t, types.StringValue("2400:8900::f03c:91ff:fe01:2345"), acl[0].IPv6)
}

func TestFlattenNFSFilesystemAccessPolicyPreserveKnown(t *testing.T) {
	ctx := context.Background()

	policy := &linodego.NFSFilesystemAccessPolicy{
		Label:        "api-label",
		Enabled:      true,
		SquashPolicy: linodego.NFSSquashPolicyAllSquash,
		Status:       linodego.NFSAccessPolicyStatusActive,
	}

	model := ResourceModel{
		Label:        types.StringUnknown(),
		Enabled:      types.BoolUnknown(),
		SquashPolicy: types.StringUnknown(),
	}

	var diags diag.Diagnostics
	model.Flatten(ctx, 12345, 67890, policy, true, &diags)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, types.StringValue("api-label"), model.Label)
	assert.Equal(t, types.BoolValue(true), model.Enabled)
	assert.Equal(t, types.StringValue("all_squash"), model.SquashPolicy)
}

func TestNFSFilesystemAccessPolicyGetUpdateOptions(t *testing.T) {
	ctx := context.Background()

	linodeIDs, _ := types.SetValueFrom(ctx, types.Int64Type, []int64{111, 222})
	protocols, _ := types.SetValueFrom(ctx, types.StringType, []string{"nfsv4"})

	plan := ResourceModel{
		Enabled:      types.BoolValue(true),
		Inherit:      types.BoolValue(false),
		LinodeIDs:    linodeIDs,
		SquashPolicy: types.StringValue("root_squash"),
		Protocols:    protocols,
	}

	t.Run("create", func(t *testing.T) {
		var diags diag.Diagnostics
		opts, shouldUpdate := plan.GetUpdateOptions(ctx, nil, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.True(t, shouldUpdate)
		assert.Equal(t, linodego.Pointer(true), opts.Enabled)
		assert.Equal(t, linodego.Pointer(false), opts.Inherit)
		assert.Equal(t, linodego.Pointer(linodego.NFSSquashPolicyRootSquash), opts.SquashPolicy)
		require.NotNil(t, opts.LinodeIDs)
		assert.ElementsMatch(t, []int{111, 222}, *opts.LinodeIDs)
		require.NotNil(t, opts.Protocols)
		assert.Equal(t, []linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}, *opts.Protocols)
	})

	t.Run("no changes", func(t *testing.T) {
		state := plan
		var diags diag.Diagnostics
		_, shouldUpdate := plan.GetUpdateOptions(ctx, &state, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
		assert.False(t, shouldUpdate)
	})

	t.Run("clear linode ids", func(t *testing.T) {
		state := plan
		updated := plan
		updated.LinodeIDs, _ = types.SetValueFrom(ctx, types.Int64Type, []int64{})

		var diags diag.Diagnostics
		opts, shouldUpdate := updated.GetUpdateOptions(ctx, &state, &diags)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.True(t, shouldUpdate)
		assert.Nil(t, opts.Enabled)
		assert.Nil(t, opts.Inherit)
		require.NotNil(t, opts.LinodeIDs)
		assert.Empty(t, *opts.LinodeIDs)
	})
}
