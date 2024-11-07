package networkreservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSourceFetch() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_reserved_ip_fetch",
				Schema: &frameworkDataSourceFetchSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data.linode_reserved_ip")

	var data DataSourceFetchModel
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
	// } else {
	// 	// List all reserved IPs
	// 	filter := ""
	// 	ips, err := d.Meta.Client.ListReservedIPAddresses(ctx, &linodego.ListOptions{Filter: filter})
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable to list Reserved IP Addresses",
	// 			err.Error(),
	// 		)
	// 		return
	// 	}

	// 	reservedIPs := make([]ReservedIPObject, len(ips))
	// 	for i, ip := range ips {
	// 		reservedIPs[i] = ReservedIPObject{
	// 			ID:         types.StringValue(ip.Address),
	// 			Address:    types.StringValue(ip.Address),
	// 			Region:     types.StringValue(ip.Region),
	// 			Gateway:    types.StringValue(ip.Gateway),
	// 			SubnetMask: types.StringValue(ip.SubnetMask),
	// 			Prefix:     types.Int64Value(int64(ip.Prefix)),
	// 			Type:       types.StringValue(string(ip.Type)),
	// 			Public:     types.BoolValue(ip.Public),
	// 			RDNS:       types.StringValue(ip.RDNS),
	// 			LinodeID:   types.Int64Value(int64(ip.LinodeID)),
	// 			Reserved:   types.BoolValue(ip.Reserved),
	// 		}
	// 	}

	// 	reservedIPsValue, diags := types.ListValueFrom(ctx, reservedIPObjectType, reservedIPs)
	// 	resp.Diagnostics.Append(diags...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// 	data.ReservedIPs = reservedIPsValue

	// 	// If there are IPs, populate the first one's details for backwards compatibility
	// 	if len(ips) > 0 {
	// 		data.parseIP(&ips[0])
	// 	}
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
