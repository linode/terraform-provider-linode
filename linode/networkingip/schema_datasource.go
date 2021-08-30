package networkingip

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"address": {
		Type:        schema.TypeString,
		Description: "The IP address.",
		Required:    true,
	},
	"gateway": {
		Type:        schema.TypeString,
		Description: "The default gateway for this address.",
		Computed:    true,
	},
	"subnet_mask": {
		Type:        schema.TypeString,
		Description: "The mask that separates host bits from network bits for this address.",
		Computed:    true,
	},
	"prefix": {
		Type:        schema.TypeInt,
		Description: "The number of bits set in the subnet mask.",
		Computed:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The type of address this is (ipv4, ipv6, ipv6/pool, ipv6/range).",
		Computed:    true,
	},
	"public": {
		Type:        schema.TypeBool,
		Description: "Whether this is a public or private IP address.",
		Computed:    true,
	},
	"rdns": {
		Type: schema.TypeString,
		Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to " +
			"a default value provided by Linode if not explicitly set.",
		Computed: true,
	},
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Linode this address currently belongs to.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The Region this IP address resides in.",
		Computed:    true,
	},
}
