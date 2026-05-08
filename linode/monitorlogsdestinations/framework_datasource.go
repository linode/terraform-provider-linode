package monitorlogsdestinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_monitor_logs_destinations",
				Schema: &frameworkDataSourceSchema,
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

	var data LogsDestinationFilterModel

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
		ctx,
		d.Meta.Client,
		data.Filters,
		listLogsDestinations,
		data.Order,
		data.OrderBy,
	)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.ParseLogsDestinations(helper.AnySliceToTyped[linodego.LogsDestination](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listLogsDestinations(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	destinations, err := client.ListLogsDestinations(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(destinations), nil
}
