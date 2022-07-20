package domainzonefile

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"domain_id": {
		Type:        schema.TypeInt,
		Description: "The domain's ID.",
		Required:    true,
	},
	"zone_file": {
		Type:        schema.TypeList,
		Description: "Lines of the zone file for the last rendered zone for this domain.",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
}
