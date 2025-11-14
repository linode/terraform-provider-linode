package producerimagesharegroupmember

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
				Name:   "linode_producer_image_share_group_member",
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

	shareGroupID := helper.FrameworkSafeInt64ToInt(
		data.ShareGroupID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenUUID := data.TokenUUID.ValueString()

	tflog.Trace(ctx, "client.ImageShareGroupGetMember(...)")
	m, err := client.ImageShareGroupGetMember(ctx, shareGroupID, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get Image Share Group Member %s", tokenUUID), err.Error(),
		)
		return
	}

	data.ParseImageShareGroupMember(m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
