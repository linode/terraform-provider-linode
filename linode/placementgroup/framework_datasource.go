package placementgroup

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
				Name:   "linode_placement_group",
				Schema: &frameworkDatasourceSchema,
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
	tflog.Debug(ctx, "Read data.linode_placement_group")

	client := d.Meta.Client

	var data PlacementGroupModel

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

	tflog.SetField(ctx, "placement_group_id", id)

	tflog.Trace(ctx, "client.GetPlacementGroup(...)")
	pg, err := client.GetPlacementGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to get placement group %v", id), err.Error(),
		)
		return
	}

	data.parsePlacementGroup(pg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
