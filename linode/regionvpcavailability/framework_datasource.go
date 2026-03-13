package regionvpcavailability

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
				Name:   "linode_region_vpc_availability",
				Schema: &FrameworkDataSourceSchema,
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
	tflog.Debug(ctx, "Read data."+r.Config.Name)

	var data RegionVPCAvailabilityModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	regionVPCAvailability, err := client.GetRegionVPCAvailability(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find the Linode Region VPC Availability",
			fmt.Sprintf(
				"Error finding the specified Linode Region VPC Availability %s: %s",
				data.ID.ValueString(),
				err.Error(),
			),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseRegionVPCAvailability(ctx, regionVPCAvailability)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
