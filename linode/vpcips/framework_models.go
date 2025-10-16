package vpcips

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type ModelVPCIP struct {
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

	IPv6Range     types.String `tfsdk:"ipv6_range"`
	IPv6IsPublic  types.Bool   `tfsdk:"ipv6_is_public"`
	IPv6Addresses types.Set    `tfsdk:"ipv6_addresses"`
}

type ModelIPv6Address struct {
	SLAACAddress types.String `tfsdk:"slaac_address"`
}

func (m *ModelVPCIP) FlattenVPCIP(ctx context.Context, vpcIp *linodego.VPCIP, preserveKnown bool) diag.Diagnostics {
	var rd diag.Diagnostics

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
	m.InterfaceID = helper.KeepOrUpdateInt64(m.InterfaceID, int64(vpcIp.InterfaceID), preserveKnown)
	m.IPv6Range = helper.KeepOrUpdateStringPointer(m.IPv6Range, vpcIp.IPv6Range, preserveKnown)
	m.IPv6IsPublic = helper.KeepOrUpdateBoolPointer(m.IPv6IsPublic, vpcIp.IPv6IsPublic, preserveKnown)

	ipv6AddressModels := helper.MapSlice(vpcIp.IPv6Addresses,
		func(address linodego.VPCIPIPv6Address) ModelIPv6Address {
			return ModelIPv6Address{
				SLAACAddress: types.StringValue(address.SLAACAddress),
			}
		},
	)

	ipv6AddressSet, diags := types.SetValueFrom(ctx, DataSourceSchemaIPv6NestedObject.Type(), ipv6AddressModels)
	if diags.HasError() {
		return diags
	}

	m.IPv6Addresses = helper.KeepOrUpdateValue(
		m.IPv6Addresses,
		ipv6AddressSet,
		preserveKnown,
	)

	var newConfigID types.Int64
	if vpcIp.ConfigID == 0 {
		newConfigID = types.Int64Null()
	} else {
		newConfigID = types.Int64Value(int64(vpcIp.ConfigID))
	}
	m.ConfigID = helper.KeepOrUpdateValue(m.ConfigID, newConfigID, preserveKnown)

	return rd
}

type Model struct {
	ID      types.String                     `tfsdk:"id"`
	VPCID   types.Int64                      `tfsdk:"vpc_id"`
	IPv6    types.Bool                       `tfsdk:"ipv6"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCIPs  []ModelVPCIP                     `tfsdk:"vpc_ips"`
}

func (model *Model) FlattenVPCIPs(
	ctx context.Context,
	vpcIps []linodego.VPCIP,
	preserveKnown bool,
) diag.Diagnostics {
	var rd diag.Diagnostics
	vpcipModels := make([]ModelVPCIP, len(vpcIps))

	for i := range vpcIps {
		var vpcIp ModelVPCIP
		rd.Append(vpcIp.FlattenVPCIP(ctx, &vpcIps[i], preserveKnown)...)
		vpcipModels[i] = vpcIp
	}

	model.VPCIPs = vpcipModels
	return rd
}
