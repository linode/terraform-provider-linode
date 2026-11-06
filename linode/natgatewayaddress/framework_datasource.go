package natgatewayaddress

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
				Name:   "linode_natgateway_address",
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

	var data NATGatewayAddressDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayID := helper.FrameworkSafeInt64ToInt(data.NATGatewayID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	address := data.Address.ValueString()

	tflog.Trace(ctx, "client.NATGatewayGetAddress(...)", map[string]any{
		"natgateway_id": natGatewayID,
		"address":       address,
	})
	natGatewayAddress, err := d.Meta.Client.NATGatewayGetAddress(ctx, natGatewayID, address)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read NAT Gateway address",
			err.Error(),
		)
		return
	}

	data.ParseNATGatewayAddress(natGatewayAddress)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
