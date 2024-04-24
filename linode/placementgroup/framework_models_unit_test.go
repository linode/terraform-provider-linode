//go:build unit

package placementgroup

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenPGModel(t *testing.T) {
	pg := linodego.PlacementGroup{
		ID:           123,
		Label:        "test",
		Region:       "us-mia",
		AffinityType: linodego.AffinityTypeAntiAffinityLocal,
		IsStrict:     false,
		IsCompliant:  false,
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

	model := PlacementGroupModel{}
	model.FlattenPlacementGroup(&pg, false)

	assert.Equal(t, "123", model.ID.ValueString())
	assert.Equal(t, "test", model.Label.ValueString())
	assert.Equal(t, "us-mia", model.Region.ValueString())
	assert.Equal(t, string(linodego.AffinityTypeAntiAffinityLocal), model.AffinityType.ValueString())
	assert.Equal(t, false, model.IsStrict.ValueBool())

	assert.Equal(t, false, model.IsCompliant.ValueBool())

	assert.Equal(t, int64(123), model.Members[0].LinodeID.ValueInt64())
	assert.Equal(t, false, model.Members[0].IsCompliant.ValueBool())

	assert.Equal(t, int64(456), model.Members[1].LinodeID.ValueInt64())
	assert.Equal(t, true, model.Members[1].IsCompliant.ValueBool())
}
