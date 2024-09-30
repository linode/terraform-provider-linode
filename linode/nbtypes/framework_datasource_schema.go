package nbtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var nodebalancerTypeSchema = schema.NestedBlockObject{
	Attributes: helper.GetPricingTypeAttributes("Node Balancer Type"),
}

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
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"types": schema.ListNestedBlock{
			Description:  "The returned list of Node Balancer types.",
			NestedObject: nodebalancerTypeSchema,
		},
	},
}
