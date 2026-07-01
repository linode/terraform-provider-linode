package monitorlogsstreams

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
				Name:   "linode_monitor_logs_streams",
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

	var data LogsStreamFilterModel

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
		listLogStreams,
		data.Order,
		data.OrderBy,
	)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.ParseLogStreams(helper.AnySliceToTyped[linodego.Stream](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listLogStreams(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	streams, err := client.ListLogStreams(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(streams), nil
}
