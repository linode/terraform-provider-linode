package networkingipassignment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "The ID of the IP assignment operation.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The region for the IP assignments.",
		},
		"assignments": schema.ListNestedAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
				listplanmodifier.RequiresReplace(),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Required:    true,
						Description: "The IPv4 address or IPv6 range to assign.",
					},
					"linode_id": schema.Int64Attribute{
						Required:    true,
						Description: "The ID of the Linode to which the IP address will be assigned.",
					},
					"reserved": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this IP address is a reserved IP.",
					},
					"tags": schema.SetAttribute{
						Computed:    true,
						Description: "A set of tags associated with this IP address.",
						ElementType: types.StringType,
					},
				},
			},
		},
	},
}
