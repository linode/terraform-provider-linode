package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
)

var ResourceSchemaIPv6NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 range assigned to this VPC.",
			Optional:    true,
			CustomType:  customtypes.LinodeAutoAllocRangeType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"allocated_range": schema.StringAttribute{
			Description: "The IPv6 range assigned to this VPC.",
			Computed:    true,
		},
		"allocation_class": schema.StringAttribute{
			Description: "The labeled IPv6 Inventory that the VPC Prefix should be allocated from.",
			Optional:    true,
			WriteOnly:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the VPC.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC. Only contains ascii letters, digits and dashes",
			Required:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region of the VPC.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": schema.StringAttribute{
			Description: "The user-defined description of this VPC.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ipv6": schema.ListNestedAttribute{
			Description: "The IPv6 configuration of this VPC.",
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
				listplanmodifier.RequiresReplace(),
			},
			NestedObject: ResourceSchemaIPv6NestedObject,
		},

		"created": schema.StringAttribute{
			Description: "The date and time when the VPC was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC was updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
