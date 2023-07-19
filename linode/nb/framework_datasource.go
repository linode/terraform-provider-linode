package nb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			"linode_nodebalancer",
			frameworkDatasourceSchema,
		),
	}
}

type DataSourceModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	Label              types.String `tfsdk:"label"`
	Region             types.String `tfsdk:"region"`
	ClientConnThrottle types.Int64  `tfsdk:"client_conn_throttle"`
	Hostname           types.String `tfsdk:"hostname"`
	Ipv4               types.String `tfsdk:"ipv4"`
	Ipv6               types.String `tfsdk:"ipv6"`
	Created            types.String `tfsdk:"created"`
	Updated            types.String `tfsdk:"updated"`
	Transfer           types.List   `tfsdk:"transfer"`
	Tags               types.Set    `tfsdk:"tags"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data DataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodebalancer, err := client.GetNodeBalancer(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer %d", int(data.ID.ValueInt64())),
			err.Error(),
		)
	}

	resp.Diagnostics.Append(data.ParseNodeBalancer(ctx, nodebalancer)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
