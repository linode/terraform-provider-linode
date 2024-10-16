package networkreservedip

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceFetchModel struct {
	ID         types.String `tfsdk:"id"`
	Address    types.String `tfsdk:"address"`
	Region     types.String `tfsdk:"region"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	Type       types.String `tfsdk:"type"`
	Public     types.Bool   `tfsdk:"public"`
	RDNS       types.String `tfsdk:"rdns"`
	LinodeID   types.Int64  `tfsdk:"linode_id"`
	Reserved   types.Bool   `tfsdk:"reserved"`
}

func (data *DataSourceFetchModel) parseIP(ip *linodego.InstanceIP) {
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
