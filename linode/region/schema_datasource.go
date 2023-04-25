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
	"capabilities": {
		Type:        schema.TypeList,
		Description: "A list of capabilities of this region.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"status": {
		Type:        schema.TypeString,
		Description: "This region’s current operational status.",
		Computed:    true,
	},
	"resolvers": {
		Computed: true,
		Type:     schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The IPv4 addresses for this region’s DNS resolvers, separated by commas.",
				},
				"ipv6": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The IPv6 addresses for this region’s DNS resolvers, separated by commas.",
				},
			},
		},
	},
}
