package firewalldevice

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceSchema = map[string]*schema.Schema{
	"firewall_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Firewall to access.",
		Required:    true,
		ForceNew:    true,
	},
	"entity_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the entity to create a Firewall device for.",
		Required:    true,
		ForceNew:    true,
	},
	"entity_type": {
		Type:        schema.TypeString,
		Description: "The type of the entity to create a Firewall device for.",
		Default:     "linode",
		Optional:    true,
		ForceNew:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this Firewall Device was created.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this Firewall Device was updated.",
		Computed:    true,
	},
}
