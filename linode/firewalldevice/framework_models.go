package firewalldevice

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type FirewallDeviceModel struct {
	ID         types.String `tfsdk:"id"`
	FirewallID types.Int64  `tfsdk:"firewall_id"`
	EntityID   types.Int64  `tfsdk:"entity_id"`
	EntityType types.String `tfsdk:"entity_type"`
	Created    types.String `tfsdk:"created"`
	Updated    types.String `tfsdk:"updated"`
}

func (fdm *FirewallDeviceModel) FlattenFirewallDevice(
	device *linodego.FirewallDevice,
	preserveKnown bool,
) {
	// firewall_id is always configured by the user and
	// never appear in linodego.FirewallDevice

	fdm.ID = helper.KeepOrUpdateString(fdm.ID, strconv.Itoa(device.ID), preserveKnown)
	fdm.EntityID = helper.KeepOrUpdateInt64(
		fdm.EntityID,
		int64(device.Entity.ID),
		preserveKnown,
	)
	fdm.EntityType = helper.KeepOrUpdateString(
		fdm.EntityType,
		string(device.Entity.Type),
		preserveKnown,
	)
	fdm.Created = helper.KeepOrUpdateString(
		fdm.Created, device.Created.Format(time.RFC3339), preserveKnown,
	)
	fdm.Updated = helper.KeepOrUpdateString(
		fdm.Updated, device.Updated.Format(time.RFC3339), preserveKnown,
	)
}

func (fdm *FirewallDeviceModel) CopyFrom(
	other FirewallDeviceModel,
	preserveKnown bool,
) {
	fdm.ID = helper.KeepOrUpdateValue(fdm.ID, other.ID, preserveKnown)
	fdm.FirewallID = helper.KeepOrUpdateValue(fdm.FirewallID, other.FirewallID, preserveKnown)
	fdm.EntityID = helper.KeepOrUpdateValue(fdm.EntityID, other.EntityID, preserveKnown)
	fdm.EntityType = helper.KeepOrUpdateValue(fdm.EntityType, other.EntityType, preserveKnown)
	fdm.Created = helper.KeepOrUpdateValue(fdm.Created, other.Created, preserveKnown)
	fdm.Updated = helper.KeepOrUpdateValue(fdm.Updated, other.Updated, preserveKnown)
}
