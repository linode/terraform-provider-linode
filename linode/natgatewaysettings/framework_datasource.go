package natgatewaysettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_natgateway_settings",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_natgateway_settings")

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.NATGatewayGetSettings(...)")
	settings, err := d.Meta.Client.NATGatewayGetSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get NAT Gateway settings",
			"An error occurred while retrieving the NAT Gateway settings: "+err.Error(),
		)
		return
	}

	data.ParseNATGatewaySettings(settings)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
