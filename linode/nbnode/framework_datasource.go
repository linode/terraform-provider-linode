package nbnode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nodebalancer_node",
				Schema: &frameworkDatasourceSchema,
			},
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

	id := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	nodeBalancerID := helper.FrameworkSafeInt64ToInt(
		data.NodeBalancerID.ValueInt64(),
		&resp.Diagnostics,
	)
	configID := helper.FrameworkSafeInt64ToInt(data.ConfigID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	node, err := client.GetNodeBalancerNode(ctx, nodeBalancerID, configID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer node with id %d:", id), err.Error(),
		)
		return
	}

	data.ParseNodeBalancerNode(node)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
