package nfsquota

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nfs_quota",
				Schema: &DataSourceSchema,
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

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	filesystemID := helper.FrameworkSafeInt64ToInt(data.FilesystemID.ValueInt64(), &resp.Diagnostics)
	id := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"space_id":      spaceID,
		"filesystem_id": filesystemID,
		"id":            id,
	})

	tflog.Trace(ctx, "client.GetNFSQuota(...)")
	quota, err := client.GetNFSQuota(ctx, spaceID, filesystemID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get NFS Quota %v", id),
			err.Error(),
		)
		return
	}

	data.ParseQuota(quota)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
