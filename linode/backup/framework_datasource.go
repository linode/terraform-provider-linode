package backup

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type DataSource struct {
	client *linodego.Client
}

type DataSourceModel struct {
	Automatic  types.List  `tfsdk:"automatic"`
	Current    types.List  `tfsdk:"current"`
	ID         types.Int64 `tfsdk:"id"`
	InProgress types.List  `tfsdk:"in_progress"`
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_instance_backups"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data DataSourceModel
	client := d.client

	linodeId := data.ID.ValueInt64()
	print("here", linodeId)
	backups, err := client.GetInstanceBackups(ctx, int(linodeId))

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get backups",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseBackups(ctx, backups, linodeId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *DataSourceModel) parseBackups(ctx context.Context, backups *linodego.InstanceBackupsResponse, linodeId int64) diag.Diagnostics {
	flattenedAutoSnapshots := make([]map[string]attr.Value, len(backups.Automatic))
	for i, snapshot := range backups.Automatic {
		flattenedAutoSnapshots[i] = flattenInstanceSnapshot(ctx, snapshot)
	}

	automatic, err := types.ListValueFrom(ctx, backupObjectType, flattenedAutoSnapshots)
	if err != nil {
		return err
	}

	current, err := types.ListValueFrom(ctx, backupObjectType, []interface{}{flattenInstanceSnapshot(ctx, backups.Snapshot.Current)})
	if err != nil {
		return err
	}

	in_progress, err := types.ListValueFrom(ctx, backupObjectType, []interface{}{flattenInstanceSnapshot(ctx, backups.Snapshot.InProgress)})
	if err != nil {
		return err
	}

	data.ID = types.Int64Value(linodeId)
	data.Automatic = automatic
	data.Current = current
	data.InProgress = in_progress

	return nil
}

func flattenInstanceSnapshot(ctx context.Context, snapshot *linodego.InstanceSnapshot) map[string]attr.Value {
	result := make(map[string]attr.Value)

	result["id"] = types.Int64Value(int64(snapshot.ID))
	result["label"] = types.StringValue(snapshot.Label)
	result["status"] = types.StringValue(string(snapshot.Status))
	result["type"] = types.StringValue(snapshot.Type)
	result["configs"], _ = types.ListValueFrom(ctx, types.StringType, snapshot.Configs)
	result["available"] = types.BoolValue(snapshot.Available)

	if snapshot.Created != nil {
		result["created"] = types.StringValue(snapshot.Created.Format(time.RFC3339))
	}

	if snapshot.Updated != nil {
		result["updated"] = types.StringValue(snapshot.Updated.Format(time.RFC3339))
	}

	if snapshot.Finished != nil {
		result["finished"] = types.StringValue(snapshot.Finished.Format(time.RFC3339))
	}

	flattenedDisks := make([]map[string]attr.Value, len(snapshot.Disks))
	for i, disk := range snapshot.Disks {
		flattenedDisks[i] = flattenSnapshotDisk(disk)
	}

	result["disks"], _ = types.ListValueFrom(ctx, diskObjectType, flattenedDisks)

	return result
}

func flattenSnapshotDisk(disk *linodego.InstanceSnapshotDisk) map[string]attr.Value {
	result := make(map[string]attr.Value)

	result["label"] = types.StringValue(disk.Label)
	result["size"] = types.Int64Value(int64(disk.Size))
	result["filesystem"] = types.StringValue(disk.Filesystem)

	return result
}
