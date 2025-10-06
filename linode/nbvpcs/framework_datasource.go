package nbvpcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_nodebalancer_vpcs",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_nodebalancer_vpcs")

	var data Model

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

	nodeBalancerId := helper.FrameworkSafeInt64ToInt(data.NodeBalancerID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, listWrapper(nodeBalancerId),
		data.Order, data.OrderBy)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.Parse(helper.AnySliceToTyped[linodego.NodeBalancerVPCConfig](result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listWrapper(
	nodeBalancerID int,
) frameworkfilter.ListFunc {
	return func(
		ctx context.Context,
		client *linodego.Client,
		filter string,
	) ([]any, error) {
		tflog.Trace(ctx, "client.ListNodeBalancerVPCConfigs(...)")

		nbs, err := client.ListNodeBalancerVPCConfigs(
			ctx,
			nodeBalancerID,
			&linodego.ListOptions{
				Filter: filter,
			},
		)
		if err != nil {
			return nil, err
		}

		return helper.TypedSliceToAny(nbs), nil
	}
}
