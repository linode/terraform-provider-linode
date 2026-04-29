package regionsvpcavailability

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
				Name:   "linode_regions_vpc_availability",
				Schema: &frameworkDataSourceSchema,
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

	var data regionsVPCAvailabilityModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	regionsVPCAvailability, err := client.ListRegionsVPCAvailability(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find the Linode Regions VPC Availability",
			fmt.Sprintf(
				"Error finding Linode Regions VPC Availability: %s",
				err.Error(),
			),
		)
		return
	}

	resp.Diagnostics.Append(data.parseRegionsVPCAvailability(ctx, regionsVPCAvailability)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
