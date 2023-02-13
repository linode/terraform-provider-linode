package region

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"country": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The country where this Region resides.",
		Computed:    true,
	},
	"id": {
		Type:        schema.TypeString,
		Description: "The unique ID of this Region.",
		Required:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "Detailed location information for this Region, including city, state or region, and country.",
		Computed:    true,
	},
}
