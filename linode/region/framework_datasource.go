package region

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_region",
				Schema: &DataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_region")

	client := r.Meta.Client

	var data RegionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	ctx = tflog.SetField(ctx, "region_id", id)

	region, err := client.GetRegion(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find the Linode Region",
			fmt.Sprintf(
				"Error finding the specified Linode Region %s: %s",
				id,
				err.Error(),
			),
		)
		return
	}

	data.ParseRegion(region)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
