package nb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of the Linode NodeBalancer.",
		Required:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode NodeBalancer.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region where this NodeBalancer will be deployed.",
		Computed:    true,
	},
	"client_conn_throttle": {
		Type:        schema.TypeInt,
		Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
		Computed:    true,
	},
	"hostname": {
		Type:        schema.TypeString,
		Description: "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
		Computed:    true,
	},
	"ipv4": {
		Type:        schema.TypeString,
		Description: "The Public IPv4 Address of this NodeBalancer",
		Computed:    true,
	},
	"ipv6": {
		Type:        schema.TypeString,
		Description: "The Public IPv6 Address of this NodeBalancer",
		Computed:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this NodeBalancer was created.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this NodeBalancer was last updated.",
		Computed:    true,
	},
	"transfer": {
		Type:        schema.TypeList,
		Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
		Computed:    true,
		Elem:        resourceTransfer(),
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
}
