package vpcsubnets

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcsubnet"
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
					"databases": schema.ListAttribute{
						Computed:    true,
						ElementType: vpcsubnet.DatabaseObjectType,
					},
					"nodebalancers": schema.ListAttribute{
						Computed:    true,
						ElementType: vpcsubnet.NodebalancerObjectType,
					},
					"natgateway": schema.SingleNestedAttribute{
						Description: "The details of the NAT Gateway currently assigned to this subnet.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Description: "The ID of the assigned NAT Gateway.",
								Computed:    true,
							},
							"label": schema.StringAttribute{
								Description: "The label of the assigned NAT Gateway.",
								Computed:    true,
							},
							"addresses": schema.ListAttribute{
								Description: "The addresses of the assigned NAT Gateway.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"portset_assignments": schema.Int64Attribute{
								Description: "The number of portset assignments on the assigned NAT Gateway.",
								Computed:    true,
							},
							"portset_capacity": schema.Int64Attribute{
								Description: "The total portset capacity of the assigned NAT Gateway.",
								Computed:    true,
							},
							"portsets": schema.ListNestedAttribute{
								Description: "The portsets of the assigned NAT Gateway.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"address": schema.StringAttribute{
											Description: "The address the portset is assigned to.",
											Computed:    true,
										},
										"ports": schema.ListNestedAttribute{
											Description: "The port ranges in the portset.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"start": schema.Int64Attribute{
														Description: "The first port in the range.",
														Computed:    true,
													},
													"end": schema.Int64Attribute{
														Description: "The last port in the range.",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
