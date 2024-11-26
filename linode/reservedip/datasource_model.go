package reservedip

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Address      types.String `tfsdk:"address"`
	Region       types.String `tfsdk:"region"`
	Gateway      types.String `tfsdk:"gateway"`
	SubnetMask   types.String `tfsdk:"subnet_mask"`
	Prefix       types.Int64  `tfsdk:"prefix"`
	Type         types.String `tfsdk:"type"`
	Public       types.Bool   `tfsdk:"public"`
	RDNS         types.String `tfsdk:"rdns"`
	LinodeID     types.Int64  `tfsdk:"linode_id"`
	Reserved     types.Bool   `tfsdk:"reserved"`
	IPVPCNAT1To1 types.List   `tfsdk:"vpc_nat_1_1"`
}

func (data *DataSourceModel) parseIP(ip *linodego.InstanceIP) {
	data.ID = types.StringValue(ip.Address)
	data.Address = types.StringValue(ip.Address)
	data.Region = types.StringValue(ip.Region)
	data.Gateway = types.StringValue(ip.Gateway)
	data.SubnetMask = types.StringValue(ip.SubnetMask)
	data.Prefix = types.Int64Value(int64(ip.Prefix))
	data.Type = types.StringValue(string(ip.Type))
	data.Public = types.BoolValue(ip.Public)
	data.RDNS = types.StringValue(ip.RDNS)
	data.LinodeID = types.Int64Value(int64(ip.LinodeID))
	data.Reserved = types.BoolValue(ip.Reserved)
	var resultList types.List
	if ip.VPCNAT1To1 == nil {
		resultList = types.ListNull(instancenetworking.VPCNAT1To1Type)
	} else {
		vpcNAT1To1, _ := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)

		resultList, _ = types.ListValue(
			instancenetworking.VPCNAT1To1Type,
			[]attr.Value{vpcNAT1To1},
		)
	}
	data.IPVPCNAT1To1 = helper.KeepOrUpdateValue(
		data.IPVPCNAT1To1,
		resultList,
		true,
	)
}
