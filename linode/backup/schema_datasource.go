package backup

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var diskSchema = map[string]*schema.Schema{
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
}

var backupSchema = map[string]*schema.Schema{
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
		Elem:        dataSourceDisk(),
		Computed:    true,
	},
}

var dataSourceSchema = map[string]*schema.Schema{
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Linode to get backups for.",
		Required:    true,
	},
	"automatic": {
		Type:        schema.TypeList,
		Description: "A list of backups or snapshots for a Linode.",
		Computed:    true,
		Elem:        dataSourceBackup(),
	},
	"current": {
		Type:        schema.TypeList,
		Description: "The current Backup for a Linode.",
		Computed:    true,
		Elem:        dataSourceBackup(),
	},
	"in_progress": {
		Type:        schema.TypeList,
		Description: "The in-progress Backup for a Linode",
		Computed:    true,
		Elem:        dataSourceBackup(),
	},
}
