package natgatewaytypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

// TODO: Replace with API filterable fields once they are available
var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"types": schema.ListNestedAttribute{
			Description: "The returned list of NAT Gateway types.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The unique ID assigned to this NAT Gateway type.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The NAT Gateway type's label.",
						Computed:    true,
					},
					"price": schema.ListAttribute{
						Description: "Cost in US dollars, broken down into hourly and monthly charges.",
						Computed:    true,
						ElementType: helper.PriceObjectType,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
