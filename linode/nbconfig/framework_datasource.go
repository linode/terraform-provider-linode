package nbconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nodebalancer_config",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_nodebalancer_config")

	client := d.Client
	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeBalancerID := helper.FrameworkSafeInt64ToInt(data.NodeBalancerID.ValueInt64(), &resp.Diagnostics)
	configID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.GetNodeBalancerConfig(...)", map[string]any{
		"nodebalancer_id": nodeBalancerID,
		"config_id":       configID,
	})

	config, err := client.GetNodeBalancerConfig(ctx, nodeBalancerID, configID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to get nodebalancer config %d", data.ID.ValueInt64()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseNodebalancerConfig(config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.Set(ctx, &data)
}
