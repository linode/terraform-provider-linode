package backup

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	Automatic  types.List  `tfsdk:"automatic"`
	Current    types.List  `tfsdk:"current"`
	Linode_ID  types.Int64 `tfsdk:"linode_id"`
	InProgress types.List  `tfsdk:"in_progress"`
	ID         types.Int64 `tfsdk:"id"`
}

func (data *DataSourceModel) parseBackups(
	ctx context.Context,
	backups *linodego.InstanceBackupsResponse,
	linodeId types.Int64,
) diag.Diagnostics {
	automatic, diag := flattenAutoSnapshots(ctx, backups.Automatic)
	if diag.HasError() {
		return diag
	}

	data.ID = linodeId
	data.Automatic = *automatic

	if backups.Snapshot.Current != nil {
		currentObj, diag := flattenInstanceSnapshot(ctx, backups.Snapshot.Current)
		if diag.HasError() {
			return diag
		}

		current, diag := basetypes.NewListValue(
			backupObjectType,
			[]attr.Value{currentObj},
		)
		if diag.HasError() {
			return diag
		}

		data.Current = current
	} else {
		data.Current = types.ListNull(backupObjectType)
	}

	if backups.Snapshot.InProgress != nil {
		inProgressObj, diag := flattenInstanceSnapshot(ctx, backups.Snapshot.InProgress)
		if diag.HasError() {
			return diag
		}

		inProgress, diag := basetypes.NewListValue(
			backupObjectType,
			[]attr.Value{inProgressObj},
		)
		if diag.HasError() {
			return diag
		}

		data.InProgress = inProgress
	} else {
		data.InProgress = types.ListNull(backupObjectType)
	}

	return nil
}

func flattenAutoSnapshots(
	ctx context.Context,
	snapshots []*linodego.InstanceSnapshot,
) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(snapshots))
	for i, snapshot := range snapshots {
		result, diag := flattenInstanceSnapshot(ctx, snapshot)
		if diag.HasError() {
			return nil, diag
		}

		resultList[i] = result
	}
	result, diag := basetypes.NewListValue(
		backupObjectType,
		resultList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}

func flattenInstanceSnapshot(ctx context.Context,
	snapshot *linodego.InstanceSnapshot,
) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["id"] = types.Int64Value(int64(snapshot.ID))
	result["label"] = types.StringValue(snapshot.Label)
	result["status"] = types.StringValue(string(snapshot.Status))
	result["type"] = types.StringValue(snapshot.Type)
	result["configs"], _ = types.ListValueFrom(ctx, types.StringType, snapshot.Configs)
	result["available"] = types.BoolValue(snapshot.Available)

	if snapshot.Created != nil {
		result["created"] = types.StringValue(snapshot.Created.Format(time.RFC3339))
	} else {
		result["created"] = types.StringNull()
	}

	if snapshot.Updated != nil {
		result["updated"] = types.StringValue(snapshot.Updated.Format(time.RFC3339))
	} else {
		result["updated"] = types.StringNull()
	}

	if snapshot.Finished != nil {
		result["finished"] = types.StringValue(snapshot.Finished.Format(time.RFC3339))
	} else {
		result["finished"] = types.StringNull()
	}

	if snapshot.Disks != nil {
		disks, diag := flattenSnapshotDisk(snapshot.Disks)
		if diag.HasError() {
			return nil, diag
		}

		result["disks"] = disks
	}

	obj, diag := types.ObjectValue(backupObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	return &obj, nil
}

func flattenSnapshotDisk(
	disks []*linodego.InstanceSnapshotDisk,
) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(disks))

	for i, field := range disks {
		valueMap := make(map[string]attr.Value)
		valueMap["label"] = types.StringValue(field.Label)
		valueMap["size"] = types.Int64Value(int64(field.Size))
		valueMap["filesystem"] = types.StringValue(field.Filesystem)

		obj, diag := types.ObjectValue(diskObjectType.AttrTypes, valueMap)
		if diag.HasError() {
			return nil, diag
		}

		resultList[i] = obj
	}

	result, diag := basetypes.NewListValue(
		diskObjectType,
		resultList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}
