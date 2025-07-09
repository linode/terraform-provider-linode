package networkingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
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

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

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
