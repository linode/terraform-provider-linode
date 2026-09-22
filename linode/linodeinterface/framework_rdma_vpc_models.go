package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/customtypes"
)

// RDMAVPCAttrModel is the model for an RDMA VPC interface (resource).
type RDMAVPCAttrModel struct {
	VPCID    types.Int64  `tfsdk:"vpc_id"`
	SubnetID types.Int64  `tfsdk:"subnet_id"`
	IPv4     types.Object `tfsdk:"ipv4"`
}

// RDMAVPCIPv4AttrModel is the model for the IPv4 block of an RDMA VPC interface.
type RDMAVPCIPv4AttrModel struct {
	Addresses types.List `tfsdk:"addresses"`
}

// RDMAVPCIPv4AddressAttrModel represents a single IPv4 address on an RDMA VPC interface.
type RDMAVPCIPv4AddressAttrModel struct {
	Address customtypes.LinodeAutoAllocIPValue `tfsdk:"address"`
	Primary types.Bool                         `tfsdk:"primary"`
}

// GetUpdateOptions returns the linodego update options for an RDMA VPC interface.
func (plan *RDMAVPCAttrModel) GetUpdateOptions(
	ctx context.Context,
	state *RDMAVPCAttrModel,
	diags *diag.Diagnostics,
) (opts linodego.RDMAVPCInterfaceUpdateOptions, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter RDMAVPCAttrModel.GetUpdateOptions")

	if !plan.SubnetID.IsUnknown() && !plan.SubnetID.IsNull() &&
		(state == nil || !state.SubnetID.Equal(plan.SubnetID)) {
		subnetID := helper.FrameworkSafeInt64ToInt(plan.SubnetID.ValueInt64(), diags)
		opts.SubnetID = &subnetID
		shouldUpdate = true
	}

	if !plan.IPv4.IsUnknown() && !plan.IPv4.IsNull() {
		var planIPv4 RDMAVPCIPv4AttrModel
		diags.Append(plan.IPv4.As(ctx, &planIPv4, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return opts, shouldUpdate
		}

		var stateIPv4 *RDMAVPCIPv4AttrModel
		if state != nil && !state.IPv4.IsNull() {
			diags.Append(state.IPv4.As(ctx, &stateIPv4, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return opts, shouldUpdate
			}
		}

		ipv4Opts, ipv4ShouldUpdate := planIPv4.GetUpdateOptions(ctx, stateIPv4, diags)
		if ipv4ShouldUpdate {
			opts.IPv4 = &ipv4Opts
			shouldUpdate = true
		}
	}

	return opts, shouldUpdate
}

// GetUpdateOptions returns the linodego IPv4 options for an RDMA VPC interface.
func (plan *RDMAVPCIPv4AttrModel) GetUpdateOptions(
	ctx context.Context,
	state *RDMAVPCIPv4AttrModel,
	diags *diag.Diagnostics,
) (opts linodego.RDMAVPCInterfaceIPv4Options, shouldUpdate bool) {
	if plan.Addresses.IsUnknown() || plan.Addresses.IsNull() {
		return opts, shouldUpdate
	}

	if state != nil && state.Addresses.Equal(plan.Addresses) {
		return opts, shouldUpdate
	}

	addresses := make([]RDMAVPCIPv4AddressAttrModel, 0, len(plan.Addresses.Elements()))
	diags.Append(plan.Addresses.ElementsAs(ctx, &addresses, false)...)
	if diags.HasError() {
		return opts, shouldUpdate
	}

	addressOpts := make([]linodego.RDMAVPCInterfaceIPv4AddressOptions, len(addresses))
	for i, a := range addresses {
		addressOpts[i] = linodego.RDMAVPCInterfaceIPv4AddressOptions{
			Address: a.Address.ValueString(),
			Primary: a.Primary.ValueBoolPointer(),
		}
	}
	opts.Addresses = addressOpts
	shouldUpdate = true
	return opts, shouldUpdate
}

// FlattenRDMAVPCInterface flattens an RDMA VPC interface into the model.
func (data *RDMAVPCAttrModel) FlattenRDMAVPCInterface(
	ctx context.Context,
	rdma linodego.RDMAVPCInterface,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	tflog.Trace(ctx, "Enter RDMAVPCAttrModel.FlattenRDMAVPCInterface")

	data.VPCID = helper.KeepOrUpdateInt64(data.VPCID, int64(rdma.VPCID), preserveKnown)
	data.SubnetID = helper.KeepOrUpdateInt64(data.SubnetID, int64(rdma.SubnetID), preserveKnown)

	flattenedIPv4 := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx,
		data.IPv4,
		resourceRDMAVPCIPv4Attribute.GetType().(basetypes.ObjectType).AttrTypes,
		preserveKnown,
		diags,
		func(ipv4 *RDMAVPCIPv4AttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			ipv4.FlattenIPv4(ctx, rdma.IPv4, pk, d)
		},
	)
	if diags.HasError() {
		return
	}
	if flattenedIPv4 == nil {
		return
	}
	data.IPv4 = *flattenedIPv4
}

// FlattenIPv4 flattens the IPv4 portion of an RDMA VPC interface.
func (data *RDMAVPCIPv4AttrModel) FlattenIPv4(
	ctx context.Context,
	ipv4 linodego.RDMAVPCInterfaceIPv4,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	addresses := helper.MapSlice(
		ipv4.Addresses,
		func(addr linodego.RDMAVPCInterfaceIPv4Address) RDMAVPCIPv4AddressAttrModel {
			return RDMAVPCIPv4AddressAttrModel{
				Address: customtypes.LinodeAutoAllocIPValueFrom(addr.Address),
				Primary: types.BoolValue(addr.Primary),
			}
		},
	)

	addressesValue, addressesDiags := types.ListValueFrom(
		ctx, configuredRDMAVPCInterfaceIPv4Address.Type(), addresses,
	)
	diags.Append(addressesDiags...)
	if diags.HasError() {
		return
	}

	// NOTE: preserveKnown is intentionally ignored (treated as false) for the
	// addresses list so the API-resolved IP address is always written to state,
	// even when the user configured "auto". The LinodeAutoAllocIPType semantic
	// equality bridges the configured "auto" and the resolved address, so this
	// does not produce a perpetual diff or an inconsistent-result error.
	data.Addresses = helper.KeepOrUpdateValue(data.Addresses, addressesValue, false)
}
