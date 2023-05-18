package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type DataSource struct {
	client *linodego.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_firewall"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data FirewallModel
	client := d.client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallId := int(data.ID.ValueInt64())
	firewall, err := client.GetFirewall(ctx, firewallId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall",
			err.Error(),
		)
		return
	}
	rules, err := client.GetFirewallRules(ctx, firewallId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall rules",
			err.Error(),
		)
		return
	}
	devices, err := client.ListFirewallDevices(ctx, firewallId, nil)
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
