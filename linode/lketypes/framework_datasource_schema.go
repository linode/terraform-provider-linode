package lketypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var lkeTypeSchema = schema.NestedBlockObject{}

var filterConfig = frameworkfilter.Config{
	"label":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"transfer": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
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
			Computed:    true,
			Description: "The returned list of LKE types.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: helper.GetPricingTypeAttributes("LKE Type"),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
