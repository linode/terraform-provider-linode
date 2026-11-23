package natgatewaytypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"label": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
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
				Attributes: helper.GetPricingTypeAttributes("NAT Gateway Type"),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
