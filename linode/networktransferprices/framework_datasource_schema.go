package networktransferprices

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var networkTransferPriceSchema = schema.NestedBlockObject{
	Attributes: helper.GetPricingTypeAttributes("Network Transfer Price"),
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
			Description:  "The returned list of Network Transfer Prices.",
			NestedObject: networkTransferPriceSchema,
		},
	},
}
