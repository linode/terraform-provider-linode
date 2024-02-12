package nbconfigs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_nodebalancer_configs",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_nodebalancer_configs")

	var data NodeBalancerConfigFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diag := filterConfig.GenerateID(data.Filters)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	data.ID = id

	result, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, data.listNodeBalancerConfigs,
		// There are no API filterable fields, so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.parseNodeBalancerConfigs(helper.AnySliceToTyped[linodego.NodeBalancerConfig](result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *NodeBalancerConfigFilterModel) listNodeBalancerConfigs(
	ctx context.Context,
	client *linodego.Client,
	filter string,
) ([]any, error) {
	nbId, err := helper.SafeInt64ToInt(data.NodeBalancerId.ValueInt64())
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "nodebalancer_id", nbId)
	tflog.Trace(ctx, "client.ListNodeBalancerConfigs(...)")

	nbs, err := client.ListNodeBalancerConfigs(ctx, nbId, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(nbs), nil
}
