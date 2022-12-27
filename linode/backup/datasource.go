package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func dataSourceDisk() *schema.Resource {
	return &schema.Resource{
		Schema: diskSchema,
	}
}

// contains a array of dataSourceDisk
func dataSourceBackup() *schema.Resource {
	return &schema.Resource{
		Schema: backupSchema,
	}
}

// contains a array of dataSourceBackup
func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	backups, err := client.GetInstanceBackups(ctx, linodeID)
	if err != nil {
		return diag.Errorf("failed to get backups: %s", err)
	}

	flattenedAutoSnapshots := make([]map[string]interface{}, len(backups.Automatic))
	for i, snapshot := range backups.Automatic {
		flattenedAutoSnapshots[i] = flattenInstanceSnapshot(snapshot)
	}

	d.SetId(fmt.Sprintf("%d", linodeID))

	d.Set("automatic", flattenedAutoSnapshots)

	if backups.Snapshot.Current != nil {
		d.Set("current", []interface{}{flattenInstanceSnapshot(backups.Snapshot.Current)})
	}

	if backups.Snapshot.InProgress != nil {
		d.Set("in_progress", []interface{}{flattenInstanceSnapshot(backups.Snapshot.InProgress)})
	}

	return nil
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
