package networktransferprices

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_network_transfer_prices",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_network_transfer_prices")

	var data NetworkTransferPriceFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d := filterConfig.GenerateID(data.Filters)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}
	data.ID = id

	result, d := filterConfig.GetAndFilter(
		ctx, r.Meta.Client, data.Filters, listNetworkTransferPrices, data.Order, data.OrderBy)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.parseNetworkTransferPrices(helper.AnySliceToTyped[linodego.NetworkTransferPrice](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listNetworkTransferPrices(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Debug(ctx, "Listing Network Transfer Prices", map[string]any{
		"filter_header": filter,
	})

	types, err := client.ListNetworkTransferPrices(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(types), nil
}
