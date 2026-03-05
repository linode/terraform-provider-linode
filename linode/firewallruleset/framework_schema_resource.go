package firewallruleset

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
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
			Description: "Controls whether traffic is accepted or dropped by this rule.",
			Required:    true,
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
		"description": schema.StringAttribute{
			Description: "Used to describe this rule. For display purposes only.",
			Optional:    true,
		},
		"ports": schema.StringAttribute{
			Description: "A string representation of ports and/or port ranges (e.g. \"443\" or \"80-90, 91\").",
			Optional:    true,
		},
		"ipv4": schema.ListAttribute{
			Description: "A list of IPv4 addresses, CIDRs, or prefix list tokens this rule applies to.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"ipv6": schema.ListAttribute{
			Description: "A list of IPv6 addresses, networks, or prefix list tokens this rule applies to.",
			Optional:    true,
			ElementType: types.StringType,
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Description: "Manages a Linode Firewall Rule Set.",
	Blocks: map[string]schema.Block{
		"rules": schema.ListNestedBlock{
			Description:  "An ordered list of firewall rules in this rule set.",
			NestedObject: ruleNestedObject,
		},
	},
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this Rule Set.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label for the Rule Set.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 32),
			},
		},
		"description": schema.StringAttribute{
			Description: "A description for this Rule Set.",
			Optional:    true,
		},
		"type": schema.StringAttribute{
			Description: "Whether the rules in this set are inbound or outbound. One of: inbound, outbound.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("inbound", "outbound"),
			},
		},
		"is_service_defined": schema.BoolAttribute{
			Description: "Whether this Rule Set is read-only and defined by a Linode service.",
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of this Rule Set, incremented on each update.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When the Rule Set was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When the Rule Set was last updated.",
			Computed:    true,
		},
	},
}
