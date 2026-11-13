package natgatewayaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_natgateway_addresses",
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
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data NATGatewayAddressesModel

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

	natGatewayID := helper.FrameworkSafeInt64ToInt(data.NATGatewayID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters,
		func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
			return listNATGatewayAddresses(ctx, client, filter, natGatewayID)
		},
		data.Order, data.OrderBy)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	resp.Diagnostics.Append(data.ParseNATGatewayAddresses(
		helper.AnySliceToTyped[linodego.NATGatewayAddressObject](result),
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listNATGatewayAddresses(
	ctx context.Context,
	client *linodego.Client,
	filter string,
	natGatewayID int,
) ([]any, error) {
	tflog.Trace(ctx, "client.NATGatewayListAddresses(...)", map[string]any{
		"natgateway_id": natGatewayID,
		"filter":        filter,
	})

	addresses, err := client.NATGatewayListAddresses(ctx, natGatewayID, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(addresses), nil
}
