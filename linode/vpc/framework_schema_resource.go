package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/stringplanmodifiers"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcsubnet"
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
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

var ResourceSchemaIPv4NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv4 range assigned to this VPC.",
			Required:    true,
		},
	},
}

var ResourceSchemaSubnetIPv6NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "An IPv6 range allocated to this subnet.",
			Computed:    true,
		},
	},
}

var ResourceSchemaSubnetNestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC Subnet.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC Subnet.",
			Computed:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The IPv4 range of this subnet in CIDR format.",
			Computed:    true,
		},
		"ipv6": schema.ListNestedAttribute{
			Description:  "The IPv6 ranges of this subnet.",
			Computed:     true,
			NestedObject: ResourceSchemaSubnetIPv6NestedObject,
		},
		"linodes": schema.ListAttribute{
			Description: "A list of Linodes assigned to this subnet.",
			Computed:    true,
			ElementType: vpcsubnet.LinodeObjectType,
		},
		"databases": schema.ListAttribute{
			Description: "A list of Managed Databases assigned to this subnet.",
			Computed:    true,
			ElementType: vpcsubnet.DatabaseObjectType,
		},
		"nodebalancers": schema.ListAttribute{
			Description: "A list of NodeBalancers assigned to this subnet.",
			Computed:    true,
			ElementType: vpcsubnet.NodebalancerObjectType,
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the VPC Subnet was updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
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
		"vpc_type": schema.StringAttribute{
			Description: "The type of the VPC. Can be either 'regular' or 'rdma'. " +
				"Defaults to 'regular'. The 'rdma' type may not be available to all users.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
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
		"ipv4": schema.ListNestedAttribute{
			Description: "The IPv4 configuration of this VPC.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			NestedObject: ResourceSchemaIPv4NestedObject,
		},
		"subnets": schema.ListNestedAttribute{
			Description: "A list of subnets under this VPC.",
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			NestedObject: ResourceSchemaSubnetNestedObject,
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifiers.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("label"),
					path.MatchRoot("description"),
					path.MatchRoot("ipv4"),
				),
			},
		},
	},
}
