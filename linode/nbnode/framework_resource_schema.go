package nbnode

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the NodeBalancer node resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"nodebalancer_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer to access.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"config_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancerConfig to access.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label for this node. This is for display purposes only.",
			Required:    true,
		},
		"address": schema.StringAttribute{
			Description: "The private IP Address and port (IP:PORT) where this backend can be reached. " +
				"This must be a private IP address.",
			Required: true,
		},
		"weight": schema.Int64Attribute{
			Description: "Used when picking a backend to serve a request and is not pinned to a single backend " +
				"yet. Nodes with a higher weight will receive more traffic. (1-255)",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 255),
			},
		},
		"mode": schema.StringAttribute{
			Description: "The mode this NodeBalancer should use when sending traffic to this backend. If set to " +
				"`accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. " +
				"If set to `drain` this backend will not receive new traffic, but connections already pinned to it will " +
				"continue to be routed to it. If set to `backup` this backend will only accept traffic if all other " +
				"nodes are down.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("accept", "reject", "drain", "backup"),
			},
		},
		"status": schema.StringAttribute{
			Description: "The current status of this node, based on the configured checks of its NodeBalancer " +
				"Config. (unknown, UP, DOWN)",
			Computed: true,
		},
	},
}
