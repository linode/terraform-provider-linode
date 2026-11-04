package nfsfilesystem

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var statsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"used_capacity_bytes": types.Int64Type,
		"used_file_count":     types.Int64Type,
		"collected_at":        timetypes.RFC3339Type{},
	},
}

// StatsModel represents the live usage statistics of an NFS Filesystem.
type StatsModel struct {
	UsedCapacityBytes types.Int64       `tfsdk:"used_capacity_bytes"`
	UsedFileCount     types.Int64       `tfsdk:"used_file_count"`
	CollectedAt       timetypes.RFC3339 `tfsdk:"collected_at"`
}

// BaseModel contains the fields shared between the resource and data source models.
type BaseModel struct {
	ID                 types.String      `tfsdk:"id"`
	SpaceID            types.Int64       `tfsdk:"space_id"`
	Label              types.String      `tfsdk:"label"`
	Region             types.String      `tfsdk:"region"`
	ProtocolVersions   types.List        `tfsdk:"protocol_versions"`
	Status             types.String      `tfsdk:"status"`
	MountTargetIPs     types.List        `tfsdk:"mount_target_ips"`
	MountTargetFQDN    types.String      `tfsdk:"mount_target_fqdn"`
	MaxCapacityBytes   types.Int64       `tfsdk:"max_capacity_bytes"`
	SnapshotUsageBytes types.Int64       `tfsdk:"snapshot_usage_bytes"`
	LDAPConfigID       types.String      `tfsdk:"ldap_config_id"`
	SourceSnapshotID   types.Int64       `tfsdk:"source_snapshot_id"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Updated            timetypes.RFC3339 `tfsdk:"updated"`
	Tags               types.Set         `tfsdk:"tags"`
	Stats              types.Object      `tfsdk:"stats"`
}

func (m *BaseModel) Flatten(
	ctx context.Context,
	fs *linodego.NFSFilesystem,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(fs.ID), preserveKnown)
	m.SpaceID = helper.KeepOrUpdateInt64(m.SpaceID, int64(fs.SpaceID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, fs.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, fs.Region, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(fs.Status), preserveKnown)
	m.MountTargetFQDN = helper.KeepOrUpdateStringPointer(m.MountTargetFQDN, fs.MountTargetFQDN, preserveKnown)
	//// The API does not reliably return max_capacity_bytes on GET. Only update when the API reports a non-zero value,
	//// otherwise preserve the existing configured/known value to avoid a
	//// perpetual diff.
	//if fs.MaxCapacityBytes != 0 {
	m.MaxCapacityBytes = helper.KeepOrUpdateInt64(m.MaxCapacityBytes, fs.MaxCapacityBytes, preserveKnown)
	//}

	m.MaxCapacityBytes = helper.KeepOrUpdateInt64(m.MaxCapacityBytes, fs.MaxCapacityBytes, preserveKnown)

	m.SnapshotUsageBytes = helper.KeepOrUpdateInt64Pointer(m.SnapshotUsageBytes, fs.SnapshotUsageBytes, preserveKnown)
	m.LDAPConfigID = helper.KeepOrUpdateStringPointer(m.LDAPConfigID, fs.LDAPConfigID, preserveKnown)
	m.SourceSnapshotID = helper.KeepOrUpdateIntPointer(m.SourceSnapshotID, fs.SourceSnapshotID, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(
		m.Created,
		timetypes.NewRFC3339TimePointerValue(fs.Created),
		preserveKnown,
	)
	m.Updated = helper.KeepOrUpdateValue(
		m.Updated,
		timetypes.NewRFC3339TimePointerValue(fs.Updated),
		preserveKnown,
	)

	protocols := helper.MapSlice(fs.ProtocolVersions, func(p linodego.NFSProtocolVersion) string {
		return string(p)
	})
	protocolList, d := types.ListValueFrom(ctx, types.StringType, protocols)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.ProtocolVersions = helper.KeepOrUpdateValue(m.ProtocolVersions, protocolList, preserveKnown)

	mountTargetIPs, d := types.ListValueFrom(ctx, types.StringType, fs.MountTargetIPs)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.MountTargetIPs = helper.KeepOrUpdateValue(m.MountTargetIPs, mountTargetIPs, preserveKnown)

	tags, d := types.SetValueFrom(ctx, types.StringType, fs.Tags)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Tags = helper.KeepOrUpdateValue(m.Tags, tags, preserveKnown)

	stats := StatsModel{
		UsedCapacityBytes: types.Int64PointerValue(fs.Stats.UsedCapacityBytes),
		UsedFileCount:     types.Int64PointerValue(fs.Stats.UsedFileCount),
		CollectedAt:       timetypes.NewRFC3339TimePointerValue(fs.Stats.CollectedAt),
	}
	statsObject, d := types.ObjectValueFrom(ctx, statsObjectType.AttrTypes, stats)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Stats = helper.KeepOrUpdateValue(m.Stats, statsObject, preserveKnown)

	return diags
}

func (m *BaseModel) CopyFrom(other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.SpaceID = helper.KeepOrUpdateValue(m.SpaceID, other.SpaceID, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.ProtocolVersions = helper.KeepOrUpdateValue(m.ProtocolVersions, other.ProtocolVersions, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.MountTargetIPs = helper.KeepOrUpdateValue(m.MountTargetIPs, other.MountTargetIPs, preserveKnown)
	m.MountTargetFQDN = helper.KeepOrUpdateValue(m.MountTargetFQDN, other.MountTargetFQDN, preserveKnown)
	m.MaxCapacityBytes = helper.KeepOrUpdateValue(m.MaxCapacityBytes, other.MaxCapacityBytes, preserveKnown)
	m.SnapshotUsageBytes = helper.KeepOrUpdateValue(m.SnapshotUsageBytes, other.SnapshotUsageBytes, preserveKnown)
	m.LDAPConfigID = helper.KeepOrUpdateValue(m.LDAPConfigID, other.LDAPConfigID, preserveKnown)
	m.SourceSnapshotID = helper.KeepOrUpdateValue(m.SourceSnapshotID, other.SourceSnapshotID, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Tags = helper.KeepOrUpdateValue(m.Tags, other.Tags, preserveKnown)
	m.Stats = helper.KeepOrUpdateValue(m.Stats, other.Stats, preserveKnown)
}

// ResourceModel represents the Terraform state for the linode_nfs_filesystem resource.
type ResourceModel struct {
	BaseModel
}

// DataSourceModel represents the Terraform state for the linode_nfs_filesystem data source.
type DataSourceModel struct {
	BaseModel
}
