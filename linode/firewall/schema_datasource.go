package firewall

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceRuleSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: `Used to identify this rule. For display purposes only.`,
		Computed:    true,
	},
	"action": {
		Type: schema.TypeString,
		Description: "Controls whether traffic is accepted or dropped by this rule. Overrides the Firewall’s " +
			"inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.",
		Computed: true,
	},
	"ports": {
		Type:        schema.TypeString,
		Description: `A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").`,
		Computed:    true,
	},
	"protocol": {
		Type:        schema.TypeString,
		Description: "The network protocol this rule controls.",
		Computed:    true,
	},
	"ipv4": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of IP addresses, CIDR blocks, or 0.0.0.0/0 (to allow all) this rule applies to.",
		Computed:    true,
	},
	"ipv6": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of IPv6 addresses or networks this rule applies to.",
		Computed:    true,
	},
}

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID assigned to this Firewall.",
		Required:    true,
	},
	"label": {
		Type: schema.TypeString,
		Description: "The label for the Firewall. For display purposes only. If no label is provided, a " +
			"default will be assigned.",
		Computed: true,
	},
	"tags": {
		Type:        schema.TypeSet,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Set:         schema.HashString,
	},
	"disabled": {
		Type:        schema.TypeBool,
		Description: "If true, the Firewall is inactive.",
		Computed:    true,
	},
	"inbound": {
		Type:        schema.TypeList,
		Elem:        dataSourceFirewallRules(),
		Description: "A firewall rule that specifies what inbound network traffic is allowed.",
		Computed:    true,
	},
	"inbound_policy": {
		Type: schema.TypeString,
		Description: "The default behavior for inbound traffic. This setting can be overridden by updating " +
			"the inbound.action property for an individual Firewall Rule.",
		Computed: true,
	},
	"outbound": {
		Type:        schema.TypeList,
		Elem:        dataSourceFirewallRules(),
		Description: "A firewall rule that specifies what outbound network traffic is allowed.",
		Computed:    true,
	},
	"outbound_policy": {
		Type: schema.TypeString,
		Description: "The default behavior for outbound traffic. This setting can be overridden by updating " +
			"the outbound.action property for an individual Firewall Rule.",
		Computed: true,
	},
	"linodes": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Description: "The IDs of Linodes to apply this firewall to.",
		Computed:    true,
		Set:         schema.HashInt,
	},
	"devices": {
		Type:        schema.TypeList,
		Elem:        resourceFirewallDevice(),
		Description: "The devices associated with this firewall.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The status of the firewall.",
		Computed:    true,
	},
}
