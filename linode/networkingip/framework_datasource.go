package networkingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_networking_ip",
				Schema: &frameworkDatasourceSchema,
			},
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
	data.Reserved = types.BoolValue(ip.Reserved)
	data.ID = types.StringValue(ip.Address)
}

type DataSourceModel struct {
	Address    types.String `tfsdk:"address"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	Type       types.String `tfsdk:"type"`
	Public     types.Bool   `tfsdk:"public"`
	RDNS       types.String `tfsdk:"rdns"`
	Reserved   types.Bool   `tfsdk:"reserved"`
	LinodeID   types.Int64  `tfsdk:"linode_id"`
	Region     types.String `tfsdk:"region"`
	ID         types.String `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_networking_ip")

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "address", data.Address.ValueString())

	// Workaround: Use GetIPAddresses and filter for the specific address
	ips, err := d.Meta.Client.ListIPAddresses(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list IP Addresses",
			err.Error(),
		)
		return
	}

	var ip *linodego.InstanceIP
	for _, candidate := range ips {
		if candidate.Address == data.Address.ValueString() {
			ip = &candidate
			break
		}
	}

	// If the IP address is not found, return an error
	if ip == nil {
		resp.Diagnostics.AddError(
			"IP Address Not Found",
			"Could not find the specified IP address in the list of retrieved IPs.",
		)
		return
	}

	data.parseIP(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
