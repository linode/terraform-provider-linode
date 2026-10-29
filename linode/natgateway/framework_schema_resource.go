package natgateway

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var vpcSubnet = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the attached VPC subnet.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the attached vpc_subnet.",
			Computed:    true,
		},
		"url": schema.StringAttribute{
			Description: "The API URL of the attached VPC subnet.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC subnet.",
			Computed:    true,
		},
		"vpc_id": schema.StringAttribute{
			Description: "The unique ID of the parent VPC.",
			Computed:    true,
		},
		"vpc_label": schema.StringAttribute{
			Description: "The label of the parent VPC.",
			Computed:    true,
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:    true,
			Description: "The unique ID of the NAT Gateway.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The region where the NAT Gateway will be created.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "The label for the NAT Gateway.",
		},
		"default_ports_per_interface": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The default number of ports per interface for the NAT Gateway.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"use_autoscaling": schema.BoolAttribute{
			Optional: true,
			Description: "Whether the NAT Gateway automatically scales its public IPv4 addresses. " +
				"If omitted, defaults to true when addresses is empty, otherwise false. Create-only and Write-only.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"vpc_subnet_id": schema.Int64Attribute{
			Optional:    true,
			Description: "The ID of the VPC subnet for this NAT Gateway. Create-only.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"address_autoscale_max": schema.Int64Attribute{
			Computed:    true,
			Description: "The maximum number of auto-allocated IP addresses allowed for this NAT Gateway instance. Reflects the state of use_autoscaling.",
		},
		"portset_assignments": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of portsets currently assigned to interfaces behind this NAT Gateway, excluding system-reserved ports.",
		},
		"portset_capacity": schema.Int64Attribute{
			Computed:    true,
			Description: "The total number of portsets available to allocate across all current and potential autoscaled public IP addresses assigned to this NAT Gateway, excluding system-reserved ports.",
		},
		"vpc_subnet": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  vpcSubnet.Attributes,
			Description: "Information about the VPC subnet attached to this NAT Gateway, if there is one attached.",
		},
		"created": schema.StringAttribute{
			Description: "When this NAT Gateway was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this NAT Gateway was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"addresses": schema.SetNestedAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The public IP addresses assigned to this NAT Gateway. Must be empty when use_autoscaling is true.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The public IP address assigned to this NAT Gateway.",
						Required:    true,
					},
				},
			},
		},
	},
}
