package balancer

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceSchemaTransfer = map[string]*schema.Schema{
	"in": {
		Type:        schema.TypeFloat,
		Description: "The total transfer, in MB, used by this NodeBalancer this month",
		Computed:    true,
	},
	"out": {
		Type:        schema.TypeFloat,
		Description: "The total inbound transfer, in MB, used for this NodeBalancer this month",
		Computed:    true,
	},
	"total": {
		Type:        schema.TypeFloat,
		Description: "The total outbound transfer, in MB, used for this NodeBalancer this month",
		Computed:    true,
	},
}

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode NodeBalancer.",
		Optional:    true,
	},
	"region": {
		Type:         schema.TypeString,
		Description:  "The region where this NodeBalancer will be deployed.",
		Required:     true,
		ForceNew:     true,
		InputDefault: "us-east",
	},
	"client_conn_throttle": {
		Type:         schema.TypeInt,
		Description:  "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
		ValidateFunc: validation.IntBetween(0, 20),
		Optional:     true,
		Default:      0,
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
		Optional:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
}
