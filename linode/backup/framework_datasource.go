package backup

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

	meta := helper.GetMetaFromProviderDataDatasource(req, resp)
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
	backups, err := client.GetInstanceBackups(ctx, int(linodeId))

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get backups",
			err.Error(),
		)
		return
	}

	data.parseBackups(ctx, backups, linodeId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *DataSourceModel) parseBackups(ctx context.Context, backups *linodego.InstanceBackupsResponse, linodeId int64) {
	flattenedAutoSnapshots := make([]map[string]interface{}, len(backups.Automatic))
	for i, snapshot := range backups.Automatic {
		flattenedAutoSnapshots[i] = flattenInstanceSnapshot(snapshot)
	}

	data.ID = types.Int64Value(linodeId)
	data.Automatic, _ = types.ListValueFrom(ctx, backupObjectType, flattenedAutoSnapshots)
	data.Current, _ = types.ListValueFrom(ctx, backupObjectType, flattenInstanceSnapshot(backups.Snapshot.Current))
	data.InProgress, _ = types.ListValueFrom(ctx, backupObjectType, flattenInstanceSnapshot(backups.Snapshot.InProgress))
}

func flattenInstanceSnapshot(snapshot *linodego.InstanceSnapshot) map[string]interface{} {
	result := make(map[string]interface{})

	result["id"] = snapshot.ID
	result["label"] = snapshot.Label
	result["status"] = snapshot.Status
	result["type"] = snapshot.Type
	result["configs"] = snapshot.Configs
	result["available"] = snapshot.Available

	if snapshot.Created != nil {
		result["created"] = snapshot.Created.Format(time.RFC3339)
	}

	if snapshot.Updated != nil {
		result["updated"] = snapshot.Updated.Format(time.RFC3339)
	}

	if snapshot.Finished != nil {
		result["finished"] = snapshot.Finished.Format(time.RFC3339)
	}

	flattenedDisks := make([]map[string]interface{}, len(snapshot.Disks))
	for i, disk := range snapshot.Disks {
		flattenedDisks[i] = flattenSnapshotDisk(disk)
	}

	result["disks"] = flattenedDisks

	return result
}

func flattenSnapshotDisk(disk *linodego.InstanceSnapshotDisk) map[string]interface{} {
	result := make(map[string]interface{})

	result["label"] = disk.Label
	result["size"] = disk.Size
	result["filesystem"] = disk.Filesystem

	return result
}
