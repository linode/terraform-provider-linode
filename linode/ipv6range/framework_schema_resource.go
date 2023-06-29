package ipv6range

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"prefix_length": schema.Int64Attribute{
			Description:   "The prefix length of the IPv6 range.",
			Required:      true,
			PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			Validators:    []validator.Int64{int64validator.OneOf(56, 64)},
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to assign this range to.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.ConflictsWith(path.Expressions{
					path.MatchRoot("route_target"),
				}...),
			},
		},
		"route_target": schema.StringAttribute{
			Description:   "The IPv6 SLAAC address to assign this range to.",
			Optional:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.Expressions{
					path.MatchRoot("linode_id"),
				}...),
			},
		},
		"is_bgp": schema.BoolAttribute{
			Description: "Whether this IPv6 range is shared.",
			Computed:    true,
		},
		"linodes": schema.SetAttribute{
			Description: "A list of Linodes targeted by this IPv6 range." +
				"Includes Linodes with IP sharing.",
			ElementType: types.Int64Type,
			Computed:    true,
		},
		"range": schema.StringAttribute{
			Description: "The IPv6 range of addresses in this pool.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region for this range of IPv6 addresses.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID for this Resource.",
			Computed:    true,
		},
	},
}
