package firewalldevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type FirewallDeviceModel struct {
	ID         types.Int64       `tfsdk:"id"`
	FirewallID types.Int64       `tfsdk:"firewall_id"`
	EntityID   types.Int64       `tfsdk:"entity_id"`
	EntityType types.String      `tfsdk:"entity_type"`
	Created    timetypes.RFC3339 `tfsdk:"created"`
	Updated    timetypes.RFC3339 `tfsdk:"updated"`
}

func (fdm *FirewallDeviceModel) FlattenFirewallDevice(
	device *linodego.FirewallDevice,
	preserveKnown bool,
) {
	// firewall_id is always configured by the user and
	// never appear in linodego.FirewallDevice

	fdm.ID = helper.KeepOrUpdateInt64(fdm.ID, int64(device.ID), preserveKnown)
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

	fdm.Created = helper.KeepOrUpdateValue(
		fdm.Created,
		timetypes.NewRFC3339TimePointerValue(device.Created),
		preserveKnown,
	)
	fdm.Updated = helper.KeepOrUpdateValue(
		fdm.Updated,
		timetypes.NewRFC3339TimePointerValue(device.Updated),
		preserveKnown,
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
