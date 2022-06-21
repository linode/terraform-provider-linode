package volume

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"source_volume_id": {
		Type:        schema.TypeInt,
		Description: "The ID of a volume to clone.",
		Optional:    true,
		ForceNew:    true,
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return newValue == ""
		},
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode Volume.",
		Required:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The status of the volume, indicating the current readiness state.",
		Computed:    true,
	},
	"region": {
		Type:         schema.TypeString,
		Description:  "The region where this volume will be deployed.",
		Required:     true,
		ForceNew:     true,
		InputDefault: "us-east",
	},
	"size": {
		Type:        schema.TypeInt,
		Description: "Size of the Volume in GB",
		Optional:    true,
		Computed:    true,
	},
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The Linode ID where the Volume should be attached.",
		Optional:    true,
		Computed:    true,
	},
	"filesystem_path": {
		Type: schema.TypeString,
		Description: "The full filesystem path for the Volume based on the Volume's label. Path is " +
			"/dev/disk/by-id/scsi-0Linode_Volume_ + Volume label.",
		Computed: true,
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
}
