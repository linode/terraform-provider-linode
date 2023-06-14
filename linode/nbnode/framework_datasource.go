package nbnode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DataSourceModel) parseNodeBalancerNode(nbnode *linodego.NodeBalancerNode) {
	data.ID = types.Int64Value(int64(nbnode.ID))
	data.NodeBalancerID = types.Int64Value(int64(nbnode.NodeBalancerID))
	data.ConfigID = types.Int64Value(int64(nbnode.ConfigID))
	data.Label = types.StringValue(nbnode.Label)
	data.Weight = types.Int64Value(int64(nbnode.Weight))
	data.Mode = types.StringValue(string(nbnode.Mode))
	data.Address = types.StringValue(nbnode.Address)
	data.Status = types.StringValue(nbnode.Status)
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

type DataSourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	ConfigID       types.Int64  `tfsdk:"config_id"`
	Label          types.String `tfsdk:"label"`
	Weight         types.Int64  `tfsdk:"weight"`
	Mode           types.String `tfsdk:"mode"`
	Address        types.String `tfsdk:"address"`
	Status         types.String `tfsdk:"status"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_nodebalancer_node"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.ID.ValueInt64())
	nodebalancerID := int(data.NodeBalancerID.ValueInt64())
	configID := int(data.ConfigID.ValueInt64())

	node, err := client.GetNodeBalancerNode(ctx, nodebalancerID, configID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get nodebalancer node with id %d:", id), err.Error(),
		)
		return
	}

	data.parseNodeBalancerNode(node)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
