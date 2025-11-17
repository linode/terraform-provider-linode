package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/instancetype"
)

var filterConfig = frameworkfilter.Config{
	"class":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"disk":        {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"gpus":        {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"memory":      {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"network_out": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"transfer":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"vcpus":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
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
			Description: "The returned list of instance types.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: instancetype.Attributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
