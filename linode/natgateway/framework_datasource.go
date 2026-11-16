package natgateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_natgateway",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data NATGatewayDataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natgatewayID := helper.FrameworkSafeInt64ToInt(
		data.ID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "natgateway_id", natgatewayID)

	natgateway, err := client.GetNATGateway(ctx, natgatewayID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get natgateway %d", natgatewayID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseNATGateway(natgateway)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
