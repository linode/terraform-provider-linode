package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"site_type":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"capabilities": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"country":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"regions": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"country": schema.StringAttribute{
						Description: "The country where this Region resides.",
						Computed:    true,
					},
					"id": schema.StringAttribute{
						Description: "The unique ID of this Region.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "Detailed location information for this Region, including city, state or region, and country.",
						Computed:    true,
					},
					"site_type": schema.StringAttribute{
						Description: "The type of this Region.",
						Computed:    true,
					},
					"capabilities": schema.SetAttribute{
						Description: "A list of capabilities of this Region.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"status": schema.StringAttribute{
						Description: "This region’s current operational status.",
						Computed:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"resolvers": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ipv4": schema.StringAttribute{
									Description: "The IPv4 addresses for this region’s DNS resolvers, separated by commas.",
									Computed:    true,
								},
								"ipv6": schema.StringAttribute{
									Description: "The IPv6 addresses for this region’s DNS resolvers, separated by commas.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	},
}
