package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type VPCAttrModel struct {
	IPv4     types.Object `tfsdk:"ipv4"`
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

func (plan *VPCAttrModel) GetCreateOptions(ctx context.Context, diags *diag.Diagnostics) (opts linodego.VPCInterfaceCreateOptions) {
	tflog.Trace(ctx, "Enter VPCAttrModel.GetCreateOptions")

	opts.SubnetID = helper.FrameworkSafeInt64ToInt(plan.SubnetID.ValueInt64(), diags)

	if !plan.IPv4.IsUnknown() && !plan.IPv4.IsNull() {
		var planIPv4 VPCIPv4AttrModel
		plan.IPv4.As(ctx, &planIPv4, basetypes.ObjectAsOptions{})
		ipv4Opts, _ := planIPv4.GetCreateOrUpdateOptions(ctx, nil)
		opts.IPv4 = &ipv4Opts
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
		plan.IPv4.As(ctx, &planIPv4, basetypes.ObjectAsOptions{})

		var stateIPv4 *VPCIPv4AttrModel
		if state != nil && !state.IPv4.IsNull() {
			state.IPv4.As(ctx, &stateIPv4, basetypes.ObjectAsOptions{})
		}

		if ipv4Opts, ipv4ShouldUpdate := planIPv4.GetCreateOrUpdateOptions(ctx, stateIPv4); ipv4ShouldUpdate {
			opts.IPv4 = &ipv4Opts
			shouldUpdate = true
		}
	}

	return opts, shouldUpdate
}

func (plan *VPCIPv4AttrModel) GetCreateOrUpdateOptions(
	ctx context.Context,
	state *VPCIPv4AttrModel,
) (opts linodego.VPCInterfaceIPv4CreateOptions, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter VPCIPv4AttrModel.GetCreateOrUpdateOptions")
	if !plan.Addresses.IsUnknown() && !plan.Addresses.IsNull() && (state == nil || !state.Addresses.Equal(plan.Addresses)) {
		length := len(plan.Addresses.Elements())
		addresses := make([]VPCIPv4AddressAttrModel, 0, length)
		plan.Addresses.ElementsAs(ctx, &addresses, false)

		addressOpts := make([]linodego.VPCInterfaceIPv4AddressCreateOptions, len(addresses))
		for i, address := range addresses {
			addressOpts[i] = address.GetCreateOptions(ctx)
		}
		opts.Addresses = &addressOpts
		shouldUpdate = true
	}

	if !plan.Ranges.IsUnknown() && !plan.Ranges.IsNull() && (state == nil || !state.Ranges.Equal(plan.Ranges)) {
		length := len(plan.Ranges.Elements())
		ranges := make([]VPCIPv4RangeAttrModel, 0, length)
		plan.Ranges.ElementsAs(ctx, &ranges, false)

		rangeOpts := make([]linodego.VPCInterfaceIPv4RangeCreateOptions, len(ranges))
		for i, r := range ranges {
			rangeOpts[i] = r.GetCreateOptions(ctx)
		}
		opts.Ranges = &rangeOpts
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

func (plan *VPCIPv4RangeAttrModel) GetCreateOptions(ctx context.Context) linodego.VPCInterfaceIPv4RangeCreateOptions {
	tflog.Trace(ctx, "Enter VPCIPv4RangeAttrModel.GetCreateOptions")

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
		ctx, data.IPv4, vpcIPv4Attribute.GetType().(basetypes.ObjectType).AttrTypes, preserveKnown, diags,
		func(ipv4 *VPCIPv4AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			ipv4.FlattenVPCIPv4(ctx, vpcInterface.IPv4, pk, d)
		},
	)

	if diags.HasError() {
		return
	}

	data.IPv4 = *flattenedIPv4
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
