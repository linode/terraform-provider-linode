package reservedip

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

var assignedEntityType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.Int64Type,
		"label": types.StringType,
		"type":  types.StringType,
		"url":   types.StringType,
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the reserved IP address, which is the address itself.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region in which to reserve the IP address.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"address": schema.StringAttribute{
			Description: "The reserved IPv4 address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"gateway": schema.StringAttribute{
			Description: "The default gateway for this address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"subnet_mask": schema.StringAttribute{
			Description: "The mask that separates host bits from network bits for this address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"prefix": schema.Int64Attribute{
			Description: "The number of bits set in the subnet mask.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "The type of IP address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"public": schema.BoolAttribute{
			Description: "Whether this is a public or private IP address.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode this address is currently assigned to, if any.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"reserved": schema.BoolAttribute{
			Description: "Whether this IP address is reserved.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"tags": schema.SetAttribute{
			Description: "Tags assigned to this reserved IP address.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
		},
		"vpc_nat_1_1": schema.ListAttribute{
			Description: "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
			Computed:    true,
			ElementType: instancenetworking.VPCNAT1To1Type,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"assigned_entity": schema.ListAttribute{
			Description: "The entity this reserved IP address is currently assigned to, if any.",
			Computed:    true,
			ElementType: assignedEntityType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
