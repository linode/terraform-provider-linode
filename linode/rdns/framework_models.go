package rdns

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"
)

type ResourceModel struct {
	Address          customtypes.IPAddrStringValue `tfsdk:"address"`
	RDNS             types.String                  `tfsdk:"rdns"`
	Reserved         types.Bool                    `tfsdk:"reserved"`
	WaitForAvailable types.Bool                    `tfsdk:"wait_for_available"`
	ID               types.String                  `tfsdk:"id"`
	Timeouts         timeouts.Value                `tfsdk:"timeouts"`
}

func (rm *ResourceModel) FlattenInstanceIP(ip *linodego.InstanceIP, preserveKnown bool) {
	rm.Address = helper.KeepOrUpdateValue(
		rm.Address, customtypes.IPAddrValue(ip.Address), preserveKnown,
	)
	rm.ID = helper.KeepOrUpdateString(rm.ID, ip.Address, preserveKnown)
	rm.RDNS = helper.KeepOrUpdateString(rm.RDNS, ip.RDNS, preserveKnown)
	rm.Reserved = helper.KeepOrUpdateValue(rm.Reserved, types.BoolValue(ip.Reserved), preserveKnown)
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.Address = helper.KeepOrUpdateValue(rm.Address, other.Address, preserveKnown)
	rm.ID = helper.KeepOrUpdateValue(rm.ID, other.ID, preserveKnown)
	rm.RDNS = helper.KeepOrUpdateValue(rm.RDNS, other.RDNS, preserveKnown)
}
