package networkingip

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

type DataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Address     types.String `tfsdk:"address"`
	Gateway     types.String `tfsdk:"gateway"`
	SubnetMask  types.String `tfsdk:"subnet_mask"`
	Prefix      types.Int64  `tfsdk:"prefix"`
	Type        types.String `tfsdk:"type"`
	Public      types.Bool   `tfsdk:"public"`
	RDNS        types.String `tfsdk:"rdns"`
	Reserved    types.Bool   `tfsdk:"reserved"`
	LinodeID    types.Int64  `tfsdk:"linode_id"`
	InterfaceID types.Int64  `tfsdk:"interface_id"`
	Region      types.String `tfsdk:"region"`
	VPCNAT1To1  types.Object `tfsdk:"vpc_nat_1_1"`
}

func (data *DataSourceModel) parseIP(ip *linodego.InstanceIP) diag.Diagnostics {
	data.ID = types.StringValue(ip.Address)
	data.Address = types.StringValue(ip.Address)
	data.Gateway = types.StringValue(ip.Gateway)
	data.SubnetMask = types.StringValue(ip.SubnetMask)
	data.Prefix = types.Int64Value(int64(ip.Prefix))
	data.Type = types.StringValue(string(ip.Type))
	data.Public = types.BoolValue(ip.Public)
	data.RDNS = types.StringValue(ip.RDNS)
	data.LinodeID = types.Int64Value(int64(ip.LinodeID))
	data.InterfaceID = types.Int64PointerValue(helper.IntPtrToInt64Ptr(ip.InterfaceID))
	data.Region = types.StringValue(ip.Region)
	data.Reserved = types.BoolValue(ip.Reserved)

	vpcNAT1To1, d := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
	if d.HasError() {
		return d
	}

	data.VPCNAT1To1 = vpcNAT1To1

	return nil
}
