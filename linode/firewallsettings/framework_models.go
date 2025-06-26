package firewallsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DefaultFirewallIDsModel struct {
	Linode          types.Int64 `tfsdk:"linode"`
	NodeBalancer    types.Int64 `tfsdk:"nodebalancer"`
	PublicInterface types.Int64 `tfsdk:"public_interface"`
	VPCInterface    types.Int64 `tfsdk:"vpc_interface"`
}

type FirewallSettingsDataSourceModel struct {
	DefaultFirewallIDs DefaultFirewallIDsModel `tfsdk:"default_firewall_ids"`
}

func (fsds *FirewallSettingsDataSourceModel) ParseFirewallSettings(settings linodego.FirewallSettings) {
	fsds.DefaultFirewallIDs.Linode = types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.Linode))
	fsds.DefaultFirewallIDs.NodeBalancer = types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.NodeBalancer))
	fsds.DefaultFirewallIDs.PublicInterface = types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.PublicInterface))
	fsds.DefaultFirewallIDs.VPCInterface = types.Int64PointerValue(helper.IntPtrToInt64Ptr(settings.DefaultFirewallIDs.VPCInterface))
}
