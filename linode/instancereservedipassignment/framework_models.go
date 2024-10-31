package instancereservedipassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

type IPVPCNAT1To1Model struct {
	Address  types.String `tfsdk:"address"`
	SubnetID types.Int64  `tfsdk:"subnet_id"`
	VPCID    types.Int64  `tfsdk:"vpc_id"`
}

type InstanceIPModel struct {
	ID               types.String `tfsdk:"id"`
	LinodeID         types.Int64  `tfsdk:"linode_id"`
	Public           types.Bool   `tfsdk:"public"`
	Address          types.String `tfsdk:"address"`
	Gateway          types.String `tfsdk:"gateway"`
	Prefix           types.Int64  `tfsdk:"prefix"`
	RDNS             types.String `tfsdk:"rdns"`
	Region           types.String `tfsdk:"region"`
	SubnetMask       types.String `tfsdk:"subnet_mask"`
	Type             types.String `tfsdk:"type"`
	ApplyImmediately types.Bool   `tfsdk:"apply_immediately"`
	IPVPCNAT1To1     types.List   `tfsdk:"vpc_nat_1_1"`
	Reserved         types.Bool   `tfsdk:"reserved"`
}

func (m *InstanceIPModel) FlattenInstanceIP(
	ctx context.Context,
	ip linodego.InstanceIP,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateString(m.ID, ip.Address, preserveKnown)
	m.LinodeID = helper.KeepOrUpdateInt64(m.LinodeID, int64(ip.LinodeID), preserveKnown)
	m.Public = helper.KeepOrUpdateBool(m.Public, ip.Public, preserveKnown)
	m.Address = helper.KeepOrUpdateString(m.Address, ip.Address, preserveKnown)
	m.Gateway = helper.KeepOrUpdateString(m.Gateway, ip.Gateway, preserveKnown)
	m.Prefix = helper.KeepOrUpdateInt64(m.Prefix, int64(ip.Prefix), preserveKnown)

	m.RDNS = helper.KeepOrUpdateString(m.RDNS, ip.RDNS, preserveKnown)
	m.Reserved = helper.KeepOrUpdateBool(m.Reserved, ip.Reserved, preserveKnown)

	m.Region = helper.KeepOrUpdateString(m.Region, ip.Region, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateString(m.SubnetMask, ip.SubnetMask, preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(ip.Type), preserveKnown)

	var resultList types.List
	if ip.VPCNAT1To1 == nil {
		resultList = types.ListNull(instancenetworking.VPCNAT1To1Type)
	} else {
		vpcNAT1To1, diags := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
		diags.Append(diags...)
		if diags.HasError() {
			return diags
		}

		var listDiags diag.Diagnostics
		resultList, listDiags = types.ListValue(
			instancenetworking.VPCNAT1To1Type,
			[]attr.Value{vpcNAT1To1},
		)
		diags.Append(listDiags...)
		if diags.HasError() {
			return diags
		}
	}
	m.IPVPCNAT1To1 = helper.KeepOrUpdateValue(
		m.IPVPCNAT1To1,
		resultList,
		preserveKnown,
	)

	return diags
}

func (m *InstanceIPModel) CopyFrom(
	ctx context.Context,
	other InstanceIPModel,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateValue(m.ID, other.Address, preserveKnown)
	m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, other.LinodeID, preserveKnown)
	m.Public = helper.KeepOrUpdateValue(m.Public, other.Public, preserveKnown)
	m.Address = helper.KeepOrUpdateValue(m.Address, other.Address, preserveKnown)
	m.Gateway = helper.KeepOrUpdateValue(m.Gateway, other.Gateway, preserveKnown)
	m.Prefix = helper.KeepOrUpdateValue(m.Prefix, other.Prefix, preserveKnown)

	m.RDNS = helper.KeepOrUpdateValue(m.RDNS, other.RDNS, preserveKnown)
	m.Reserved = helper.KeepOrUpdateValue(m.Reserved, other.Reserved, preserveKnown)

	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateValue(m.SubnetMask, other.SubnetMask, preserveKnown)
	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
	m.IPVPCNAT1To1 = helper.KeepOrUpdateValue(
		m.IPVPCNAT1To1,
		other.IPVPCNAT1To1,
		preserveKnown,
	)

	return diags
}
