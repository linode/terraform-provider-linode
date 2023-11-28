package instancedisk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Disk’s label is for display purposes only.",
	},
	"linode_id": {
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the Linode to assign this disk to.",
	},
	"size": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The size of the Disk in MB.",
	},

	"authorized_keys": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		ForceNew: true,
		Description: "A list of public SSH keys that will be automatically appended to the root " +
			"user’s ~/.ssh/authorized_keys file when deploying from an Image.",
	},
	"authorized_users": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		ForceNew: true,
		Description: "A list of usernames. If the usernames have associated SSH keys, the keys will be appended to the " +
			"root users ~/.ssh/authorized_keys file automatically when deploying from an Image.",
	},
	"filesystem": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: "The filesystem of this disk.",
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(
			[]string{"raw", "swap", "ext3", "ext4", "initrd"}, true)),
	},
	"image": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "An Image ID to deploy the Linode Disk from.",
	},
	"root_pass": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Sensitive:   true,
		Description: "This sets the root user’s password on a newly-created Linode Disk when deploying from an Image.",
		ValidateFunc: validation.StringLenBetween(
			helper.RootPassMinimumCharacters,
			helper.RootPassMaximumCharacters),
	},
	"stackscript_data": {
		Type: schema.TypeMap,
		Description: "An object containing responses to any User Defined Fields present in the StackScript " +
			"being deployed to this Disk. Only accepted if 'stackscript_id' is given. The required values depend " +
			"on the StackScript being deployed.",
		Optional:  true,
		ForceNew:  true,
		Sensitive: true,
	},
	"stackscript_id": {
		Type: schema.TypeInt,
		Description: "A StackScript ID that will cause the referenced StackScript " +
			"to be run during deployment of this Linode.",
		Optional: true,
		ForceNew: true,
	},

	"created": {
		Type:        schema.TypeString,
		Description: "When this disk was created.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "A brief description of this Disk's current state.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this disk was last updated.",
		Computed:    true,
	},
}
