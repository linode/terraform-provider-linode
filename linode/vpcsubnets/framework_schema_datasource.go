package vpcsubnets

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcsubnet"
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
		"vpc_subnets": schema.ListNestedAttribute{
			Description: "The returned list of subnets under a VPC.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
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
						Description: "The IPv6 ranges of this subnet.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"range": schema.StringAttribute{
									Description: "An IPv6 range allocated to this subnet.",
									Computed:    true,
									CustomType:  customtypes.LinodeAutoAllocRangeType{},
								},
							},
						},
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
					"linodes": schema.ListAttribute{
						Computed:    true,
						ElementType: vpcsubnet.LinodeObjectType,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
