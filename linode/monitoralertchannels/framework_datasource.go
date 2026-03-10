package monitoralertchannels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_monitor_alert_channels",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_monitor_alert_channels")

	var data MonitorAlertChannelFilterModel

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
		ctx, d.Meta.Client, data.Filters, listMonitorAlertChannels,
		types.StringNull(), types.StringNull())
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.parseMonitorAlertChannels(helper.AnySliceToTyped[linodego.AlertChannel](result))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listMonitorAlertChannels(
	ctx context.Context,
	client *linodego.Client,
	_ string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListAlertChannels(...)")

	channels, err := client.ListAlertChannels(ctx, nil)
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(channels), nil
}
