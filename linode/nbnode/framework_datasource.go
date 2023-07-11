package nbnode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			"linode_nodebalancer_node",
			frameworkDatasourceSchema,
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.ID.ValueInt64())
	nodebalancerID := int(data.NodeBalancerID.ValueInt64())
	configID := int(data.ConfigID.ValueInt64())

	node, err := client.GetNodeBalancerNode(ctx, nodebalancerID, configID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer node with id %d:", id), err.Error(),
		)
		return
	}

	data.ParseNodeBalancerNode(node)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
