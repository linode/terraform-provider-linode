package networkingip

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the IPv4 address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to allocate an IPv4 address for. Required when reserved is false or not set.",
			Optional:    true,
			Computed:    true,
		},
		"reserved": schema.BoolAttribute{
			Description: "Whether the IPv4 address should be reserved.",
			Optional:    true,
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region for the reserved IPv4 address. Required when reserved is true and linode_id is not set.",
			Optional:    true,
			Computed:    true,
		},
		"public": schema.BoolAttribute{
			Description: "Whether the IPv4 address is public or private.",
			Optional:    true,
			Computed:    true,
		},
		"address": schema.StringAttribute{
			Description: "The allocated IPv4 address.",
			Computed:    true,
			Optional:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of IP address (ipv4).",
			Computed:    true,
			Optional:    true,
		},
	},
}
