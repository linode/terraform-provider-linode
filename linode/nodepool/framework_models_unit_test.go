//go:build unit

package nodepool

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseNodePool(t *testing.T) {
	lkeNodePool := linodego.LKENodePool{
		ID:    123,
		Count: 3,
		Type:  "g6-standard-2",
		Disks: []linodego.LKENodePoolDisk{
			{Size: 50, Type: "ssd"},
		},
		Linodes: []linodego.LKENodePoolLinode{
			{InstanceID: 1, ID: "linode123", Status: "running"},
			{InstanceID: 2, ID: "linode124", Status: "running"},
			{InstanceID: 3, ID: "linode125", Status: "running"},
		},
		Tags: []string{"production", "web-server"},
		Autoscaler: linodego.LKENodePoolAutoscaler{
			Enabled: true,
			Min:     1,
			Max:     5,
		},
	}

	clusterID := 1
	nodePoolModel := NodePoolModel{}
	var diags diag.Diagnostics

	nodePoolModel.ParseNodePool(context.Background(), clusterID, &lkeNodePool, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "1:123", nodePoolModel.ID.ValueString())
	assert.Equal(t, int64(1), nodePoolModel.ClusterID.ValueInt64())
	assert.Equal(t, int64(123), nodePoolModel.PoolID.ValueInt64())
	assert.Equal(t, int64(3), nodePoolModel.Count.ValueInt64())
	assert.Equal(t, "g6-standard-2", nodePoolModel.Type.ValueString())
	assert.Len(t, nodePoolModel.Nodes, 3) // Asserting that there are 3 nodes

	// Checking Tags - converting types.List to []string for assertion
	tags := make([]string, len(nodePoolModel.Tags.Elements()))
	for i, v := range nodePoolModel.Tags.Elements() {
		tags[i] = v.(types.String).ValueString()
	}
	assert.Contains(t, tags, "production")
	assert.Contains(t, tags, "web-server")

	// Example of asserting autoscaler values
	assert.NotNil(t, nodePoolModel.Autoscaler)
	assert.Equal(t, int64(1), nodePoolModel.Autoscaler.Min.ValueInt64())
	assert.Equal(t, int64(5), nodePoolModel.Autoscaler.Max.ValueInt64())
}

func TestSetNodePoolCreateOptions(t *testing.T) {
	nodePoolModel := createNodePoolModel()

	var createOpts linodego.LKENodePoolCreateOptions
	var diags diag.Diagnostics

	nodePoolModel.SetNodePoolCreateOptions(context.Background(), &createOpts, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, 3, createOpts.Count)
	assert.Equal(t, "g6-standard-2", createOpts.Type)
	assert.Contains(t, createOpts.Tags, "production")
	assert.Contains(t, createOpts.Tags, "web-server")

	assert.True(t, createOpts.Autoscaler.Enabled)
	assert.Equal(t, 1, createOpts.Autoscaler.Min)
	assert.Equal(t, 5, createOpts.Autoscaler.Max)
}

func TestSetNodePoolUpdateOptions(t *testing.T) {
	nodePoolModel := createNodePoolModel()

	var updateOpts linodego.LKENodePoolUpdateOptions
	var diags diag.Diagnostics

	nodePoolModel.SetNodePoolUpdateOptions(context.Background(), &updateOpts, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, 3, updateOpts.Count)
	assert.Contains(t, *updateOpts.Tags, "production")
	assert.Contains(t, *updateOpts.Tags, "web-server")

	assert.True(t, updateOpts.Autoscaler.Enabled)
	assert.Equal(t, 1, updateOpts.Autoscaler.Min)
	assert.Equal(t, 5, updateOpts.Autoscaler.Max)
}

func createNodePoolModel() *NodePoolModel {
	tags, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"production", "web-server"})
	nodePoolModel := NodePoolModel{
		ClusterID: types.Int64Value(1),
		Count:     types.Int64Value(3),
		Type:      types.StringValue("g6-standard-2"),
		Nodes: []NodePoolNodeModel{
			{InstanceID: types.Int64Value(1), ID: types.StringValue("linode123"), Status: types.StringValue("running")},
			{InstanceID: types.Int64Value(2), ID: types.StringValue("linode124"), Status: types.StringValue("running")},
			{InstanceID: types.Int64Value(3), ID: types.StringValue("linode125"), Status: types.StringValue("running")},
		},
		Tags: tags,
		Autoscaler: &NodePoolAutoscalerModel{
			Min: types.Int64Value(1),
			Max: types.Int64Value(5),
		},
	}
	return &nodePoolModel
}
