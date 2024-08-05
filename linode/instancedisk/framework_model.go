package instancedisk

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ResourceModel struct {
	ID              types.String      `tfsdk:"id"`
	Label           types.String      `tfsdk:"label"`
	LinodeID        types.Int64       `tfsdk:"linode_id"`
	Size            types.Int64       `tfsdk:"size"`
	AuthorizedKeys  types.Set         `tfsdk:"authorized_keys"`
	AuthorizedUsers types.Set         `tfsdk:"authorized_users"`
	Filesystem      types.String      `tfsdk:"filesystem"`
	Image           types.String      `tfsdk:"image"`
	RootPass        types.String      `tfsdk:"root_pass"`
	StackScriptData types.Map         `tfsdk:"stackscript_data"`
	StackScriptID   types.Int64       `tfsdk:"stackscript_id"`
	Created         timetypes.RFC3339 `tfsdk:"created"`
	Updated         timetypes.RFC3339 `tfsdk:"updated"`
	Status          types.String      `tfsdk:"status"`
	DiskEncryption  types.String      `tfsdk:"disk_encryption"`
	Timeouts        timeouts.Value    `tfsdk:"timeouts"`
}

func (data *ResourceModel) FlattenDisk(disk *linodego.InstanceDisk, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(disk.ID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, disk.Label, preserveKnown)
	data.Size = helper.KeepOrUpdateInt64(data.Size, int64(disk.Size), preserveKnown)
	data.Filesystem = helper.KeepOrUpdateString(
		data.Filesystem, string(disk.Filesystem), preserveKnown,
	)

	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(disk.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(disk.Updated), preserveKnown,
	)

	data.Status = helper.KeepOrUpdateString(data.Status, string(disk.Status), preserveKnown)
	data.DiskEncryption = helper.KeepOrUpdateString(data.DiskEncryption, string(disk.DiskEncryption), preserveKnown)
}

func (data *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.LinodeID = helper.KeepOrUpdateValue(data.LinodeID, other.LinodeID, preserveKnown)
	data.AuthorizedKeys = helper.KeepOrUpdateValue(
		data.AuthorizedKeys, other.AuthorizedKeys, preserveKnown,
	)
	data.AuthorizedUsers = helper.KeepOrUpdateValue(
		data.AuthorizedUsers, other.AuthorizedUsers, preserveKnown,
	)
	data.Filesystem = helper.KeepOrUpdateValue(data.Filesystem, other.Filesystem, preserveKnown)
	data.Image = helper.KeepOrUpdateValue(data.Image, other.Image, preserveKnown)
	data.RootPass = helper.KeepOrUpdateValue(data.RootPass, other.RootPass, preserveKnown)
	data.StackScriptData = helper.KeepOrUpdateValue(
		data.StackScriptData, other.StackScriptData, preserveKnown,
	)
	data.StackScriptID = helper.KeepOrUpdateValue(
		data.StackScriptID, other.StackScriptID, preserveKnown,
	)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, other.Updated, preserveKnown)
	data.Status = helper.KeepOrUpdateValue(data.Status, other.Status, preserveKnown)
	data.DiskEncryption = helper.KeepOrUpdateValue(data.DiskEncryption, other.DiskEncryption, preserveKnown)
}
