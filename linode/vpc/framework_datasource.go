package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_vpc",
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
	tflog.Debug(ctx, "Read data.linode_vpc")

	client := d.Meta.Client

	var data VPCModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetVPC(...)")
	vpc, err := client.GetVPC(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read VPC %v", id),
			err.Error(),
		)
		return
	}

	data.FlattenVPC(ctx, vpc, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
