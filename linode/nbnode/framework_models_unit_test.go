package nbnode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseNodeBalancerNode(t *testing.T) {
	mockNode := &linodego.NodeBalancerNode{
		ID:             54321,
		Address:        "192.168.210.120:80",
		Label:          "node54321",
		Status:         "UP",
		Weight:         50,
		Mode:           "accept",
		ConfigID:       4567,
		NodeBalancerID: 12345,
	}

	data := &DataSourceModel{}

	data.ParseNodeBalancerNode(mockNode)

	assert.Equal(t, types.Int64Value(54321), data.ID)
	assert.Equal(t, types.Int64Value(12345), data.NodeBalancerID)
	assert.Equal(t, types.Int64Value(4567), data.ConfigID)
	assert.Equal(t, types.StringValue("node54321"), data.Label)
	assert.Equal(t, types.Int64Value(50), data.Weight)
	assert.Equal(t, types.StringValue("accept"), data.Mode)
	assert.Equal(t, types.StringValue("192.168.210.120:80"), data.Address)
	assert.Equal(t, types.StringValue("UP"), data.Status)
}
