package volume

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique id of this Volume.",
		Required:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The Volume's label. For display purposes only.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The datacenter where this Volume is located.",
		Computed:    true,
	},
	"size": {
		Type:        schema.TypeInt,
		Description: "The size of this Volume in GiB.",
		Computed:    true,
	},
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "If a Volume is attached to a specific Linode, the ID of that Linode will be displayed here.",
		Computed:    true,
	},
	"filesystem_path": {
		Type: schema.TypeString,
		Description: "The full filesystem path for the Volume based on the Volume's label. Path is " +
			"/dev/disk/by-id/scsi-0LinodeVolume + Volume label.",
		Computed: true,
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "An array of tags applied to this Volume. Tags are for organizational purposes only.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The status of the Volume. Can be one of active | creating | resizing | contact_support",
		Computed:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "Datetime string representing when the Volume was created.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "Datetime string representing when the Volume was last updated.",
		Computed:    true,
	},
}
