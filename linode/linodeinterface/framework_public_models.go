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

type PublicAttrModel struct {
	IPv4 types.Object `tfsdk:"ipv4"`
	IPv6 types.Object `tfsdk:"ipv6"`
}

type PublicIPv4AttrModel struct {
	Addresses         types.List `tfsdk:"addresses"`
	AssignedAddresses types.Set  `tfsdk:"assigned_addresses"`
	Shared            types.Set  `tfsdk:"shared"`
}

type PublicIPv6AttrModel struct {
	Ranges         types.List `tfsdk:"ranges"`
	AssignedRanges types.Set  `tfsdk:"assigned_ranges"`
	Shared         types.Set  `tfsdk:"shared"`
	SLAAC          types.Set  `tfsdk:"slaac"`
}

type SharedPublicIPv4AddressAttrModel struct {
	Address  types.String `tfsdk:"address"`
	LinodeID types.Int64  `tfsdk:"linode_id"`
}

// PublicIPv4AddressAttrModel is a shared model between `configuredPublicInterfaceIPv4Address` and
// `computedPublicInterfaceIPv4Address` schemas.
type PublicIPv4AddressAttrModel struct {
	Address types.String `tfsdk:"address"`
	Primary types.Bool   `tfsdk:"primary"`
}

type ConfiguredPublicIPv6RangeAttrModel struct {
	Range types.String `tfsdk:"range"`
}

type ComputedPublicIPv6RangeAttrModel struct {
	Range       types.String `tfsdk:"range"`
	RouteTarget types.String `tfsdk:"route_target"`
}

type PublicIPv6SLAACAttrModel struct {
	Address types.String `tfsdk:"address"`
	Prefix  types.Int64  `tfsdk:"prefix"`
}

func (plan *PublicAttrModel) GetCreateOrUpdateOptions(
	ctx context.Context,
	state *PublicAttrModel,
) (opts linodego.PublicInterfaceCreateOptions, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter PublicAttrModel.GetCreateOrUpdateOptions")

	if !plan.IPv4.IsUnknown() && !plan.IPv4.IsNull() && (state == nil || !state.IPv4.Equal(plan.IPv4)) {
		var planPublicIPv4 PublicIPv4AttrModel
		plan.IPv4.As(ctx, &planPublicIPv4, basetypes.ObjectAsOptions{})
		opts.IPv4 = linodego.Pointer(planPublicIPv4.GetCreateOptions(ctx))
		shouldUpdate = true
	}

	if !plan.IPv6.IsUnknown() && !plan.IPv6.IsNull() && (state == nil || !state.IPv6.Equal(plan.IPv6)) {
		var planPublicIPv6 PublicIPv6AttrModel
		plan.IPv6.As(ctx, &planPublicIPv6, basetypes.ObjectAsOptions{})
		opts.IPv6 = linodego.Pointer(planPublicIPv6.GetCreateOptions(ctx))
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func (plan *PublicIPv4AttrModel) GetCreateOptions(ctx context.Context) (opts linodego.PublicInterfaceIPv4CreateOptions) {
	tflog.Trace(ctx, "Enter PublicIPv4AttrModel.GetCreateOptions")

	if !plan.Addresses.IsNull() && !plan.Addresses.IsUnknown() {
		length := len(plan.Addresses.Elements())
		addressesOpts := make([]linodego.PublicInterfaceIPv4AddressCreateOptions, 0, length)

		addresses := make([]PublicIPv4AddressAttrModel, 0, length)

		// Since `addresses` list is with a default value of empty list, we may safely assume
		// its elements won't contain any unknown and null
		plan.Addresses.ElementsAs(ctx, &addresses, false)

		for _, v := range addresses {
			addressesOpts = append(addressesOpts, v.GetCreateOptions(ctx))
		}
		opts.Addresses = linodego.Pointer(addressesOpts)
	}

	return opts
}

func (data *PublicIPv4AttrModel) FlattenPublicIPv4(ctx context.Context, ipv4 linodego.PublicInterfaceIPv4, preserveKnown bool, diags *diag.Diagnostics) {
	tflog.Trace(ctx, "Enter PublicIPv4AttrModel.FlattenPublicIPv4")

	// data.Address should never need to be flattened from a linodego struct because its values can
	// either be configured by the TF practitioner or defaulted to an empty list

	// when it's null, the types of attributes of the object won't be filled by object.As(...), resetting manually here
	// TODO: filing a bugfix/feature request to HashiCorp
	if data.Addresses.IsNull() {
		data.Addresses = types.ListNull(configuredPublicInterfaceIPv4Address.Type())
	}

	var newDiags diag.Diagnostics
	assignedAddresses := make([]PublicIPv4AddressAttrModel, len(ipv4.Addresses))
	for i, v := range ipv4.Addresses {
		assignedAddresses[i] = PublicIPv4AddressAttrModel{
			Address: types.StringValue(v.Address),
			Primary: types.BoolValue(v.Primary),
		}
	}

	// Each object in the `assigned_addresses` set attribute is computed-only without UseStateForUnknown plan modifier,
	// so it's always an unknown. Thus, no need to check the nested attributes of these objects
	newAssignedAddresses, newDiags := types.SetValueFrom(ctx, computedPublicInterfaceIPv4Address.Type(), assignedAddresses)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.AssignedAddresses = helper.KeepOrUpdateValue(data.AssignedAddresses, newAssignedAddresses, preserveKnown)

	shared := make([]SharedPublicIPv4AddressAttrModel, len(ipv4.Shared))
	for i, v := range ipv4.Shared {
		shared[i] = SharedPublicIPv4AddressAttrModel{
			Address:  types.StringValue(v.Address),
			LinodeID: types.Int64Value(int64(v.LinodeID)),
		}
	}

	// Each object in the `shared` set attribute is computed-only without UseStateForUnknown plan modifier,
	// so it's always an unknown. Thus, no need to check the nested attributes of these objects.
	newShared, newDiags := types.SetValueFrom(ctx, sharedPublicInterfaceIPv4Address.Type(), shared)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	data.Shared = helper.KeepOrUpdateValue(data.Shared, newShared, preserveKnown)
}

func (data *PublicIPv6AttrModel) FlattenPublicIPv6(ctx context.Context, ipv6 linodego.PublicInterfaceIPv6, preserveKnown bool, diags *diag.Diagnostics) {
	tflog.Trace(ctx, "Enter PublicIPv6AttrModel.FlattenPublicIPv6")

	// data.Ranges should never need to be flattened from a linodego struct because its values can
	// either be configured by the TF practitioner or defaulted to an empty list

	// when it's null, the types of attributes of the object won't be filled by object.As(...), resetting manually here
	// TODO: filing a bugfix/feature request to HashiCorp
	if data.Ranges.IsNull() {
		data.Ranges = types.ListNull(configuredPublicInterfaceIPv6Range.Type())
	}

	var newDiags diag.Diagnostics
	assignedRanges := make([]ComputedPublicIPv6RangeAttrModel, len(ipv6.Ranges))
	for i, v := range ipv6.Ranges {
		assignedRanges[i] = ComputedPublicIPv6RangeAttrModel{
			Range:       types.StringValue(v.Range),
			RouteTarget: types.StringPointerValue(v.RouteTarget),
		}
	}

	// `assigned_ranges` attribute is computed-only so it's always an unknown when the ipv6 is being flattened
	data.AssignedRanges, newDiags = types.SetValueFrom(ctx, computedPublicInterfaceIPv6Range.Type(), assignedRanges)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	shared := make([]ComputedPublicIPv6RangeAttrModel, len(ipv6.Shared))
	for i, v := range ipv6.Shared {
		shared[i] = ComputedPublicIPv6RangeAttrModel{
			Range:       types.StringValue(v.Range),
			RouteTarget: types.StringPointerValue(v.RouteTarget),
		}
	}

	// `shared` attribute is computed-only so it's always an unknown when the ipv6 is being flattened
	data.Shared, newDiags = types.SetValueFrom(ctx, computedPublicInterfaceIPv6Range.Type(), shared)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	slaac := make([]PublicIPv6SLAACAttrModel, len(ipv6.SLAAC))
	for i, v := range ipv6.SLAAC {
		slaac[i] = PublicIPv6SLAACAttrModel{
			Address: types.StringValue(v.Address),
			Prefix:  types.Int64Value(int64(v.Prefix)),
		}
	}

	// `slaac` attribute is computed-only so it's always an unknown when the ipv6 is being flattened
	data.SLAAC, newDiags = types.SetValueFrom(ctx, resourcePublicInterfaceIPv6SLAAC.Type(), slaac)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}
}

func (plan *PublicIPv6AttrModel) GetCreateOptions(ctx context.Context) (opts linodego.PublicInterfaceIPv6CreateOptions) {
	tflog.Trace(ctx, "Enter PublicIPv6AttrModel.GetCreateOptions")

	if !plan.Ranges.IsNull() && !plan.Ranges.IsUnknown() {
		length := len(plan.Ranges.Elements())

		rangesOpts := make([]linodego.PublicInterfaceIPv6RangeCreateOptions, 0, length)
		ranges := make([]ConfiguredPublicIPv6RangeAttrModel, 0, length)

		// Since `ranges` list is with a default value of empty list, we may safely assume
		// its elements won't contain any unknown and null
		plan.Ranges.ElementsAs(ctx, &ranges, false)

		for _, v := range ranges {
			rangesOpts = append(rangesOpts, v.GetCreateOptions(ctx))
		}
		opts.Ranges = linodego.Pointer(rangesOpts)
	}

	return opts
}

func (plan *PublicIPv4AddressAttrModel) GetCreateOptions(ctx context.Context) (opts linodego.PublicInterfaceIPv4AddressCreateOptions) {
	tflog.Trace(ctx, "Enter PublicIPv4AddressAttrModel.GetCreateOptions")

	opts.Address = helper.ValueStringPointerWithUnknownToNil(plan.Address)
	opts.Primary = helper.ValueBoolPointerWithUnknownToNil(plan.Primary)
	return opts
}

func (plan *ConfiguredPublicIPv6RangeAttrModel) GetCreateOptions(ctx context.Context) (opts linodego.PublicInterfaceIPv6RangeCreateOptions) {
	tflog.Trace(ctx, "Enter ConfiguredPublicIPv6RangeAttrModel.GetCreateOptions")

	opts.Range = plan.Range.ValueString()
	return opts
}

func (data *PublicAttrModel) FlattenPublicInterface(
	ctx context.Context, publicInterface linodego.PublicInterface, preserveKnown bool, diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter PublicAttrModel.FlattenPublicInterface")

	flattenedPublicIPv4 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv4, resourcePublicIPv4Attribute.GetType().(types.ObjectType).AttrTypes, preserveKnown, diags,
		func(publicIPv4 *PublicIPv4AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			if publicInterface.IPv4 == nil {
				*isNull = true
				publicIPv4.Addresses = helper.KeepOrUpdateValue(
					publicIPv4.Addresses, types.ListNull(configuredPublicInterfaceIPv4Address.GetAttributes().Type()), pk,
				)
				publicIPv4.AssignedAddresses = helper.KeepOrUpdateValue(
					publicIPv4.AssignedAddresses, types.SetNull(computedPublicInterfaceIPv4Address.GetAttributes().Type()), pk,
				)
				publicIPv4.Shared = helper.KeepOrUpdateValue(
					publicIPv4.Shared, types.SetNull(sharedPublicInterfaceIPv4Address.GetAttributes().Type()), pk,
				)
				return
			}

			publicIPv4.FlattenPublicIPv4(ctx, *publicInterface.IPv4, pk, d)
		},
	)
	if diags.HasError() {
		return
	}

	data.IPv4 = *flattenedPublicIPv4

	flattenedPublicIPv6 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.IPv6, resourcePublicIPv6Attribute.GetType().(types.ObjectType).AttrTypes, preserveKnown, diags,
		func(publicIPv6 *PublicIPv6AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			if publicInterface.IPv6 == nil {
				*isNull = true
				publicIPv6.Ranges = helper.KeepOrUpdateValue(
					publicIPv6.Ranges, types.ListNull(configuredPublicInterfaceIPv6Range.GetAttributes().Type()), pk,
				)
				publicIPv6.AssignedRanges = helper.KeepOrUpdateValue(
					publicIPv6.AssignedRanges, types.SetNull(computedPublicInterfaceIPv6Range.GetAttributes().Type()), pk,
				)
				publicIPv6.Shared = helper.KeepOrUpdateValue(
					publicIPv6.Shared, types.SetNull(computedPublicInterfaceIPv6Range.GetAttributes().Type()), pk,
				)
				publicIPv6.SLAAC = helper.KeepOrUpdateValue(
					publicIPv6.SLAAC, types.SetNull(resourcePublicInterfaceIPv6SLAAC.GetAttributes().Type()), pk,
				)
				return
			}
			publicIPv6.FlattenPublicIPv6(ctx, *publicInterface.IPv6, pk, d)
		},
	)
	if diags.HasError() {
		return
	}

	data.IPv6 = *flattenedPublicIPv6
}
