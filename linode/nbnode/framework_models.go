package nbnode

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

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

func (data *DataSourceModel) ParseNodeBalancerNode(nbnode *linodego.NodeBalancerNode) {
	data.ID = types.Int64Value(int64(nbnode.ID))
	data.NodeBalancerID = types.Int64Value(int64(nbnode.NodeBalancerID))
	data.ConfigID = types.Int64Value(int64(nbnode.ConfigID))
	data.Label = types.StringValue(nbnode.Label)
	data.Weight = types.Int64Value(int64(nbnode.Weight))
	data.Mode = types.StringValue(string(nbnode.Mode))
	data.Address = types.StringValue(nbnode.Address)
	data.Status = types.StringValue(nbnode.Status)
}
