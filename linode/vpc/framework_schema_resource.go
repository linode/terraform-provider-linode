package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC. Only contains ascii letters, digits and dashes",
			Required:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region of the VPC.",
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: "The user-defined description of this VPC.",
			Optional:    true,
			Computed:    true,
		},
		"subnets": schema.ListAttribute{
			Description: "A list of subnets under this VPC.",
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			ElementType: subnetObjectType,
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the VPC was created.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC was updated.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
	},
	Blocks: map[string]schema.Block{
		"subnets_create_options": schema.ListNestedBlock{
			Description: "A list of create options to create a list of VPC subnets. " +
				"Configure this block when creating if you want to have a list of VPC subnets under this VPC.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"label": schema.StringAttribute{
						Description: "The label of the VPC subnet.",
						Required:    true,
					},
					"ipv4": schema.StringAttribute{
						Description: "The IPv4 range of this subnet in CIDR format.",
						Required:    true,
					},
				},
			},
		},
	},
}
