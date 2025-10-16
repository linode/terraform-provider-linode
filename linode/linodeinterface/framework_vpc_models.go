package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type VPCAttrModel struct {
	IPv4     types.Object `tfsdk:"ipv4"`
	IPv6     types.Object `tfsdk:"ipv6"`
	SubnetID types.Int64  `tfsdk:"subnet_id"`
}

type VPCIPv4AttrModel struct {
	Addresses         types.List `tfsdk:"addresses"`
	AssignedAddresses types.Set  `tfsdk:"assigned_addresses"`
	Ranges            types.List `tfsdk:"ranges"`
	AssignedRanges    types.Set  `tfsdk:"assigned_ranges"`
}

// VPCIPv4AddressAttrModel is a shared model between `configuredVPCInterfaceIPv4Address` and
// `computedVPCInterfaceIPv4Address`
type VPCIPv4AddressAttrModel struct {
	Address      types.String `tfsdk:"address"`
	Primary      types.Bool   `tfsdk:"primary"`
	Nat11Address types.String `tfsdk:"nat_1_1_address"`
}

// VPCIPv4RangeAttrModel is a shared model between `configuredVPCInterfaceIPv4Range` and
// `computedVPCInterfaceIPv4Range`
type VPCIPv4RangeAttrModel struct {
	Range types.String `tfsdk:"range"`
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
	Range   types.String `tfsdk:"range"`
	Address types.String `tfsdk:"address"`
}

type VPCIPv6RangeAttrModel struct {
	Range types.String `tfsdk:"range"`
}

func (plan *VPCAttrModel) GetCreateOptions(ctx context.Context, diags *diag.Diagnostics) (opts linodego.VPCInterfaceCreateOptions) {
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
	if !plan.Addresses.IsUnknown() && !plan.Addresses.IsNull() && (state == nil || !state.Addresses.Equal(plan.Addresses)) {
		length := len(plan.Addresses.Elements())
		addresses := make([]VPCIPv4AddressAttrModel, 0, length)
		diags.Append(plan.Addresses.ElementsAs(ctx, &addresses, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		addressOpts := make([]linodego.VPCInterfaceIPv4AddressCreateOptions, len(addresses))
		for i, address := range addresses {
			addressOpts[i] = address.GetCreateOptions()
		}
		opts.Addresses = &addressOpts
		shouldUpdate = true
	}

	if !plan.Ranges.IsUnknown() && !plan.Ranges.IsNull() && (state == nil || !state.Ranges.Equal(plan.Ranges)) {
		length := len(plan.Ranges.Elements())
		ranges := make([]VPCIPv4RangeAttrModel, 0, length)
		diags.Append(plan.Ranges.ElementsAs(ctx, &ranges, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		rangeOpts := make([]linodego.VPCInterfaceIPv4RangeCreateOptions, len(ranges))
		for i, r := range ranges {
			rangeOpts[i] = r.GetCreateOptions()
		}
		opts.Ranges = &rangeOpts
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func (plan *VPCIPv4AddressAttrModel) GetCreateOptions() linodego.VPCInterfaceIPv4AddressCreateOptions {
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

func (plan *VPCIPv4RangeAttrModel) GetCreateOptions() linodego.VPCInterfaceIPv4RangeCreateOptions {
	return linodego.VPCInterfaceIPv4RangeCreateOptions{
		Range: plan.Range.ValueString(),
	}
}

func (data *VPCAttrModel) FlattenVPCInterface(
	ctx context.Context, vpcInterface linodego.VPCInterface, preserveKnown bool, diags *diag.Diagnostics,
) {
	data.SubnetID = helper.KeepOrUpdateInt64(data.SubnetID, int64(vpcInterface.SubnetID), preserveKnown)

	flattenedIPv4 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv4, vpcIPv4Attribute.GetType().(basetypes.ObjectType).AttrTypes, preserveKnown, diags,
		func(ipv4 *VPCIPv4AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			ipv4.FlattenVPCIPv4(ctx, vpcInterface.IPv4, pk, d)
		},
	)

	if diags.HasError() {
		return
	}

	data.IPv4 = *flattenedIPv4

	flattenedIPv6 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv6, vpcIPv6Attribute.GetType().(basetypes.ObjectType).AttrTypes, preserveKnown, diags,
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

	assignedRanges := make([]VPCIPv4RangeAttrModel, len(ipv4.Ranges))
	for i, r := range ipv4.Ranges {
		assignedRanges[i] = VPCIPv4RangeAttrModel{
			Range: types.StringValue(r.Range),
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
		opts.SLAAC = &slaacOpts
		shouldUpdate = true
	}

	if !plan.Ranges.IsUnknown() && !plan.Ranges.IsNull() && (state == nil || !state.Ranges.Equal(plan.Ranges)) {
		length := len(plan.Ranges.Elements())
		ranges := make([]VPCIPv6RangeAttrModel, 0, length)
		diags.Append(plan.Ranges.ElementsAs(ctx, &ranges, false)...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		rangeOpts := make([]linodego.VPCInterfaceIPv6RangeCreateOptions, len(ranges))
		for i, r := range ranges {
			rangeOpts[i] = r.GetCreateOptions()
		}
		opts.Ranges = &rangeOpts
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

func (plan *VPCIPv6RangeAttrModel) GetCreateOptions() linodego.VPCInterfaceIPv6RangeCreateOptions {
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
				Range:   types.StringValue(slaac.Range),
				Address: types.StringValue(slaac.Address),
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
		func(r linodego.VPCInterfaceIPv6Range) VPCIPv6RangeAttrModel {
			return VPCIPv6RangeAttrModel{
				Range: types.StringValue(r.Range),
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
