package nbs

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/nb"
)

var filterConfig = frameworkfilter.Config{
	"label":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"tags":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"ipv4":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"region": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},

	"hostname":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"ipv6":                 {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"client_conn_throttle": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"nodebalancers": schema.ListNestedAttribute{
			Description: "The returned list of NodeBalancers.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nb.DataSourceAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
