package instancesharedips

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Linode to share these IP addresses with.",
		Required:    true,
		ForceNew:    true,
	},
	"addresses": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A set of IP addresses to share to the Linode",
		Required:    true,
	},
}
