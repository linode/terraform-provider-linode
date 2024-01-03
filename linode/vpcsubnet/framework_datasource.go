package vpcsubnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_vpc_subnet",
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
	client := d.Meta.Client

	var data VPCSubnetModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpcId := helper.FrameworkSafeInt64ToInt(data.VPCId.ValueInt64(), &resp.Diagnostics)
	id := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	vpcSubnet, err := client.GetVPCSubnet(ctx, vpcId, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read VPC Subnet %v", data.ID.ValueInt64()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenSubnet(ctx, vpcSubnet, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
