//go:build unit

package lkenodepool

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
		ID:             123,
		Count:          3,
		Type:           "g6-standard-2",
		DiskEncryption: linodego.InstanceDiskEncryptionEnabled,
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

	nodePoolModel := NodePoolModel{}
	var diags diag.Diagnostics

	nodePoolModel.FlattenLKENodePool(&lkeNodePool, false, &diags)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", nodePoolModel.ID.ValueString())
	assert.Equal(t, int64(3), nodePoolModel.Count.ValueInt64())
	assert.Equal(t, "g6-standard-2", nodePoolModel.Type.ValueString())
	assert.Equal(t, "enabled", nodePoolModel.DiskEncryption.ValueString())
	assert.Len(t, nodePoolModel.Nodes.Elements(), 3)

	tags := make([]string, len(nodePoolModel.Tags.Elements()))
	for i, v := range nodePoolModel.Tags.Elements() {
		tags[i] = v.(types.String).ValueString()
	}
	assert.Contains(t, tags, "production")
	assert.Contains(t, tags, "web-server")

	// Example of asserting autoscaler values
	assert.NotNil(t, nodePoolModel.Autoscaler)
	assert.Equal(t, int64(1), nodePoolModel.Autoscaler[0].Min.ValueInt64())
	assert.Equal(t, int64(5), nodePoolModel.Autoscaler[0].Max.ValueInt64())
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
	tags, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"production", "web-server"})
	nodes, _ := flattenLKENodePoolLinodeList([]linodego.LKENodePoolLinode{
		{InstanceID: 1, ID: "linode123", Status: "running"},
		{InstanceID: 2, ID: "linode124", Status: "running"},
		{InstanceID: 3, ID: "linode125", Status: "running"},
	})

	nodePoolModel := NodePoolModel{
		ClusterID: types.Int64Value(1),
		Count:     types.Int64Value(3),
		Type:      types.StringValue("g6-standard-2"),
		Nodes:     *nodes,
		Tags:      tags,
		Autoscaler: []NodePoolAutoscalerModel{
			{
				Min: types.Int64Value(1),
				Max: types.Int64Value(5),
			},
		},
	}
	return &nodePoolModel
}
