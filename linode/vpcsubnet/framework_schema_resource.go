package vpcsubnet

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var subnetLinodesAttribute = schema.ListNestedAttribute{
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of a Linode attached to this subnet.",
			},
			"interfaces": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The ID of an interface that references this VPC subnet.",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this interface is active",
						},
					},
				},
				Computed: true,
			},
		},
	},
	Computed: true,
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC Subnet.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The id of the parent VPC for this VPC Subnet",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC subnet.",
			Required:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The IPv4 range of this subnet in CIDR format.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"linodes": subnetLinodesAttribute,
		"created": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
