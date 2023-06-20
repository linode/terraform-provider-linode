package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
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

var instanceTypeSchema = schema.NestedBlockObject{
	Attributes: instancetype.Attributes,
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
			Description:  "The returned list of instance types.",
			NestedObject: instanceTypeSchema,
		},
	},
}
