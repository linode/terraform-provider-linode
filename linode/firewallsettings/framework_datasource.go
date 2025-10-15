package firewallsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_firewall_settings",
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
	var state FirewallSettingsModel

	client := d.Meta.Client

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)

	firewallSettings, err := client.GetFirewallSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall settings",
			"An error occurred while retrieving the firewall settings: "+err.Error(),
		)
		return
	}

	state.FlattenFirewallSettings(ctx, *firewallSettings, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
