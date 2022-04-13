package ipv6range

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"range": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The IPv6 range to retrieve information about.",
	},

	"is_bgp": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether this IPv6 range is shared.",
	},
	"linodes": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Computed:    true,
		Description: "A list of Linodes targeted by this IPv6 range. Includes Linodes with IP sharing.",
	},
	"prefix": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The prefix length of the address, denoting how many addresses can be assigned from this range.",
	},
	"region": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The region for this range of IPv6 addresses.",
	},
}
