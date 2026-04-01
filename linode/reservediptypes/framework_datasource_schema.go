package reservediptypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"label": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order_by": filterConfig.OrderBySchema(),
		"order":    filterConfig.OrderSchema(),
		"types": schema.ListNestedAttribute{
			Description: "The returned list of Reserved IP types.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The unique ID assigned to this Reserved IP Type.",
						Required:    true,
					},
					"label": schema.StringAttribute{
						Description: "The Reserved IP Type's label.",
						Computed:    true,
						Optional:    true,
					},
					"price": schema.ListAttribute{
						Description: "Cost in US dollars, broken down into hourly and monthly charges.",
						Computed:    true,
						ElementType: helper.PriceObjectType,
					},
					"region_prices": schema.ListAttribute{
						Description: "A list of region-specific prices for this Reserved IP Type.",
						Computed:    true,
						ElementType: helper.RegionPriceObjectType,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
