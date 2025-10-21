package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type VLANAttrModel struct {
	IPAMAddress cidrtypes.IPv4Prefix `tfsdk:"ipam_address"`
	VLANLabel   types.String         `tfsdk:"vlan_label"`
}

func (data *VLANAttrModel) CopyFrom(ctx context.Context, other VLANAttrModel, preserveKnown bool) {
	tflog.Trace(ctx, "Enter VLANAttrModel.CopyFrom")

	data.IPAMAddress = helper.KeepOrUpdateValue(data.IPAMAddress, other.IPAMAddress, preserveKnown)
	data.VLANLabel = helper.KeepOrUpdateValue(data.VLANLabel, other.VLANLabel, preserveKnown)
}

func (plan *VLANAttrModel) GetCreateOptions(ctx context.Context) (vlan linodego.VLANInterface) {
	tflog.Trace(ctx, "Enter VLANAttrModel.GetCreateOptions")
	if !plan.IPAMAddress.IsUnknown() {
		vlan.IPAMAddress = plan.IPAMAddress.ValueStringPointer()
	}
	if !plan.VLANLabel.IsUnknown() {
		vlan.VLANLabel = plan.VLANLabel.ValueString()
	}
	return vlan
}

func (data *VLANAttrModel) FlattenVLANInterface(
	ctx context.Context, vlanInterface linodego.VLANInterface, preserveKnown bool,
) {
	tflog.Trace(ctx, "Enter VLANAttrModel.FlattenVLANInterface")

	data.VLANLabel = helper.KeepOrUpdateString(data.VLANLabel, vlanInterface.VLANLabel, preserveKnown)
	data.IPAMAddress = helper.KeepOrUpdateValue(data.IPAMAddress, cidrtypes.NewIPv4PrefixPointerValue(vlanInterface.IPAMAddress), preserveKnown)
}
