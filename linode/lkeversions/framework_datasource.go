package lkeversions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_versions",
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
	tflog.Debug(ctx, "Read data.linode_lke_versions")

	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Tier is populated
	if data.Tier.IsNull() {
		// If Tier is not populated, use ListLKEVersions
		tflog.Trace(ctx, "client.ListLKEVersions(...)")
		versions, err := client.ListLKEVersions(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get LKE Versions: %s", err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(data.parseLKEVersions(versions)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		// If Tier is populated, use ListLKETierVersions
		tflog.Trace(ctx, "client.ListLKETierVersions(...)")
		versions, err := client.ListLKETierVersions(ctx, data.Tier.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to get LKE Tier Versions", err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(data.parseLKETierVersions(versions)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
