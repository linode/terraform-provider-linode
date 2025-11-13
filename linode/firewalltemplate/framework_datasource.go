package firewalltemplate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_firewall_template",
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
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data FirewallTemplateDataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	slug := data.Slug.ValueString()
	ctx = tflog.SetField(ctx, "slug", slug)

	tflog.Trace(ctx, "client.GetFirewallTemplate(...)", map[string]any{
		"slug": slug,
	})
	firewallTemplate, err := client.GetFirewallTemplate(ctx, slug)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to Get the Firewall Template %q", slug),
			err.Error(),
		)
		return
	}

	data.parseFirewallTemplate(ctx, *firewallTemplate, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
