package maintenancepolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_maintenance_policies",
				Schema: &frameworkDataSourceSchema,
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

	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.ListMaintenancePolicies(...)")
	policies, err := client.ListMaintenancePolicies(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Maintenance Policies", err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseMaintenancePolicies(policies)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
