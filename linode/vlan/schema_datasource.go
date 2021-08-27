package vlan

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The unique label of this VLAN.",
		Computed:    true,
	},
	"linodes": {
		Type:        schema.TypeList,
		Description: "The Linodes currently attached to this VLAN.",
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region this VLAN is located in.",
		Computed:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this VLAN was created.",
		Computed:    true,
	},
}
