package vpcips

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type VPCIPModel struct {
	Address      types.String `tfsdk:"address"`
	AddressRange types.String `tfsdk:"address_range"`
	Gateway      types.String `tfsdk:"gateway"`
	SubnetMask   types.String `tfsdk:"subnet_mask"`
	Prefix       types.Int64  `tfsdk:"prefix"`
	LinodeID     types.Int64  `tfsdk:"linode_id"`
	Region       types.String `tfsdk:"region"`
	Active       types.Bool   `tfsdk:"active"`
	NAT1To1      types.String `tfsdk:"nat_1_1"`
	VPCID        types.Int64  `tfsdk:"vpc_id"`
	SubnetID     types.Int64  `tfsdk:"subnet_id"`
	ConfigID     types.Int64  `tfsdk:"config_id"`
	InterfaceID  types.Int64  `tfsdk:"interface_id"`
}

func (m *VPCIPModel) FlattenVPCIP(vpcIp *linodego.VPCIP, preserveKnown bool) {
	m.Address = helper.KeepOrUpdateStringPointer(m.Address, vpcIp.Address, preserveKnown)
	m.AddressRange = helper.KeepOrUpdateStringPointer(m.AddressRange, vpcIp.AddressRange, preserveKnown)
	m.Gateway = helper.KeepOrUpdateString(m.Gateway, vpcIp.Gateway, preserveKnown)
	m.SubnetMask = helper.KeepOrUpdateString(m.SubnetMask, vpcIp.SubnetMask, preserveKnown)
	m.Prefix = helper.KeepOrUpdateInt64(m.Prefix, int64(vpcIp.Prefix), preserveKnown)
	m.LinodeID = helper.KeepOrUpdateInt64(m.LinodeID, int64(vpcIp.LinodeID), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, vpcIp.Region, preserveKnown)
	m.Active = helper.KeepOrUpdateBool(m.Active, vpcIp.Active, preserveKnown)
	m.NAT1To1 = helper.KeepOrUpdateStringPointer(m.NAT1To1, vpcIp.NAT1To1, preserveKnown)
	m.VPCID = helper.KeepOrUpdateInt64(m.VPCID, int64(vpcIp.VPCID), preserveKnown)
	m.SubnetID = helper.KeepOrUpdateInt64(m.SubnetID, int64(vpcIp.SubnetID), preserveKnown)
	m.ConfigID = helper.KeepOrUpdateInt64(m.ConfigID, int64(vpcIp.ConfigID), preserveKnown)
	m.InterfaceID = helper.KeepOrUpdateInt64(m.InterfaceID, int64(vpcIp.InterfaceID), preserveKnown)
}

type VPCIPFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	VPCID   types.Int64                      `tfsdk:"vpc_id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCIPs  []VPCIPModel                     `tfsdk:"vpc_ips"`
}

func (model *VPCIPFilterModel) FlattenVPCIPs(
	vpcIps []linodego.VPCIP,
	preserveKnown bool,
) {
	vpcipModels := make([]VPCIPModel, len(vpcIps))

	for i := range vpcIps {
		var vpcIp VPCIPModel
		vpcIp.FlattenVPCIP(&vpcIps[i], preserveKnown)
		vpcipModels[i] = vpcIp
	}

	model.VPCIPs = vpcipModels
}
