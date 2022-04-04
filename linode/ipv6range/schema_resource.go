package ipv6range

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceSchema = map[string]*schema.Schema{
	"prefix_length": {
		Type:         schema.TypeInt,
		Description:  "The prefix length of the IPv6 range.",
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.IntInSlice([]int{56, 64}),
	},
	"linode_id": {
		Type:          schema.TypeInt,
		Description:   "The ID of the Linode to assign this range to.",
		Optional:      true,
		ConflictsWith: []string{"route_target"},
	},
	"route_target": {
		Type:          schema.TypeString,
		Description:   "The IPv6 SLAAC address to assign this range to.",
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"linode_id"},
	},
	"is_bgp": {
		Type:        schema.TypeBool,
		Description: "Whether this IPv6 range is shared.",
		Computed:    true,
	},
	"linodes": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Description: "A list of Linodes targeted by this IPv6 range. Includes Linodes with IP sharing.",
		Computed:    true,
	},
	"range": {
		Type:        schema.TypeString,
		Description: "The IPv6 range of addresses in this pool.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region for this range of IPv6 addresses.",
		Computed:    true,
	},
}
