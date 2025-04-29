package objquota

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_object_storage_quota",
				Schema: &FrameworkDatasourceSchema,
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

	id := data.ID.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"quota_id": id,
	})

	quota, err := client.GetObjectStorageQuota(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get the Object Storage Quota: %s", err.Error(),
		)
		return
	}

	usage, err := client.GetObjectStorageQuotaUsage(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get the Object Storage Quota Usage: %s", err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseObjectStorageQuota(ctx, quota, usage)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
