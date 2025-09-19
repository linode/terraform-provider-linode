package linodeinterface

import (
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type VLANAttrModel struct {
	IPAMAddress cidrtypes.IPv4Prefix `tfsdk:"ipam_address"`
	VLANLabel   types.String         `tfsdk:"vlan_label"`
}

func (data *VLANAttrModel) CopyFrom(other VLANAttrModel, preserveKnown bool) {
	data.IPAMAddress = helper.KeepOrUpdateValue(data.IPAMAddress, other.IPAMAddress, preserveKnown)
	data.VLANLabel = helper.KeepOrUpdateValue(data.VLANLabel, other.VLANLabel, preserveKnown)
}

func (plan *VLANAttrModel) GetCreateOptions() (vlan linodego.VLANInterface) {
	if !plan.IPAMAddress.IsUnknown() {
		vlan.IPAMAddress = plan.IPAMAddress.ValueStringPointer()
	}
	if !plan.VLANLabel.IsUnknown() {
		vlan.VLANLabel = plan.VLANLabel.ValueString()
	}
	return
}

func (data *VLANAttrModel) FlattenVLANInterface(
	vlanInterface linodego.VLANInterface, preserveKnown bool,
) {
	data.VLANLabel = helper.KeepOrUpdateString(data.VLANLabel, vlanInterface.VLANLabel, preserveKnown)
	data.IPAMAddress = helper.KeepOrUpdateValue(data.IPAMAddress, cidrtypes.NewIPv4PrefixPointerValue(vlanInterface.IPAMAddress), preserveKnown)
}
