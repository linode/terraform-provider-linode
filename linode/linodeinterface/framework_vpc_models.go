package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

type VPCAttrModel struct {
	IPv4     types.Object `tfsdk:"ipv4"`
	IPv6     types.Object `tfsdk:"ipv6"`
	SubnetID types.Int64  `tfsdk:"subnet_id"`
}

type VPCIPv4AttrModel struct {
	Addresses         types.List   `tfsdk:"addresses"`
	AssignedAddresses types.Set    `tfsdk:"assigned_addresses"`
	Ranges            types.List   `tfsdk:"ranges"`
	AssignedRanges    types.Set    `tfsdk:"assigned_ranges"`
	NATGateway        types.Object `tfsdk:"natgateway"`
}

// vpcInterfaceNATGatewayPortAttrTypes is the attribute type map of a NAT gateway portset port.
var vpcInterfaceNATGatewayPortAttrTypes = map[string]attr.Type{
	"start": types.Int64Type,
	"end":   types.Int64Type,
}

// vpcInterfaceNATGatewayPortsetAttrTypes is the attribute type map of a NAT gateway portset.
var vpcInterfaceNATGatewayPortsetAttrTypes = map[string]attr.Type{
	"address": types.StringType,
	"ports":   types.ListType{ElemType: types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortAttrTypes}},
}

// vpcInterfaceNATGatewayAttrTypes is the attribute type map of the natgateway block on a VPC interface's IPv4 config.
// It must match the schema of `natgateway` in both the resource and the data source.
var vpcInterfaceNATGatewayAttrTypes = map[string]attr.Type{
	"id":                  types.Int64Type,
	"label":               types.StringType,
	"addresses":           types.ListType{ElemType: types.StringType},
	"portset_assignments": types.Int64Type,
	"portset_capacity":    types.Int64Type,
	"portsets":            types.ListType{ElemType: types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortsetAttrTypes}},
}

type VPCInterfaceNATGatewayPortModel struct {
	Start types.Int64 `tfsdk:"start"`
	End   types.Int64 `tfsdk:"end"`
}

type VPCInterfaceNATGatewayPortsetModel struct {
	Address types.String `tfsdk:"address"`
	Ports   types.List   `tfsdk:"ports"`
}

type VPCInterfaceNATGatewayModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	Label              types.String `tfsdk:"label"`
	Addresses          types.List   `tfsdk:"addresses"`
	PortsetAssignments types.Int64  `tfsdk:"portset_assignments"`
	PortsetCapacity    types.Int64  `tfsdk:"portset_capacity"`
	Portsets           types.List   `tfsdk:"portsets"`
}

// flattenVPCInterfaceNATGateway converts the NAT gateway of a VPC interface's IPv4 config into a
// natgateway object. A nil gateway yields a null object.
func flattenVPCInterfaceNATGateway(
	ctx context.Context,
	gateway *linodego.VPCInterfaceIPv4NATGateway,
	diags *diag.Diagnostics,
) types.Object {
	if gateway == nil {
		return types.ObjectNull(vpcInterfaceNATGatewayAttrTypes)
	}

	addresses, addressesDiags := types.ListValueFrom(ctx, types.StringType, gateway.Addresses)
	diags.Append(addressesDiags...)

	portsets := types.ListNull(types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortsetAttrTypes})
	// Portsets are only returned for some customers, so keep them null when absent.
	if gateway.Portsets != nil {
		portsetModels := make([]VPCInterfaceNATGatewayPortsetModel, len(gateway.Portsets))
		for i, portset := range gateway.Portsets {
			ports := types.ListNull(types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortAttrTypes})
			if portset.Ports != nil {
				portModels := make([]VPCInterfaceNATGatewayPortModel, len(portset.Ports))
				for j, port := range portset.Ports {
					portModels[j] = VPCInterfaceNATGatewayPortModel{
						Start: types.Int64Value(int64(port.Start)),
						End:   types.Int64Value(int64(port.End)),
					}
				}

				var portDiags diag.Diagnostics
				ports, portDiags = types.ListValueFrom(
					ctx, types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortAttrTypes}, portModels,
				)
				diags.Append(portDiags...)
			}

			portsetModels[i] = VPCInterfaceNATGatewayPortsetModel{
				Address: types.StringValue(portset.Address),
				Ports:   ports,
			}
		}

		var portsetDiags diag.Diagnostics
		portsets, portsetDiags = types.ListValueFrom(
			ctx, types.ObjectType{AttrTypes: vpcInterfaceNATGatewayPortsetAttrTypes}, portsetModels,
		)
		diags.Append(portsetDiags...)
	}

	obj, objDiags := types.ObjectValueFrom(ctx, vpcInterfaceNATGatewayAttrTypes, VPCInterfaceNATGatewayModel{
		ID:                 types.Int64Value(int64(gateway.ID)),
		Label:              types.StringValue(gateway.Label),
		Addresses:          addresses,
		PortsetAssignments: types.Int64Value(int64(gateway.PortsetAssignments)),
		PortsetCapacity:    types.Int64Value(int64(gateway.PortsetCapacity)),
		Portsets:           portsets,
	})
	diags.Append(objDiags...)

	return obj
}

// VPCIPv4AddressAttrModel is a shared model between `configuredVPCInterfaceIPv4Address` and
// `computedVPCInterfaceIPv4Address`
type VPCIPv4AddressAttrModel struct {
	Address      types.String `tfsdk:"address"`
	Primary      types.Bool   `tfsdk:"primary"`
	Nat11Address types.String `tfsdk:"nat_1_1_address"`
}

type ConfiguredVPCIPv4RangeAttrModel struct {
	Range types.String `tfsdk:"range"`
}

type ComputedVPCIPv4RangeAttrModel struct {
	Range cidrtypes.IPv4Prefix `tfsdk:"range"`
}

type VPCIPv6AttrModel struct {
	IsPublic       types.Bool `tfsdk:"is_public"`
	SLAAC          types.List `tfsdk:"slaac"`
	AssignedSLAAC  types.Set  `tfsdk:"assigned_slaac"`
	Ranges         types.List `tfsdk:"ranges"`
	AssignedRanges types.Set  `tfsdk:"assigned_ranges"`
}

type VPCIPv6SLAACAttrModel struct {
	Range types.String `tfsdk:"range"`
}

type VPCIPv6SLAACAttrComputedModel struct {
	Range   cidrtypes.IPv6Prefix `tfsdk:"range"`
	Address iptypes.IPv6Address  `tfsdk:"address"`
}

type ConfiguredVPCIPv6RangeAttrModel struct {
	Range types.String `tfsdk:"range"`
}

type ComputedVPCIPv6RangeAttrModel struct {
	Range cidrtypes.IPv6Prefix `tfsdk:"range"`
}

func (plan *VPCAttrModel) GetCreateOptions(ctx context.Context, diags *diag.Diagnostics) (opts linodego.VPCInterfaceCreateOptions) {
	tflog.Trace(ctx, "Enter VPCAttrModel.GetCreateOptions")

	opts.SubnetID = helper.FrameworkSafeInt64ToInt(plan.SubnetID.ValueInt64(), diags)

	if !plan.IPv4.IsUnknown() && !plan.IPv4.IsNull() {
		var planIPv4 VPCIPv4AttrModel
		diags.Append(plan.IPv4.As(ctx, &planIPv4, basetypes.ObjectAsOptions{})...)
		ipv4Opts, _ := planIPv4.GetCreateOrUpdateOptions(ctx, nil, diags)
		opts.IPv4 = &ipv4Opts
	}

	if !plan.IPv6.IsUnknown() && !plan.IPv6.IsNull() {
		var planIPv6 VPCIPv6AttrModel
		diags.Append(plan.IPv6.As(ctx, &planIPv6, basetypes.ObjectAsOptions{})...)
		ipv6Opts, _ := planIPv6.GetCreateOrUpdateOptions(ctx, nil, diags)
		opts.IPv6 = &ipv6Opts
	}

	return opts
}

func (plan *VPCAttrModel) GetUpdateOptions(
	ctx context.Context,
	state *VPCAttrModel,
	diags *diag.Diagnostics,
) (opts linodego.VPCInterfaceUpdateOptions, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter VPCAttrModel.GetUpdateOptions")

	// Note: SubnetID cannot be updated according to the API, so we don't include it

	if !plan.IPv4.IsUnknown() && !plan.IPv4.IsNull() {
		var planIPv4 VPCIPv4AttrModel
		diags.Append(plan.IPv4.As(ctx, &planIPv4, basetypes.ObjectAsOptions{})...)

		var stateIPv4 *VPCIPv4AttrModel
		if state != nil && !state.IPv4.IsNull() {
			diags.Append(state.IPv4.As(ctx, &stateIPv4, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return opts, shouldUpdate
			}
		}

		if ipv4Opts, ipv4ShouldUpdate := planIPv4.GetCreateOrUpdateOptions(ctx, stateIPv4, diags); ipv4ShouldUpdate {
			opts.IPv4 = &ipv4Opts
			shouldUpdate = true
		}
	}

	if !plan.IPv6.IsUnknown() && !plan.IPv6.IsNull() {
		var planIPv6 VPCIPv6AttrModel
		diags.Append(plan.IPv6.As(ctx, &planIPv6, basetypes.ObjectAsOptions{})...)

		var stateIPv6 *VPCIPv6AttrModel
		if state != nil && !state.IPv6.IsNull() {
			diags.Append(state.IPv6.As(ctx, &stateIPv6, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return opts, shouldUpdate
			}
		}

		if ipv6Opts, ipv6ShouldUpdate := planIPv6.GetCreateOrUpdateOptions(ctx, stateIPv6, diags); ipv6ShouldUpdate {
			opts.IPv6 = &ipv6Opts
			shouldUpdate = true
		}
	}

	return opts, shouldUpdate
}

func (plan *VPCIPv4AttrModel) GetCreateOrUpdateOptions(
	ctx context.Context,
	state *VPCIPv4AttrModel,
	diags *diag.Diagnostics,
) (opts linodego.VPCInterfaceIPv4CreateOptions, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter VPCIPv4AttrModel.GetCreateOrUpdateOptions")
	if !plan.Addresses.IsUnknown() && !plan.Addresses.IsNull() && (state == nil || !state.Addresses.Equal(plan.Addresses)) {
		length := len(plan.Addresses.Elements())
		addresses := make([]VPCIPv4AddressAttrModel, 0, length)
		diags.Append(plan.Addresses.ElementsAs(ctx, &addresses, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		addressOpts := make([]linodego.VPCInterfaceIPv4AddressCreateOptions, len(addresses))
		for i, address := range addresses {
			addressOpts[i] = address.GetCreateOptions(ctx)
		}
		opts.Addresses = addressOpts
		shouldUpdate = true
	}

	if !plan.Ranges.IsUnknown() && !plan.Ranges.IsNull() && (state == nil || !state.Ranges.Equal(plan.Ranges)) {
		length := len(plan.Ranges.Elements())
		ranges := make([]ConfiguredVPCIPv4RangeAttrModel, 0, length)
		diags.Append(plan.Ranges.ElementsAs(ctx, &ranges, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		rangeOpts := make([]linodego.VPCInterfaceIPv4RangeCreateOptions, len(ranges))
		for i, r := range ranges {
			rangeOpts[i] = r.GetCreateOptions(ctx)
		}
		opts.Ranges = rangeOpts
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func (plan *VPCIPv4AddressAttrModel) GetCreateOptions(ctx context.Context) linodego.VPCInterfaceIPv4AddressCreateOptions {
	tflog.Trace(ctx, "Enter VPCIPv4AddressAttrModel.GetCreateOptions")

	opts := linodego.VPCInterfaceIPv4AddressCreateOptions{}

	if !plan.Address.IsUnknown() {
		opts.Address = plan.Address.ValueStringPointer()
	}

	if !plan.Primary.IsUnknown() {
		opts.Primary = plan.Primary.ValueBoolPointer()
	}

	if !plan.Nat11Address.IsUnknown() {
		opts.NAT1To1Address = plan.Nat11Address.ValueStringPointer()
	}

	return opts
}

func (plan *ConfiguredVPCIPv4RangeAttrModel) GetCreateOptions(ctx context.Context) linodego.VPCInterfaceIPv4RangeCreateOptions {
	tflog.Trace(ctx, "Enter ConfiguredVPCIPv4RangeAttrModel.GetCreateOptions")

	return linodego.VPCInterfaceIPv4RangeCreateOptions{
		Range: plan.Range.ValueString(),
	}
}

func (data *VPCAttrModel) FlattenVPCInterface(
	ctx context.Context, vpcInterface linodego.VPCInterface, preserveKnown bool, diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter VPCAttrModel.FlattenVPCInterface")

	data.SubnetID = helper.KeepOrUpdateInt64(data.SubnetID, int64(vpcInterface.SubnetID), preserveKnown)

	flattenedIPv4 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv4, resourceVPCIPv4Attribute.GetType().(basetypes.ObjectType).AttrTypes, preserveKnown, diags,
		func(ipv4 *VPCIPv4AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			ipv4.FlattenVPCIPv4(ctx, vpcInterface.IPv4, pk, d)
		},
	)

	if diags.HasError() {
		return
	}

	data.IPv4 = *flattenedIPv4

	flattenedIPv6 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv6, resourceVPCIPv6Attribute.GetType().(basetypes.ObjectType).AttrTypes, preserveKnown, diags,
		func(ipv6 *VPCIPv6AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			ipv6.FlattenVPCIPv6(ctx, vpcInterface.IPv6, pk, d)
		},
	)

	if diags.HasError() {
		return
	}

	data.IPv6 = *flattenedIPv6
}

func (data *VPCIPv4AttrModel) FlattenVPCIPv4(ctx context.Context, ipv4 linodego.VPCInterfaceIPv4, preserveKnown bool, diags *diag.Diagnostics) {
	tflog.Trace(ctx, "Enter VPCIPv4AttrModel.FlattenVPCIPv4")

	// When the object is null/unknown, the types of attributes of the object won't be filled by object.As(...) in the
	// helper function `KeepOrUpdateSingleNestedAttributeWithTypes`, so resetting manually here.
	if data.Addresses.IsNull() {
		data.Addresses = types.ListNull(configuredVPCInterfaceIPv4Address.Type())
	}
	if data.Ranges.IsNull() {
		data.Ranges = types.ListNull(configuredVPCInterfaceIPv4Range.Type())
	}

	assignedAddresses := make([]VPCIPv4AddressAttrModel, len(ipv4.Addresses))
	for i, addr := range ipv4.Addresses {
		assignedAddresses[i] = VPCIPv4AddressAttrModel{
			Address:      types.StringValue(addr.Address),
			Primary:      types.BoolValue(addr.Primary),
			Nat11Address: types.StringPointerValue(addr.NAT1To1Address),
		}
	}

	assignedAddressesValue, assignedAddressesDiags := types.SetValueFrom(
		ctx, computedVPCInterfaceIPv4Address.GetAttributes().Type(), assignedAddresses,
	)
	diags.Append(assignedAddressesDiags...)
	if diags.HasError() {
		return
	}

	data.AssignedAddresses = helper.KeepOrUpdateValue(data.AssignedAddresses, assignedAddressesValue, preserveKnown)

	assignedRanges := make([]ComputedVPCIPv4RangeAttrModel, len(ipv4.Ranges))
	for i, r := range ipv4.Ranges {
		assignedRanges[i] = ComputedVPCIPv4RangeAttrModel{
			Range: cidrtypes.NewIPv4PrefixValue(r.Range),
		}
	}

	assignedRangesValue, assignedRangesDiags := types.SetValueFrom(
		ctx, computedVPCInterfaceIPv4Range.GetAttributes().Type(), assignedRanges,
	)
	diags.Append(assignedRangesDiags...)
	if diags.HasError() {
		return
	}

	data.AssignedRanges = helper.KeepOrUpdateValue(data.AssignedRanges, assignedRangesValue, preserveKnown)

	natGateway := flattenVPCInterfaceNATGateway(ctx, ipv4.NATGateway, diags)
	if diags.HasError() {
		return
	}
	data.NATGateway = helper.KeepOrUpdateValue(data.NATGateway, natGateway, preserveKnown)
}

func (plan *VPCIPv6AttrModel) GetCreateOrUpdateOptions(
	ctx context.Context,
	state *VPCIPv6AttrModel,
	diags *diag.Diagnostics,
) (opts linodego.VPCInterfaceIPv6CreateOptions, shouldUpdate bool) {
	if !plan.IsPublic.IsUnknown() &&
		!plan.IsPublic.IsNull() && (state == nil || !state.IsPublic.Equal(plan.IsPublic)) {
		opts.IsPublic = plan.IsPublic.ValueBoolPointer()
		shouldUpdate = true
	}

	if !plan.SLAAC.IsUnknown() && !plan.SLAAC.IsNull() && (state == nil || !state.SLAAC.Equal(plan.SLAAC)) {
		length := len(plan.SLAAC.Elements())
		slaac := make([]VPCIPv6SLAACAttrModel, 0, length)
		diags.Append(plan.SLAAC.ElementsAs(ctx, &slaac, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		slaacOpts := helper.MapSlice(
			slaac,
			func(entry VPCIPv6SLAACAttrModel) linodego.VPCInterfaceIPv6SLAACCreateOptions {
				return entry.GetCreateOptions()
			},
		)
		opts.SLAAC = slaacOpts
		shouldUpdate = true
	}

	if !plan.Ranges.IsUnknown() && !plan.Ranges.IsNull() && (state == nil || !state.Ranges.Equal(plan.Ranges)) {
		length := len(plan.Ranges.Elements())
		ranges := make([]ConfiguredVPCIPv6RangeAttrModel, 0, length)
		diags.Append(plan.Ranges.ElementsAs(ctx, &ranges, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		rangeOpts := make([]linodego.VPCInterfaceIPv6RangeCreateOptions, len(ranges))
		for i, r := range ranges {
			rangeOpts[i] = r.GetCreateOptions()
		}
		opts.Ranges = rangeOpts
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func (plan *VPCIPv6SLAACAttrModel) GetCreateOptions() linodego.VPCInterfaceIPv6SLAACCreateOptions {
	opts := linodego.VPCInterfaceIPv6SLAACCreateOptions{}

	if !plan.Range.IsUnknown() {
		opts.Range = plan.Range.ValueString()
	}

	return opts
}

func (plan *ConfiguredVPCIPv6RangeAttrModel) GetCreateOptions() linodego.VPCInterfaceIPv6RangeCreateOptions {
	return linodego.VPCInterfaceIPv6RangeCreateOptions{
		Range: plan.Range.ValueString(),
	}
}

func (data *VPCIPv6AttrModel) FlattenVPCIPv6(ctx context.Context, ipv6 linodego.VPCInterfaceIPv6, preserveKnown bool, diags *diag.Diagnostics) {
	data.IsPublic = helper.KeepOrUpdateValue(data.IsPublic, types.BoolPointerValue(ipv6.IsPublic), preserveKnown)

	// When the object is null/unknown, the types of attributes of the object won't be filled by object.As(...) in the
	// helper function `KeepOrUpdateSingleNestedAttributeWithTypes`, so resetting manually here.
	if data.SLAAC.IsNull() {
		data.SLAAC = types.ListNull(configuredVPCInterfaceIPv6SLAAC.Type())
	}
	if data.Ranges.IsNull() {
		data.Ranges = types.ListNull(configuredVPCInterfaceIPv6Range.Type())
	}

	assignedSLAAC := helper.MapSlice(
		ipv6.SLAAC,
		func(slaac linodego.VPCInterfaceIPv6SLAAC) VPCIPv6SLAACAttrComputedModel {
			return VPCIPv6SLAACAttrComputedModel{
				Range:   cidrtypes.NewIPv6PrefixValue(slaac.Range),
				Address: iptypes.NewIPv6AddressValue(slaac.Address),
			}
		},
	)

	assignedSLAACValue, assignedSLAACDiags := types.SetValueFrom(
		ctx, computedVPCInterfaceIPv6SLAAC.GetAttributes().Type(), assignedSLAAC,
	)
	diags.Append(assignedSLAACDiags...)
	if diags.HasError() {
		return
	}

	data.AssignedSLAAC = helper.KeepOrUpdateValue(data.AssignedSLAAC, assignedSLAACValue, preserveKnown)

	assignedRanges := helper.MapSlice(
		ipv6.Ranges,
		func(r linodego.VPCInterfaceIPv6Range) ComputedVPCIPv6RangeAttrModel {
			return ComputedVPCIPv6RangeAttrModel{
				Range: cidrtypes.NewIPv6PrefixValue(r.Range),
			}
		},
	)

	assignedRangesValue, assignedRangesDiags := types.SetValueFrom(
		ctx, computedVPCInterfaceIPv6Range.GetAttributes().Type(), assignedRanges,
	)
	diags.Append(assignedRangesDiags...)
	if diags.HasError() {
		return
	}

	data.AssignedRanges = helper.KeepOrUpdateValue(data.AssignedRanges, assignedRangesValue, preserveKnown)
}
