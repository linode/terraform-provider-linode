package networkreservedips

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
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

type DataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Address     types.String `tfsdk:"address"`
	Region      types.String `tfsdk:"region"`
	Gateway     types.String `tfsdk:"gateway"`
	SubnetMask  types.String `tfsdk:"subnet_mask"`
	Prefix      types.Int64  `tfsdk:"prefix"`
	Type        types.String `tfsdk:"type"`
	Public      types.Bool   `tfsdk:"public"`
	RDNS        types.String `tfsdk:"rdns"`
	LinodeID    types.Int64  `tfsdk:"linode_id"`
	Reserved    types.Bool   `tfsdk:"reserved"`
	ReservedIPs types.List   `tfsdk:"reserved_ips"`
}

func (data *DataSourceModel) parseIP(ip *linodego.InstanceIP) {
	data.ID = types.StringValue(ip.Address)
	data.Address = types.StringValue(ip.Address)
	data.Region = types.StringValue(ip.Region)
	data.Gateway = types.StringValue(ip.Gateway)
	data.SubnetMask = types.StringValue(ip.SubnetMask)
	data.Prefix = types.Int64Value(int64(ip.Prefix))
	data.Type = types.StringValue(string(ip.Type))
	data.Public = types.BoolValue(ip.Public)
	data.RDNS = types.StringValue(ip.RDNS)
	data.LinodeID = types.Int64Value(int64(ip.LinodeID))
	data.Reserved = types.BoolValue(ip.Reserved)
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
	} else {
		// List all reserved IPs
		filter := ""
		if !data.Region.IsNull() {
			filter = fmt.Sprintf("{\"region\":\"%s\"}", data.Region.ValueString())
		}
		ips, err := d.Meta.Client.ListReservedIPAddresses(ctx, &linodego.ListOptions{Filter: filter})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to list Reserved IP Addresses",
				err.Error(),
			)
			return
		}

		reservedIPs := make([]ReservedIPObject, len(ips))
		for i, ip := range ips {
			reservedIPs[i] = ReservedIPObject{
				ID:         types.StringValue(ip.Address),
				Address:    types.StringValue(ip.Address),
				Region:     types.StringValue(ip.Region),
				Gateway:    types.StringValue(ip.Gateway),
				SubnetMask: types.StringValue(ip.SubnetMask),
				Prefix:     types.Int64Value(int64(ip.Prefix)),
				Type:       types.StringValue(string(ip.Type)),
				Public:     types.BoolValue(ip.Public),
				RDNS:       types.StringValue(ip.RDNS),
				LinodeID:   types.Int64Value(int64(ip.LinodeID)),
				Reserved:   types.BoolValue(ip.Reserved),
			}
		}

		reservedIPsValue, diags := types.ListValueFrom(ctx, reservedIPObjectType, reservedIPs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ReservedIPs = reservedIPsValue

		// If there are IPs, populate the first one's details for backwards compatibility
		if len(ips) > 0 {
			data.parseIP(&ips[0])
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
