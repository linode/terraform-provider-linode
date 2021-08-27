package sshkey

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode SSH Key.",
		Required:    true,
	},
	"ssh_key": {
		Type:        schema.TypeString,
		Description: "The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.",
		Required:    true,
		ForceNew:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "The date this key was added.",
		Computed:    true,
	},
}
