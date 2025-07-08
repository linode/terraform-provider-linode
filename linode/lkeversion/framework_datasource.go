package lkeversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_version",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	client := d.Meta.Client

	var data DataSourceModel

	// Get the config data
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Tier is populated
	if data.Tier.IsNull() {
		// If Tier is not populated, use GetLKEVersion
		versionInfo, err := client.GetLKEVersion(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting lke version %s: ", data.ID.ValueString()), err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(data.ParseLKEVersion(versionInfo)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		// If Tier is populated, use GetLKETierVersion
		versionInfo, err := client.GetLKETierVersion(ctx, data.Tier.ValueString(), data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error getting lke version %s %s: ", data.Tier.ValueString(), data.ID.ValueString()), err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(data.ParseLKETierVersion(versionInfo)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
