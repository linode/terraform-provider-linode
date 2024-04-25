package nb

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
				Name:   "linode_nodebalancer",
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
	tflog.Debug(ctx, "Read data.linode_nodebalancer")

	var data NodeBalancerDataSourceModel
	client := d.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeBalancerID := helper.FrameworkSafeInt64ToInt(
		data.ID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "nodebalancer_id", nodeBalancerID)
	tflog.Trace(ctx, "client.GetNodeBalancer(...)")

	nodeBalancer, err := client.GetNodeBalancer(ctx, nodeBalancerID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer %d", nodeBalancerID),
			err.Error(),
		)
	}

	tflog.Trace(ctx, "client.ListNodeBalancerFirewalls(...)")

	firewalls, err := client.ListNodeBalancerFirewalls(ctx, nodeBalancerID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to list firewalls assgiend to nodebalancer %d", nodeBalancerID),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(data.flattenNodeBalancer(ctx, nodeBalancer, firewalls)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
