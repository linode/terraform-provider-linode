package natgateway

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var (
	natGatewayAddressAttrTypes = map[string]attr.Type{
		"address": types.StringType,
	}
	natGatewayAddressObjectType  = types.ObjectType{AttrTypes: natGatewayAddressAttrTypes}
	natGatewayVPCSubnetAttrTypes = map[string]attr.Type{
		"id":        types.StringType,
		"type":      types.StringType,
		"url":       types.StringType,
		"label":     types.StringType,
		"vpc_id":    types.StringType,
		"vpc_label": types.StringType,
	}
)

type NATGatewayModel struct {
	ID                       types.Int64       `tfsdk:"id"`
	Region                   types.String      `tfsdk:"region"`
	Label                    types.String      `tfsdk:"label"`
	DefaultPortsPerInterface types.Int64       `tfsdk:"default_ports_per_interface"`
	UseAutoscaling           types.Bool        `tfsdk:"use_autoscaling"`
	AddressAutoscaleMax      types.Int64       `tfsdk:"address_autoscale_max"`
	PortsetAssignments       types.Int64       `tfsdk:"portset_assignments"`
	PortsetCapacity          types.Int64       `tfsdk:"portset_capacity"`
	VPCSubnet                types.Object      `tfsdk:"vpc_subnet"`
	Created                  timetypes.RFC3339 `tfsdk:"created"`
	Updated                  timetypes.RFC3339 `tfsdk:"updated"`
	Addresses                types.Set         `tfsdk:"addresses"`
}

func natGatewayVPCSubnetObject(gw *linodego.NATGateway) (types.Object, diag.Diagnostics) {
	if gw.VPCSubnet.ID == 0 {
		return types.ObjectNull(natGatewayVPCSubnetAttrTypes), nil
	}

	return types.ObjectValue(natGatewayVPCSubnetAttrTypes, map[string]attr.Value{
		"id":        types.StringValue(strconv.Itoa(gw.VPCSubnet.ID)),
		"type":      types.StringValue(gw.VPCSubnet.Type),
		"url":       types.StringValue(gw.VPCSubnet.URL),
		"label":     types.StringValue(gw.VPCSubnet.Label),
		"vpc_id":    types.StringValue(strconv.Itoa(gw.VPCSubnet.VPCID)),
		"vpc_label": types.StringValue(gw.VPCSubnet.VPCLabel),
	})
}

func natGatewayAddressesSet(gw *linodego.NATGateway) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	addresses := make([]attr.Value, 0, len(gw.Addresses))
	for _, a := range gw.Addresses {
		obj, d := types.ObjectValue(natGatewayAddressAttrTypes, map[string]attr.Value{
			"address": types.StringValue(a.Address),
		})
		diags.Append(d...)
		addresses = append(addresses, obj)
	}

	set, d := types.SetValue(natGatewayAddressObjectType, addresses)
	diags.Append(d...)

	return set, diags
}

func (m *NATGatewayModel) FlattenNATGateway(gw *linodego.NATGateway, preserveKnown bool) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateInt64(m.ID, int64(gw.ID), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, gw.Region, preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, gw.Label, preserveKnown)
	m.DefaultPortsPerInterface = helper.KeepOrUpdateInt64(
		m.DefaultPortsPerInterface, int64(gw.DefaultPortsPerInterface), preserveKnown,
	)
	m.AddressAutoscaleMax = helper.KeepOrUpdateInt64(m.AddressAutoscaleMax, int64(gw.AddressAutoscaleMax), preserveKnown)

	// The API does not return use_autoscaling; address_autoscale_max is non-zero only when it is enabled.
	m.UseAutoscaling = helper.KeepOrUpdateValue(m.UseAutoscaling, types.BoolValue(gw.AddressAutoscaleMax > 0), preserveKnown)

	m.PortsetAssignments = helper.KeepOrUpdateInt64(m.PortsetAssignments, int64(gw.PortsetAssignments), preserveKnown)
	m.PortsetCapacity = helper.KeepOrUpdateInt64(m.PortsetCapacity, int64(gw.PortsetCapacity), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(gw.Created), preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(gw.Updated), preserveKnown)

	vpcSubnet, d := natGatewayVPCSubnetObject(gw)
	diags.Append(d...)
	m.VPCSubnet = helper.KeepOrUpdateValue(m.VPCSubnet, vpcSubnet, preserveKnown)

	addresses, d := natGatewayAddressesSet(gw)
	diags.Append(d...)
	m.Addresses = helper.KeepOrUpdateValue(m.Addresses, addresses, preserveKnown)

	return diags
}

type natGatewayAddressModel struct {
	Address types.String `tfsdk:"address"`
}

func (m *NATGatewayModel) ExpandCreateOptions(ctx context.Context) (linodego.NATGatewayCreateOptions, diag.Diagnostics) {
	var diags diag.Diagnostics

	opts := linodego.NATGatewayCreateOptions{
		Region:    m.Region.ValueString(),
		Label:     m.Label.ValueString(),
		Addresses: []linodego.NATGatewayAddressCreateOptions{},
	}

	if !m.Addresses.IsNull() && !m.Addresses.IsUnknown() {
		var addresses []natGatewayAddressModel
		diags.Append(m.Addresses.ElementsAs(ctx, &addresses, false)...)
		for _, a := range addresses {
			opts.Addresses = append(opts.Addresses, linodego.NATGatewayAddressCreateOptions{
				Address: a.Address.ValueString(),
			})
		}
	}

	if !m.DefaultPortsPerInterface.IsNull() && !m.DefaultPortsPerInterface.IsUnknown() {
		v := int(m.DefaultPortsPerInterface.ValueInt64())
		opts.DefaultPortsPerInterface = &v
	}

	if !m.UseAutoscaling.IsNull() && !m.UseAutoscaling.IsUnknown() {
		v := m.UseAutoscaling.ValueBool()
		opts.UseAutoscaling = &v
	}

	return opts, diags
}

type NATGatewayDataSourceModel struct {
	ID                       types.Int64       `tfsdk:"id"`
	Region                   types.String      `tfsdk:"region"`
	Label                    types.String      `tfsdk:"label"`
	DefaultPortsPerInterface types.Int64       `tfsdk:"default_ports_per_interface"`
	AddressAutoscaleMax      types.Int64       `tfsdk:"address_autoscale_max"`
	PortsetAssignments       types.Int64       `tfsdk:"portset_assignments"`
	PortsetCapacity          types.Int64       `tfsdk:"portset_capacity"`
	VPCSubnet                types.Object      `tfsdk:"vpc_subnet"`
	Created                  timetypes.RFC3339 `tfsdk:"created"`
	Updated                  timetypes.RFC3339 `tfsdk:"updated"`
	Addresses                types.Set         `tfsdk:"addresses"`
}

func (m *NATGatewayDataSourceModel) ParseNATGateway(gw *linodego.NATGateway) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.Int64Value(int64(gw.ID))
	m.Region = types.StringValue(gw.Region)
	m.Label = types.StringValue(gw.Label)
	m.DefaultPortsPerInterface = types.Int64Value(int64(gw.DefaultPortsPerInterface))
	m.AddressAutoscaleMax = types.Int64Value(int64(gw.AddressAutoscaleMax))
	m.PortsetAssignments = types.Int64Value(int64(gw.PortsetAssignments))
	m.PortsetCapacity = types.Int64Value(int64(gw.PortsetCapacity))
	m.Created = timetypes.NewRFC3339TimePointerValue(gw.Created)
	m.Updated = timetypes.NewRFC3339TimePointerValue(gw.Updated)

	vpcSubnet, d := natGatewayVPCSubnetObject(gw)
	diags.Append(d...)
	m.VPCSubnet = vpcSubnet

	addresses, d := natGatewayAddressesSet(gw)
	diags.Append(d...)
	m.Addresses = addresses

	return diags
}
