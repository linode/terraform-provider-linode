package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/linode/helper"
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
	var data FirewallModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	firewall, err := client.GetFirewall(ctx, firewallID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall",
			err.Error(),
		)
		return
	}
	rules, err := client.GetFirewallRules(ctx, firewallID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall rules",
			err.Error(),
		)
		return
	}
	devices, err := client.ListFirewallDevices(ctx, firewallID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall devices",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseComputedAttributes(ctx, firewall, rules, devices)...)
	resp.Diagnostics.Append(data.parseNonComputedAttributes(ctx, firewall, rules, devices)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
