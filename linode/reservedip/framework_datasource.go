package reservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_reserved_ip",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data.linode_reserved_ip")

	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Address.IsNull() {
		// Fetch a specific reserved IP
		ip, err := d.Meta.Client.GetReservedIPAddress(ctx, data.Address.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get Reserved IP Address",
				err.Error(),
			)
			return
		}
		data.parseIP(ip)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
