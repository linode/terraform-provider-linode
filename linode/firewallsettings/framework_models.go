package firewallsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DefaultFirewallIDsAttributeModel struct {
	Linode          types.Int64 `tfsdk:"linode"`
	NodeBalancer    types.Int64 `tfsdk:"nodebalancer"`
	PublicInterface types.Int64 `tfsdk:"public_interface"`
	VPCInterface    types.Int64 `tfsdk:"vpc_interface"`
}

type FirewallSettingsModel struct {
	DefaultFirewallIDs types.Object `tfsdk:"default_firewall_ids"`
}

func (fsds *FirewallSettingsModel) GetUpdateOptions(
	ctx context.Context,
	diags *diag.Diagnostics,
) (opts linodego.FirewallSettingsUpdateOptions) {
	var defaultFirewallIDs DefaultFirewallIDsAttributeModel
	diags.Append(fsds.DefaultFirewallIDs.As(ctx, &defaultFirewallIDs, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	if diags.HasError() {
		return
	}

	if !defaultFirewallIDs.Linode.IsUnknown() {
		opts.DefaultFirewallIDs.Linode = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDs.Linode.ValueInt64Pointer(), diags),
		)
	}

	if !defaultFirewallIDs.NodeBalancer.IsUnknown() {
		opts.DefaultFirewallIDs.NodeBalancer = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDs.NodeBalancer.ValueInt64Pointer(), diags),
		)
	}

	if !defaultFirewallIDs.PublicInterface.IsUnknown() {
		opts.DefaultFirewallIDs.PublicInterface = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDs.PublicInterface.ValueInt64Pointer(), diags),
		)
	}

	if !defaultFirewallIDs.VPCInterface.IsUnknown() {
		opts.DefaultFirewallIDs.VPCInterface = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDs.VPCInterface.ValueInt64Pointer(), diags),
		)
	}

	return
}

func (fsds *FirewallSettingsModel) FlattenFirewallSettings(
	ctx context.Context,
	settings linodego.FirewallSettings,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	if preserveKnown && fsds.DefaultFirewallIDs.IsNull() {
		return
	}

	var defaultFirewallIDs DefaultFirewallIDsAttributeModel

	if !fsds.DefaultFirewallIDs.IsUnknown() && !fsds.DefaultFirewallIDs.IsNull() {
		diags.Append(
			fsds.DefaultFirewallIDs.As(ctx, &defaultFirewallIDs, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...,
		)

		if diags.HasError() {
			return
		}
	}

	// When the DefaultFirewallIDs wrapper object is unknown,
	// we need to override all nested known values (not to preserve them).
	preserveKnown = preserveKnown && !fsds.DefaultFirewallIDs.IsUnknown()

	defaultFirewallIDs.FlattenFirewallSettings(settings, preserveKnown)

	defaultIDs, newDiags := types.ObjectValueFrom(
		ctx,
		fsds.DefaultFirewallIDs.AttributeTypes(ctx),
		defaultFirewallIDs,
	)
	diags.Append(newDiags...)
	if diags.HasError() {
		return
	}

	fsds.DefaultFirewallIDs = defaultIDs
}

func (dfiam *DefaultFirewallIDsAttributeModel) FlattenFirewallSettings(settings linodego.FirewallSettings, preserveKnown bool) {
	dfiam.Linode = helper.KeepOrUpdateInt64Pointer(dfiam.Linode, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.Linode), preserveKnown)
	dfiam.NodeBalancer = helper.KeepOrUpdateInt64Pointer(dfiam.NodeBalancer, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.NodeBalancer), preserveKnown)
	dfiam.PublicInterface = helper.KeepOrUpdateInt64Pointer(
		dfiam.PublicInterface,
		helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.PublicInterface),
		preserveKnown,
	)
	dfiam.VPCInterface = helper.KeepOrUpdateInt64Pointer(dfiam.VPCInterface, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.VPCInterface), preserveKnown)
}
