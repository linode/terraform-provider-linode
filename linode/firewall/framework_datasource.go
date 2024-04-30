package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_firewall",
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
	tflog.Debug(ctx, "Read data.linode_firewall")

	var data FirewallDataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "firewall_id", firewallID)

	tflog.Trace(ctx, "client.GetFirewall(...)")
	firewall, err := client.GetFirewall(ctx, firewallID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.GetFirewallRules(...)")
	rules, err := client.GetFirewallRules(ctx, firewallID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall rules",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "client.ListFirewallDevices(...)")
	devices, err := client.ListFirewallDevices(ctx, firewallID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall devices",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.flattenFirewallForDataSource(ctx, firewall, devices, rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
