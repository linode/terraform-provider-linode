//go:build unit

package placementgroup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestFlattenPGModel(t *testing.T) {
	label := "test-pgmodel"

	pg := linodego.PlacementGroup{
		ID:                   123,
		Label:                "test-pgmodel",
		Region:               "us-mia",
		PlacementGroupType:   linodego.PlacementGroupTypeAntiAffinityLocal,
		PlacementGroupPolicy: "flexible",
		IsCompliant:          false,
		Members: []linodego.PlacementGroupMember{
			{
				LinodeID:    123,
				IsCompliant: false,
			},
			{
				LinodeID:    456,
				IsCompliant: true,
			},
		},
	}

	model := &PlacementGroupResourceModel{
		ID:    types.StringValue("123"),
		Label: types.StringValue(label),
	}

	model.FlattenPlacementGroup(context.Background(), &pg, false)

	require.Equal(t, "123", model.ID.ValueString())
	require.Equal(t, label, model.Label.ValueString())
	require.Equal(t, "us-mia", model.Region.ValueString())
	require.Equal(t, string(linodego.PlacementGroupTypeAntiAffinityLocal), model.PlacementGroupType.ValueString())
	require.Equal(t, "flexible", model.PlacementGroupPolicy.ValueString())

	require.Equal(t, false, model.IsCompliant.ValueBool())

	members := make([]PlacementGroupMemberModel, 0)

	d := model.Members.ElementsAs(context.Background(), &members, false)
	require.False(t, d.HasError(), d.Errors())

	require.Equal(t, int64(123), members[0].LinodeID.ValueInt64())
	require.Equal(t, false, members[0].IsCompliant.ValueBool())

	require.Equal(t, int64(456), members[1].LinodeID.ValueInt64())
	require.Equal(t, true, members[1].IsCompliant.ValueBool())
}
