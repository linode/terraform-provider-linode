package vpcsubnets

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/vpcsubnet"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"ipv4":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"vpc_id": schema.Int64Attribute{
			Description: "The id of the parent VPC for the list of VPC subnets",
			Required:    true,
		},
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"vpc_subnets": schema.ListNestedBlock{
			Description: "The returned list of subnets under a VPC.",
			NestedObject: schema.NestedBlockObject{
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

					"linodes": vpcsubnet.LinodesSchema,
				},
			},
		},
	},
}
