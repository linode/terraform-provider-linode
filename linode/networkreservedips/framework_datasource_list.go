package networkreservedips

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSourceList() datasource.DataSource {
	return &DataSourceList{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_reserved_ip_list",
				Schema: &frameworkDataSourceListSchema,
			},
		),
	}
}

type DataSourceList struct {
	helper.BaseDataSource
}

func (d *DataSourceList) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data.linode_reserved_ip_list")

	var data DataSourceListModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := d.Meta.Client

	// List all reserved IPs
	ips, err := client.ListReservedIPAddresses(ctx, nil)
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
			Reserved:   types.BoolValue(true), // All IPs in this list are reserved
		}
	}

	reservedIPsValue, diags := types.ListValueFrom(ctx, reservedIPObjectType, reservedIPs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ReservedIPs = reservedIPsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
