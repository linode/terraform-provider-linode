package monitorlogsstreamsquotas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_monitor_logs_stream_quotas",
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

	var data LogStreamQuotaListModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	quotas, err := d.Meta.Client.ListLogStreamQuotas(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list log stream quotas", err.Error())
		return
	}

	data.parseQuotas(quotas)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
