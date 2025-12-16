package producerimagesharegroup

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
				Name:   "linode_producer_image_share_group",
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

	id := helper.FrameworkSafeInt64ToInt(
		data.ID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetImageShareGroup(...)")
	sg, err := client.GetImageShareGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get Image Share Group %v", id), err.Error(),
		)
		return
	}

	data.ParseImageShareGroup(sg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
