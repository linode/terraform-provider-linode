package databasevalkeyconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
		Name: "linode_database_valkey_config", Schema: &frameworkDataSourceSchema,
	})}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config, err := d.Meta.Client.GetValkeyDatabaseConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Valkey Config", "Unable to retrieve Valkey Config: "+err.Error())
		return
	}
	data.ParseValkeyConfig(ctx, config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
