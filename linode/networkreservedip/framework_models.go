package networkreservedip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ReservedIPModel struct {
	ID         types.String `tfsdk:"id"`
	Region     types.String `tfsdk:"region"`
	Address    types.String `tfsdk:"address"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	Type       types.String `tfsdk:"type"`
	Public     types.Bool   `tfsdk:"public"`
	RDNS       types.String `tfsdk:"rdns"`
	LinodeID   types.Int64  `tfsdk:"linode_id"`
	Reserved   types.Bool   `tfsdk:"reserved"`
}

func (m *ReservedIPModel) FlattenReservedIP(
	ctx context.Context,
	ip linodego.InstanceIP,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateString(m.ID, ip.Address, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, ip.Region, preserveKnown)
	m.Address = helper.KeepOrUpdateString(m.Address, ip.Address, preserveKnown)
	m.Gateway = helper.KeepOrUpdateString(m.Gateway, ip.Gateway, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateString(m.SubnetMask, ip.SubnetMask, preserveKnown)
	m.Prefix = helper.KeepOrUpdateInt64(m.Prefix, int64(ip.Prefix), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(ip.Type), preserveKnown)
	m.Public = helper.KeepOrUpdateBool(m.Public, ip.Public, preserveKnown)
	m.RDNS = helper.KeepOrUpdateString(m.RDNS, ip.RDNS, preserveKnown)
	m.LinodeID = helper.KeepOrUpdateInt64(m.LinodeID, int64(ip.LinodeID), preserveKnown)
	m.Reserved = helper.KeepOrUpdateBool(m.Reserved, ip.Reserved, preserveKnown)

	return diags
}

// func (m *ReservedIPModel) CopyFrom(
// 	ctx context.Context,
// 	other ReservedIPModel,
// 	preserveKnown bool,
// ) diag.Diagnostics {
// 	var diags diag.Diagnostics

// 	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
// 	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
// 	m.Address = helper.KeepOrUpdateValue(m.Address, other.Address, preserveKnown)
// 	m.Gateway = helper.KeepOrUpdateValue(m.Gateway, other.Gateway, preserveKnown)
// 	m.SubnetMask = helper.KeepOrUpdateValue(m.SubnetMask, other.SubnetMask, preserveKnown)
// 	m.Prefix = helper.KeepOrUpdateValue(m.Prefix, other.Prefix, preserveKnown)
// 	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
// 	m.Public = helper.KeepOrUpdateValue(m.Public, other.Public, preserveKnown)
// 	m.RDNS = helper.KeepOrUpdateValue(m.RDNS, other.RDNS, preserveKnown)
// 	m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, other.LinodeID, preserveKnown)
// 	m.Reserved = helper.KeepOrUpdateValue(m.Reserved, other.Reserved, preserveKnown)

// 	return diags
// }
