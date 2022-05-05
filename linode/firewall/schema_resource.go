package firewall

import (
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var resourceRuleSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: `Used to identify this rule. For display purposes only.`,
		Required:    true,
	},
	"action": {
		Type: schema.TypeString,
		Description: "Controls whether traffic is accepted or dropped by this rule. Overrides the Firewall’s " +
			"inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.",
		Required: true,
	},
	"ports": {
		Type:        schema.TypeString,
		Description: `A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").`,
		Optional:    true,
	},
	"protocol": {
		Type:        schema.TypeString,
		Description: "The network protocol this rule controls.",
		StateFunc: func(val interface{}) string {
			return strings.ToUpper(val.(string))
		},
		Required: true,
	},
	"ipv4": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "A list of IP addresses, CIDR blocks, or 0.0.0.0/0 (to allow all) this rule applies to.",
		Optional:    true,
	},
	"ipv6": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
				err := helper.ValidateIPv6Range(i.(string))
				if err != nil {
					return diag.FromErr(err)
				}

				return nil
			},
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// We handle validation separately
				result, _ := helper.CompareIPv6Ranges(oldValue, newValue)
				return result
			},
		},
		Description: "A list of IPv6 addresses or networks this rule applies to.",
		Optional:    true,
	},
}

var resourceDeviceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The ID of the firewall device.",
		Computed:    true,
	},
	"entity_id": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The ID of the underlying entity for the firewall device (e.g. the Linode's ID).",
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The type of firewall device.",
		Computed:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the underlying entity for the firewall device.",
		Computed:    true,
	},
	"url": {
		Type:        schema.TypeString,
		Description: "The URL of the underlying entity for the firewall device.",
		Computed:    true,
	},
}

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type: schema.TypeString,
		Description: "The label for the Firewall. For display purposes only. If no label is provided, a " +
			"default will be assigned.",
		Required:     true,
		ValidateFunc: validation.StringLenBetween(3, 32),
	},
	"tags": {
		Type:        schema.TypeSet,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Set:         schema.HashString,
	},
	"disabled": {
		Type:        schema.TypeBool,
		Description: "If true, the Firewall is inactive.",
		Optional:    true,
		Default:     false,
	},
	"inbound": {
		Type:        schema.TypeList,
		Elem:        resourceFirewallRules(),
		Description: "A firewall rule that specifies what inbound network traffic is allowed.",
		Optional:    true,
	},
	"inbound_policy": {
		Type: schema.TypeString,
		Description: "The default behavior for inbound traffic. This setting can be overridden by updating " +
			"the inbound.action property for an individual Firewall Rule.",
		Required: true,
	},
	"outbound": {
		Type:        schema.TypeList,
		Elem:        resourceFirewallRules(),
		Description: "A firewall rule that specifies what outbound network traffic is allowed.",
		Optional:    true,
	},
	"outbound_policy": {
		Type: schema.TypeString,
		Description: "The default behavior for outbound traffic. This setting can be overridden by updating " +
			"the outbound.action property for an individual Firewall Rule.",
		Required: true,
	},
	"linodes": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeInt},
		Description: "The IDs of Linodes to apply this firewall to.",
		Optional:    true,
		Computed:    true,
		Set:         schema.HashInt,
	},
	"devices": {
		Type:        schema.TypeList,
		Elem:        resourceFirewallDevice(),
		Computed:    true,
		Description: "The devices associated with this firewall.",
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The status of the firewall.",
		Computed:    true,
	},
}
