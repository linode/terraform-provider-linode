package consumerimagesharegrouptoken

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
				Name:   "linode_consumer_image_share_group_token",
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

	tokenUUID := data.TokenUUID.ValueString()

	tflog.Trace(ctx, "client.ImageShareGroupGetToken(...)")
	m, err := client.ImageShareGroupGetToken(ctx, tokenUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get Image Share Group Token %s", tokenUUID), err.Error(),
		)
		return
	}

	data.ParseImageShareGroupToken(m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
