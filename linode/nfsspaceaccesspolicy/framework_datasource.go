package nfsspaceaccesspolicy

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:          "linode_nfs_space_access_policy",
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

	var data ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "space_id", strconv.FormatInt(data.SpaceID.ValueInt64(), 10))

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetNFSSpaceAccessPolicy(...)")
	policy, err := client.GetNFSSpaceAccessPolicy(ctx, spaceID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read the access policy of NFS Space %d.", spaceID),
			err.Error(),
		)
		return
	}

	data.Flatten(ctx, spaceID, policy, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
