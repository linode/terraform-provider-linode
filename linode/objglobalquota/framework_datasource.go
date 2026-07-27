package objglobalquota

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_object_storage_global_quota",
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

	quotaID := data.QuotaID.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"quota_id": quotaID,
	})

	quota, err := client.GetObjectStorageGlobalQuota(ctx, quotaID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get the Object Storage global quota",
			err.Error(),
		)
		return
	}

	var usage *linodego.ObjectStorageGlobalQuotaUsage
	if quota.HasUsage {
		usage, err = client.GetObjectStorageGlobalQuotaUsage(ctx, quotaID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to get the Object Storage global quota usage",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(data.parseObjectStorageGlobalQuota(ctx, quota, usage)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
