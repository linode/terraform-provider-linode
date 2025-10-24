package vpcsubnet

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
)

var LinodeInterfaceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":        types.Int64Type,
		"config_id": types.Int64Type,
		"active":    types.BoolType,
	},
}

var LinodeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.Int64Type,
		"interfaces": types.ListType{
			ElemType: LinodeInterfaceObjectType,
		},
	},
}

var ResourceSchemaIPv6NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "An existing IPv6 prefix owned by the current account or a " +
				"forward slash (/) followed by a valid prefix length. " +
				"If unspecified, a range with the default prefix will be " +
				"allocated for this VPC.",
			Optional:   true,
			Computed:   true,
			CustomType: customtypes.LinodeAutoAllocRangeType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"allocated_range": schema.StringAttribute{
			Description: "The IPv6 range assigned to this subnet.",
			Computed:    true,
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the VPC Subnet.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
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
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"ipv6": schema.ListNestedAttribute{
			Description: "The IPv6 ranges of this subnet.",
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.RequiresReplace(),
				listplanmodifier.UseStateForUnknown(),
			},
			NestedObject: ResourceSchemaIPv6NestedObject,
		},

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
		"linodes": schema.ListAttribute{
			Computed:    true,
			ElementType: LinodeObjectType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
