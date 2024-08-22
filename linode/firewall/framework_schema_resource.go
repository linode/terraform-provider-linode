package firewall

import (
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	linodeplanmodifiers "github.com/linode/terraform-provider-linode/v2/linode/helper/planmodifiers"
)

var ruleNestedObject = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "Used to identify this rule. For display purposes only.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 32),
			},
		},
		"action": schema.StringAttribute{
			Description: "Controls whether traffic is accepted or dropped by this rule. " +
				"Overrides the Firewall's inbound_policy if this is an inbound rule, or " +
				"the outbound_policy if this is an outbound rule.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACCEPT", "DROP"),
			},
		},
		"protocol": schema.StringAttribute{
			Description: "The network protocol this rule controls.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.TCP),
					string(linodego.UDP),
					string(linodego.ICMP),
					string(linodego.IPENCAP),
				),
			},
		},
		"ports": schema.StringAttribute{
			Description: "A string representation of ports and/or port ranges " +
				"(i.e. \"443\" or \"80-90, 91\").",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"ipv4": schema.ListAttribute{
			Description: "A list of CIDR blocks or 0.0.0.0/0 (to allow all) this rule applies to.",
			Optional:    true,
			ElementType: cidrtypes.IPv4PrefixType{},
		},
		"ipv6": schema.ListAttribute{
			Description: "A list of IPv6 addresses or networks this rule applies to.",
			Optional:    true,
			ElementType: cidrtypes.IPv6PrefixType{},
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Blocks: map[string]schema.Block{
		"inbound": schema.ListNestedBlock{
			Description:  "A firewall rule that specifies what inbound network traffic is allowed.",
			NestedObject: ruleNestedObject,
		},
		"outbound": schema.ListNestedBlock{
			Description:  "A firewall rule that specifies what outbound network traffic is allowed.",
			NestedObject: ruleNestedObject,
		},
	},
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this Object Storage key.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label for the Firewall. For display purposes only." +
				" If no label is provided, a default will be assigned.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 32),
			},
		},
		"tags": schema.SetAttribute{
			Description: "An array of tags applied to the firewall. Tags are for organizational purposes only.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Set{
				linodeplanmodifiers.CaseInsensitiveSet(),
			},
			Default: helper.EmptySetDefault(types.StringType),
		},
		"disabled": schema.BoolAttribute{
			Description: "If true, the Firewall is inactive.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"inbound_policy": schema.StringAttribute{
			Description: "The default behavior for inbound traffic. This setting can be overridden by updating " +
				"the inbound.action property for an individual Firewall Rule.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACCEPT", "DROP"),
			},
		},
		"outbound_policy": schema.StringAttribute{
			Description: "The default behavior for outbound traffic. This setting can be overridden by updating " +
				"the outbound.action property for an individual Firewall Rule.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("ACCEPT", "DROP"),
			},
		},
		"linodes": schema.SetAttribute{
			Description: "The IDs of Linodes to apply this firewall to.",
			Optional:    true,
			Computed:    true,
			ElementType: types.Int64Type,
			Default:     helper.EmptySetDefault(types.Int64Type),
		},
		"nodebalancers": schema.SetAttribute{
			Description: "The IDs of NodeBalancers to apply this firewall to.",
			Optional:    true,
			Computed:    true,
			ElementType: types.Int64Type,
			Default:     helper.EmptySetDefault(types.Int64Type),
		},
		"devices": schema.ListAttribute{
			Description: "The devices associated with this firewall.",
			Computed:    true,
			ElementType: deviceObjectType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			Description: "The status of the firewall.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this firewall was created",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this firewall was last updated",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
