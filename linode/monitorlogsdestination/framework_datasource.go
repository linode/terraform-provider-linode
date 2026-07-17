package monitorlogsdestination

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_monitor_logs_destination",
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
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	client := d.Meta.Client

	var data LogsDestinationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", id)

	tflog.Trace(ctx, "client.GetLogsDestination(...)")
	dest, err := client.GetLogsDestination(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get logs destination %v", id),
			err.Error(),
		)
		return
	}

	data.ParseLogsDestination(dest)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
