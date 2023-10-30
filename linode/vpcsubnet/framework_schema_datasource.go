package vpcsubnet

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The id of the VPC Subnet.",
			Required:    true,
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The id of the parent VPC for this VPC Subnet",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC Subnet.",
			Computed:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The IPv4 range of this subnet in CIDR format.",
			Computed:    true,
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
	Blocks: map[string]schema.Block{
		"linodes": LinodesSchema,
	},
}
