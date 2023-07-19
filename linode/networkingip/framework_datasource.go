package networkingip

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			"linode_networking_ip",
			frameworkDatasourceSchema,
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseIP(ip *linodego.InstanceIP) {
	data.Address = types.StringValue(ip.Address)
	data.Gateway = types.StringValue(ip.Gateway)
	data.SubnetMask = types.StringValue(ip.SubnetMask)
	data.Prefix = types.Int64Value(int64(ip.Prefix))
	data.Type = types.StringValue(string(ip.Type))
	data.Public = types.BoolValue(ip.Public)
	data.RDNS = types.StringValue(ip.RDNS)
	data.LinodeID = types.Int64Value(int64(ip.LinodeID))
	data.Region = types.StringValue(ip.Region)

	id, _ := json.Marshal(ip)

	data.ID = types.StringValue(string(id))
}

type DataSourceModel struct {
	Address    types.String `tfsdk:"address"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	Type       types.String `tfsdk:"type"`
	Public     types.Bool   `tfsdk:"public"`
	RDNS       types.String `tfsdk:"rdns"`
	LinodeID   types.Int64  `tfsdk:"linode_id"`
	Region     types.String `tfsdk:"region"`
	ID         types.String `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip, err := d.Meta.Client.GetIPAddress(ctx, data.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get IP Address: %s", err.Error(),
		)
		return
	}

	data.parseIP(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
