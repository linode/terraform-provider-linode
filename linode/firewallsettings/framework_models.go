package firewallsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (fsds *FirewallSettingsModel) ParseFirewallSettings(ctx context.Context, settings linodego.FirewallSettings, diags *diag.Diagnostics) {
	defaultIDs, newDiags := types.ObjectValueFrom(ctx, fsds.DefaultFirewallIDs.AttributeTypes(ctx), DefaultFirewallIDsAttributeModel{
		Linode:          types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.Linode)),
		NodeBalancer:    types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.NodeBalancer)),
		PublicInterface: types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.PublicInterface)),
		VPCInterface:    types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.VPCInterface)),
	})
	diags.Append(newDiags...)

	fsds.DefaultFirewallIDs = defaultIDs
}
