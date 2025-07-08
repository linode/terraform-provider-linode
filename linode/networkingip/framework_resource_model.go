package networkingip

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

type ResourceModel struct {
	ID types.String `tfsdk:"id"`

	LinodeID types.Int64  `tfsdk:"linode_id"`
	Reserved types.Bool   `tfsdk:"reserved"`
	Region   types.String `tfsdk:"region"`
	Public   types.Bool   `tfsdk:"public"`
	Type     types.String `tfsdk:"type"`
	RDNS     types.String `tfsdk:"rdns"`

	Address    types.String `tfsdk:"address"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	VPCNAT1To1 types.Object `tfsdk:"vpc_nat_1_1"`
}

func (m *ResourceModel) FlattenIPAddress(
	ip *linodego.InstanceIP,
	preserveKnown bool,
) diag.Diagnostics {
	m.ID = helper.KeepOrUpdateString(m.ID, ip.Address, preserveKnown)

	if ip.LinodeID != 0 {
		m.LinodeID = helper.KeepOrUpdateValue(
			m.LinodeID,
			types.Int64Value(int64(ip.LinodeID)),
			preserveKnown,
		)
	} else {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Null(), preserveKnown)
	}

	m.RDNS = helper.KeepOrUpdateValue(m.RDNS, types.StringValue(ip.RDNS), preserveKnown)
	m.Reserved = helper.KeepOrUpdateValue(m.Reserved, types.BoolValue(ip.Reserved), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, ip.Region, preserveKnown)
	m.Public = helper.KeepOrUpdateValue(m.Public, types.BoolValue(ip.Public), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(ip.Type), preserveKnown)

	m.Address = helper.KeepOrUpdateString(m.Address, ip.Address, preserveKnown)
	m.Gateway = helper.KeepOrUpdateString(m.Gateway, ip.Gateway, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateString(m.SubnetMask, ip.SubnetMask, preserveKnown)
	m.Prefix = helper.KeepOrUpdateInt64(m.Prefix, int64(ip.Prefix), preserveKnown)

	vpcNAT1To1, d := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
	if d.HasError() {
		return d
	}

	m.VPCNAT1To1 = helper.KeepOrUpdateValue(m.VPCNAT1To1, vpcNAT1To1, preserveKnown)

	return nil
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, other.LinodeID, preserveKnown)
	m.Reserved = helper.KeepOrUpdateValue(m.Reserved, other.Reserved, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.Public = helper.KeepOrUpdateValue(m.Public, other.Public, preserveKnown)
	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
	m.RDNS = helper.KeepOrUpdateValue(m.RDNS, other.RDNS, preserveKnown)
	m.Address = helper.KeepOrUpdateValue(m.Address, other.Address, preserveKnown)
	m.Gateway = helper.KeepOrUpdateValue(m.Gateway, other.Gateway, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateValue(m.SubnetMask, other.SubnetMask, preserveKnown)
	m.Prefix = helper.KeepOrUpdateValue(m.Prefix, other.Prefix, preserveKnown)
	m.VPCNAT1To1 = helper.KeepOrUpdateValue(m.VPCNAT1To1, other.VPCNAT1To1, preserveKnown)
}
