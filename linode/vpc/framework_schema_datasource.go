package vpc

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcsubnet"
)

var DataSourceSchemaIPv6NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 range assigned to this VPC.",
			Computed:    true,
		},
	},
}

var DataSourceSchemaIPv4NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv4 range assigned to this VPC.",
			Computed:    true,
		},
	},
}

var DataSourceSchemaSubnetIPv6NestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "An IPv6 range allocated to this subnet.",
			Computed:    true,
		},
	},
}

var DataSourceSchemaSubnetNestedObject = schema.NestedAttributeObject{
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
			NestedObject: DataSourceSchemaSubnetIPv6NestedObject,
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

var VPCAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The id of the VPC.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the VPC.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "The user-defined description of this VPC.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "The region of the VPC.",
		Computed:    true,
	},
	"vpc_type": schema.StringAttribute{
		Description: "The type of the VPC ('regular' or 'rdma'). " +
			"Omitted if the requesting account does not have access to the GPUDirect RDMA functionality.",
		Computed: true,
	},
	"ipv4": schema.ListNestedAttribute{
		Description:  "The IPv4 configuration of this VPC.",
		Computed:     true,
		NestedObject: DataSourceSchemaIPv4NestedObject,
	},
	"ipv6": schema.ListNestedAttribute{
		Description:  "The IPv6 configuration of this VPC.",
		Computed:     true,
		NestedObject: DataSourceSchemaIPv6NestedObject,
	},
	"subnets": schema.ListNestedAttribute{
		Description:  "A list of subnets under this VPC.",
		Computed:     true,
		NestedObject: DataSourceSchemaSubnetNestedObject,
	},
	"created": schema.StringAttribute{
		Description: "The date and time when the VPC was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "The date and time when the VPC was updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: VPCAttrs,
}
