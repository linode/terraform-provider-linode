package balancernode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceSchema = map[string]*schema.Schema{
	"nodebalancer_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the NodeBalancer to access.",
		Required:    true,
		ForceNew:    true,
	},
	"config_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the NodeBalancerConfig to access.",
		Required:    true,
		ForceNew:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The label for this node. This is for display purposes only.",
		Required:    true,
	},
	"weight": {
		Type: schema.TypeInt,
		Description: "Used when picking a backend to serve a request and is not pinned to a single backend " +
			"yet. Nodes with a higher weight will receive more traffic. (1-255)",
		ValidateFunc: validation.IntBetween(1, 255),
		Optional:     true,
		Computed:     true,
	},
	"mode": {
		Type: schema.TypeString,
		Description: "The mode this NodeBalancer should use when sending traffic to this backend. If set to " +
			"`accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. " +
			"If set to `drain` this backend will not receive new traffic, but connections already pinned to it will " +
			"continue to be routed to it. If set to `backup` this backend will only accept traffic if all other " +
			"nodes are down.",
		ValidateFunc: validation.StringInSlice([]string{"accept", "reject", "drain", "backup"}, false),
		Optional:     true,
		Computed:     true,
	},
	"address": {
		Type: schema.TypeString,
		Description: "The private IP Address and port (IP:PORT) where this backend can be reached. " +
			"This must be a private IP address.",
		Required: true,
	},
	"status": {
		Type: schema.TypeString,
		Description: "The current status of this node, based on the configured checks of its NodeBalancer " +
			"Config. (unknown, UP, DOWN)",
		Computed: true,
	},
}
