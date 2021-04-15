package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"

	"context"
	"fmt"
	"time"
)

func dataSourceLinodeInstanceBackupSnapshotDisk() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of this disk.",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The size of this disk.",
				Computed:    true,
			},
			"filesystem": {
				Type:        schema.TypeString,
				Description: "The filesystem of this disk.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeInstanceBackup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The unique ID of this Backup.",
				Computed:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "A label for Backups that are of type snapshot.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The current state of a specific Backup.",
				Computed:    true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "This indicates whether the Backup is an automatic Backup or manual snapshot" +
					" taken by the User at a specific point in time.",
				Computed: true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "The date the Backup was taken.",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "The date the Backup was most recently updated.",
				Computed:    true,
			},
			"finished": {
				Type:        schema.TypeString,
				Description: "The date the Backup completed.",
				Computed:    true,
			},
			"configs": {
				Type:        schema.TypeList,
				Description: "A list of the labels of the Configuration profiles that are part of the Backup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"disks": {
				Type:        schema.TypeList,
				Description: "A list of the disks that are part of the Backup.",
				Elem:        dataSourceLinodeInstanceBackupSnapshotDisk(),
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeInstanceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinodeInstanceBackupsRead,

		Schema: map[string]*schema.Schema{
			"linode_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Linode to get backups for.",
				Required:    true,
			},
			"automatic": {
				Type:        schema.TypeList,
				Description: "A list of backups or snapshots for a Linode.",
				Computed:    true,
				Elem:        dataSourceLinodeInstanceBackup(),
			},
			"current": {
				Type:        schema.TypeList,
				Description: "The current Backup for a Linode.",
				Computed:    true,
				Elem:        dataSourceLinodeInstanceBackup(),
			},
			"in_progress": {
				Type:        schema.TypeList,
				Description: "The in-progress Backup for a Linode",
				Computed:    true,
				Elem:        dataSourceLinodeInstanceBackup(),
			},
		},
	}
}

func dataSourceLinodeInstanceBackupsRead(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

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

	if snapshot.Created != nil {
		result["created"] = snapshot.Created.Format(time.RFC3339)
	}

	if snapshot.Updated != nil {
		result["updated"] = snapshot.Updated.Format(time.RFC3339)
	}

	if snapshot.Finished != nil {
		result["finished"] = snapshot.Finished.Format(time.RFC3339)
	}

	result["configs"] = snapshot.Configs

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
