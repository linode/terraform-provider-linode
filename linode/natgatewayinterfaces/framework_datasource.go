package natgatewayinterfaces

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_natgateway_interfaces",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_natgateway_interfaces")

	var data Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayID := int(data.NATGatewayID.ValueInt64())

	var interfaces []linodego.NATGatewayInterface
	var err error

	if data.Address.IsNull() || data.Address.IsUnknown() {
		data.ID = types.StringValue(strconv.Itoa(natGatewayID))

		tflog.Trace(ctx, "client.NATGatewayListInterfaces(...)", map[string]any{
			"natgateway_id": natGatewayID,
		})

		interfaces, err = d.Meta.Client.NATGatewayListInterfaces(ctx, natGatewayID, nil)
	} else {
		address := data.Address.ValueString()
		data.ID = types.StringValue(fmt.Sprintf("%d/%s", natGatewayID, address))

		tflog.Trace(ctx, "client.NATGatewayListAddressInterfaces(...)", map[string]any{
			"natgateway_id": natGatewayID,
			"address":       address,
		})

		interfaces, err = d.Meta.Client.NATGatewayListAddressInterfaces(ctx, natGatewayID, address, nil)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to list NAT Gateway interfaces",
			fmt.Sprintf("Failed to list interfaces for NAT Gateway %d: %s", natGatewayID, err),
		)
		return
	}

	data.ParseInterfaces(interfaces)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
