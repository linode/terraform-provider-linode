package image

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "A short description of the Image. Labels cannot contain special characters.",
		Required:    true,
	},
	"disk_id": {
		Type:          schema.TypeInt,
		Description:   "The ID of the Linode Disk that this Image will be created from.",
		RequiredWith:  []string{"linode_id"},
		ConflictsWith: []string{"file_path"},
		Optional:      true,
		ForceNew:      true,
	},
	"linode_id": {
		Type:          schema.TypeInt,
		Description:   "The ID of the Linode that this Image will be created from.",
		RequiredWith:  []string{"disk_id"},
		ConflictsWith: []string{"file_path"},
		Optional:      true,
		ForceNew:      true,
	},
	"file_path": {
		Type:          schema.TypeString,
		Description:   "The name of the file to upload to this image.",
		ConflictsWith: []string{"linode_id", "disk_id"},
		RequiredWith:  []string{"region"},
		Optional:      true,
		ForceNew:      true,
	},
	"region": {
		Type:         schema.TypeString,
		Description:  "The region to upload to.",
		RequiredWith: []string{"file_path"},
		Optional:     true,
	},
	"file_hash": {
		Type:        schema.TypeString,
		Description: "The MD5 hash of the image file.",
		Computed:    true,
		Optional:    true,
		ForceNew:    true,
	},
	"description": {
		Type:        schema.TypeString,
		Description: "A detailed description of this Image.",
		Optional:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this Image was created.",
		Computed:    true,
	},
	"created_by": {
		Type:        schema.TypeString,
		Description: "The name of the User who created this Image.",
		Computed:    true,
	},
	"deprecated": {
		Type:        schema.TypeBool,
		Description: "Whether or not this Image is deprecated. Will only be True for deprecated public Images.",
		Computed:    true,
	},
	"is_public": {
		Type:        schema.TypeBool,
		Description: "True if the Image is public.",
		Computed:    true,
	},
	"size": {
		Type:        schema.TypeInt,
		Description: "The minimum size this Image needs to deploy. Size is in MB.",
		Computed:    true,
	},
	"type": {
		Type: schema.TypeString,
		Description: "How the Image was created. 'Manual' Images can be created at any time. 'Automatic' " +
			"images are created automatically from a deleted Linode.",
		Computed: true,
	},
	"expiry": {
		Type:        schema.TypeString,
		Description: "Only Images created automatically (from a deleted Linode; type=automatic) will expire.",
		Computed:    true,
	},
	"vendor": {
		Type:        schema.TypeString,
		Description: "The upstream distribution vendor. Nil for private Images.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The current status of this Image.",
		Computed:    true,
	},
}
