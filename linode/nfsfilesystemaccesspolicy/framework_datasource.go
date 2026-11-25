package nfsfilesystemaccesspolicy

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
				Name:          "linode_nfs_filesystem_access_policy",
				Schema:        &frameworkDatasourceSchema,
				IsEarlyAccess: true,
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

	ctx = populateLogAttributes(ctx, data.ResourceModel)

	spaceID, filesystemID := getIDs(data.ResourceModel, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetNFSFilesystemAccessPolicy(...)")
	policy, err := client.GetNFSFilesystemAccessPolicy(ctx, spaceID, filesystemID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to read the access policy of NFS Filesystem %d (NFS Space %d).",
				filesystemID, spaceID,
			),
			err.Error(),
		)
		return
	}

	data.Flatten(ctx, spaceID, filesystemID, policy, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
