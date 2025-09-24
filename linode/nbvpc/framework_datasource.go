package nbvpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nodebalancer_vpc",
				Schema: &DataSourceSchema,
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
	tflog.Debug(ctx, "Read data.linode_nodebalancer_vpc")

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
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"id":              id,
		"nodebalancer_id": nodeBalancerID,
	})

	vpcConfig, err := client.GetNodeBalancerVPCConfig(ctx, nodeBalancerID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get NodeBalancer-VPC config with id %d:", id), err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data.Flatten(vpcConfig))...)
}
