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
	var defaultFirewallIDsModel DefaultFirewallIDsAttributeModel
	var defaultFirewallIDs linodego.DefaultFirewallIDsOptions

	diags.Append(fsds.DefaultFirewallIDs.As(ctx, &defaultFirewallIDsModel, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	if diags.HasError() {
		return opts
	}

	shouldUpdateDefaultFirewallIDs := false

	if !defaultFirewallIDsModel.Linode.IsUnknown() {
		defaultFirewallIDs.Linode = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDsModel.Linode.ValueInt64Pointer(), diags),
		)
		shouldUpdateDefaultFirewallIDs = true
	}

	if !defaultFirewallIDsModel.NodeBalancer.IsUnknown() {
		defaultFirewallIDs.NodeBalancer = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDsModel.NodeBalancer.ValueInt64Pointer(), diags),
		)
		shouldUpdateDefaultFirewallIDs = true
	}

	if !defaultFirewallIDsModel.PublicInterface.IsUnknown() {
		defaultFirewallIDs.PublicInterface = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDsModel.PublicInterface.ValueInt64Pointer(), diags),
		)
		shouldUpdateDefaultFirewallIDs = true
	}

	if !defaultFirewallIDsModel.VPCInterface.IsUnknown() {
		defaultFirewallIDs.VPCInterface = linodego.Pointer(
			helper.FrameworkSafeInt64PointerToIntPointer(defaultFirewallIDsModel.VPCInterface.ValueInt64Pointer(), diags),
		)
		shouldUpdateDefaultFirewallIDs = true
	}

	if diags.HasError() {
		return opts
	}

	if shouldUpdateDefaultFirewallIDs {
		opts.DefaultFirewallIDs = &defaultFirewallIDs
	}

	return opts
}

func (fsds *FirewallSettingsModel) FlattenFirewallSettings(
	ctx context.Context,
	settings linodego.FirewallSettings,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	defaultFirewallIDs := helper.KeepOrUpdateSingleNestedAttribute(
		ctx,
		fsds.DefaultFirewallIDs,
		preserveKnown,
		diags,
		func(defaultFirewallIDsAttrsModel *DefaultFirewallIDsAttributeModel, _ *bool, preserveKnown bool, _ *diag.Diagnostics) {
			defaultFirewallIDsAttrsModel.FlattenDefaultFirewallIDs(settings, preserveKnown)
		},
	)

	if diags.HasError() {
		return
	}

	fsds.DefaultFirewallIDs = *defaultFirewallIDs
}

func (dfiam *DefaultFirewallIDsAttributeModel) FlattenDefaultFirewallIDs(settings linodego.FirewallSettings, preserveKnown bool) {
	dfiam.Linode = helper.KeepOrUpdateInt64Pointer(dfiam.Linode, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.Linode), preserveKnown)
	dfiam.NodeBalancer = helper.KeepOrUpdateInt64Pointer(dfiam.NodeBalancer, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.NodeBalancer), preserveKnown)
	dfiam.PublicInterface = helper.KeepOrUpdateInt64Pointer(
		dfiam.PublicInterface,
		helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.PublicInterface),
		preserveKnown,
	)
	dfiam.VPCInterface = helper.KeepOrUpdateInt64Pointer(dfiam.VPCInterface, helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.VPCInterface), preserveKnown)
}
