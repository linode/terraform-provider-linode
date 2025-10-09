package nbvpc

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	ID             types.Int64  `tfsdk:"id"`
	IPv4Range      types.String `tfsdk:"ipv4_range"`
	SubnetID       types.Int64  `tfsdk:"subnet_id"`
	VPCID          types.Int64  `tfsdk:"vpc_id"`
}

func (m *DataSourceModel) Flatten(vpcConfig *linodego.NodeBalancerVPCConfig) *DataSourceModel {
	m.NodeBalancerID = types.Int64Value(int64(vpcConfig.NodeBalancerID))
	m.ID = types.Int64Value(int64(vpcConfig.ID))

	m.IPv4Range = types.StringValue(vpcConfig.IPv4Range)

	m.VPCID = types.Int64Value(int64(vpcConfig.VPCID))
	m.SubnetID = types.Int64Value(int64(vpcConfig.SubnetID))

	return m
}
