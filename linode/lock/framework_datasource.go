package lock

import (
	"context"

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
				Name:          "linode_lock",
				Schema:        &frameworkDatasourceSchema,
				IsEarlyAccess: true,
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

	var data DataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lockID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "lock_id", lockID)

	lock, err := client.GetLock(ctx, lockID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get lock",
			err.Error(),
		)
		return
	}

	data.ParseLock(lock)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
