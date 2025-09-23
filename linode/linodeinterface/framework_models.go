package linodeinterface

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type LinodeInterfaceModel struct {
	ID           types.String `tfsdk:"id"`
	LinodeID     types.Int64  `tfsdk:"linode_id"`
	FirewallID   types.Int64  `tfsdk:"firewall_id"`
	DefaultRoute types.Object `tfsdk:"default_route"`
	Public       types.Object `tfsdk:"public"`
	VLAN         types.Object `tfsdk:"vlan"`
	VPC          types.Object `tfsdk:"vpc"`
}

// Structs for Public Interfaces

// Structs for VLAN Interfaces

// Structs for VPC Interfaces

// type VPCAttrModel struct {
// 	IPv4 types.Bool `tfsdk:"ipv4"`
// 	IPv6 types.Bool `tfsdk:"ipv6"`
// }

func (state *LinodeInterfaceModel) GetIDs(diags *diag.Diagnostics) (linodeID int, id int) {
	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		diags.AddError(
			"Failed to Convert ID Type",
			fmt.Sprintf(
				"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					"Failed to convert string ID %q to an integer.\n", id,
			),
		)
	}
	linodeID = helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), diags)
	return linodeID, id
}

func (plan *LinodeInterfaceModel) GetCreateOptions(ctx context.Context, diags *diag.Diagnostics) (opts linodego.LinodeInterfaceCreateOptions, linodeID int) {
	if !plan.DefaultRoute.IsUnknown() && !plan.DefaultRoute.IsNull() {
		var planDefaultRoute DefaultRouteAttrModel
		plan.DefaultRoute.As(ctx, &planDefaultRoute, basetypes.ObjectAsOptions{})
		defaultRouteOpts, _ := planDefaultRoute.GetCreateOrUpdateOptions(nil)
		opts.DefaultRoute = linodego.Pointer(defaultRouteOpts)
	}

	if !plan.FirewallID.IsUnknown() {
		opts.FirewallID = helper.FrameworkSafeInt64ValueToIntDoublePointerWithUnknownToNil(plan.FirewallID, diags)
		if diags.HasError() {
			return opts, linodeID
		}
	}

	if !plan.Public.IsUnknown() && !plan.Public.IsNull() {
		var planPublicInterface PublicAttrModel
		plan.Public.As(ctx, &planPublicInterface, basetypes.ObjectAsOptions{})
		publicOpts, _ := planPublicInterface.GetCreateOrUpdateOptions(ctx, nil)
		opts.Public = linodego.Pointer(publicOpts)
	} else if !plan.VLAN.IsUnknown() && !plan.VLAN.IsNull() {
		var planVLANInterface VLANAttrModel
		plan.VLAN.As(ctx, &planVLANInterface, basetypes.ObjectAsOptions{})
		opts.VLAN = linodego.Pointer(planVLANInterface.GetCreateOptions())
	} else if !plan.VPC.IsUnknown() && !plan.VPC.IsNull() {
		var planVPCInterface VPCAttrModel
		plan.VPC.As(ctx, &planVPCInterface, basetypes.ObjectAsOptions{})
		vpc := planVPCInterface.GetCreateOptions(ctx, diags)
		opts.VPC = linodego.Pointer(vpc)
	}

	linodeID = helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), diags)
	return opts, linodeID
}

func (plan *LinodeInterfaceModel) GetUpdateOptions(
	ctx context.Context,
	state LinodeInterfaceModel,
	diags *diag.Diagnostics,
) (opts linodego.LinodeInterfaceUpdateOptions) {
	if !plan.DefaultRoute.IsUnknown() && !plan.DefaultRoute.IsNull() {
		var planDefaultRoute DefaultRouteAttrModel
		var stateDefaultRoute *DefaultRouteAttrModel
		plan.DefaultRoute.As(ctx, &planDefaultRoute, basetypes.ObjectAsOptions{})

		// state can't be unknown, checking null is enough here
		if !state.DefaultRoute.IsNull() {
			state.DefaultRoute.As(ctx, &stateDefaultRoute, basetypes.ObjectAsOptions{})
		}

		if updatedDefaultRoute, ok := planDefaultRoute.GetCreateOrUpdateOptions(stateDefaultRoute); ok {
			opts.DefaultRoute = linodego.Pointer(updatedDefaultRoute)
		}
	}

	if !plan.Public.IsUnknown() && !plan.Public.IsNull() {
		var planPublicInterface PublicAttrModel
		var statePublicInterface *PublicAttrModel
		plan.Public.As(ctx, &planPublicInterface, basetypes.ObjectAsOptions{})

		// state can't be unknown, checking null is enough here
		if !state.Public.IsNull() {
			state.Public.As(ctx, &statePublicInterface, basetypes.ObjectAsOptions{})
		}

		if updatedPublicInterface, shouldUpdate := planPublicInterface.GetCreateOrUpdateOptions(ctx, statePublicInterface); shouldUpdate {
			opts.Public = linodego.Pointer(updatedPublicInterface)
		}
	}

	if !plan.VPC.IsUnknown() && !plan.VPC.IsNull() {
		var planVPCInterface VPCAttrModel
		var stateVPCInterface *VPCAttrModel
		plan.VPC.As(ctx, &planVPCInterface, basetypes.ObjectAsOptions{})

		// state can't be unknown, checking null is enough here
		if !state.VPC.IsNull() {
			state.VPC.As(ctx, &stateVPCInterface, basetypes.ObjectAsOptions{})
		}

		if updatedVPCInterface, ok := planVPCInterface.GetUpdateOptions(ctx, stateVPCInterface, diags); ok {
			opts.VPC = linodego.Pointer(updatedVPCInterface)
		}
	}

	// VLAN interface can't be updated, so no need to check it here

	return opts
}

func (data *LinodeInterfaceModel) FlattenInterface(
	ctx context.Context, i linodego.LinodeInterface, preserveKnown bool, diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(i.ID), preserveKnown)

	flattenedDefaultRoute := helper.KeepOrUpdateSingleNestedAttributes(
		ctx, data.DefaultRoute, preserveKnown, diags, func(dr *DefaultRouteAttrModel, isNull *bool, pk bool, _ *diag.Diagnostics) {
			if i.DefaultRoute == nil {
				dr.IPv4 = helper.KeepOrUpdateValue(dr.IPv4, types.BoolNull(), pk)
				dr.IPv6 = helper.KeepOrUpdateValue(dr.IPv6, types.BoolNull(), pk)
				*isNull = true
				return
			}
			dr.FlattenInterfaceDefaultRoute(*i.DefaultRoute, pk)
		},
	)

	if diags.HasError() {
		return
	}

	data.DefaultRoute = *flattenedDefaultRoute

	flattenedVLAN := helper.KeepOrUpdateSingleNestedAttributes(
		ctx, data.VLAN, preserveKnown, diags, func(vlan *VLANAttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			if i.VLAN == nil {
				*isNull = true
				vlan.IPAMAddress = helper.KeepOrUpdateValue(vlan.IPAMAddress, cidrtypes.NewIPv4PrefixNull(), pk)
				vlan.VLANLabel = helper.KeepOrUpdateValue(vlan.VLANLabel, types.StringNull(), pk)
				return
			}
			vlan.FlattenVLANInterface(*i.VLAN, pk)
		},
	)
	if diags.HasError() {
		return
	}

	data.VLAN = *flattenedVLAN
	flattenedPublic := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.Public, publicInterfaceSchema.GetType().(types.ObjectType).AttrTypes, preserveKnown, diags,
		func(public *PublicAttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			if i.Public == nil {
				*isNull = true
				public.IPv4 = helper.KeepOrUpdateValue(public.IPv4, types.ObjectNull(publicIPv4Attribute.GetType().(types.ObjectType).AttrTypes), pk)
				public.IPv6 = helper.KeepOrUpdateValue(public.IPv6, types.ObjectNull(publicIPv6Attribute.GetType().(types.ObjectType).AttrTypes), pk)
				return
			}
			public.FlattenPublicInterface(ctx, *i.Public, pk, d)
		},
	)
	if diags.HasError() {
		return
	}

	data.Public = *flattenedPublic

	// Flatten VPC interface
	flattenedVPC := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, data.VPC, vpcInterfaceSchema.GetType().(types.ObjectType).AttrTypes, preserveKnown, diags,
		func(vpc *VPCAttrModel, isNull *bool, pk bool, d *diag.Diagnostics) {
			if i.VPC == nil {
				*isNull = true
				vpc.IPv4 = helper.KeepOrUpdateValue(
					vpc.IPv4, types.ObjectNull(vpcIPv4Attribute.GetType().(types.ObjectType).AttrTypes), pk,
				)
				vpc.SubnetID = helper.KeepOrUpdateValue(vpc.SubnetID, types.Int64Null(), pk)
				return
			}
			vpc.FlattenVPCInterface(ctx, *i.VPC, pk, d)
		},
	)
	if diags.HasError() {
		return
	}

	data.VPC = *flattenedVPC
}
