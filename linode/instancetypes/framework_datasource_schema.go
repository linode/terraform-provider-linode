package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

var filterConfig = frameworkfilter.Config{
	"class":       {APIFilterable: true},
	"disk":        {APIFilterable: true},
	"gpus":        {APIFilterable: true},
	"label":       {APIFilterable: true},
	"memory":      {APIFilterable: true},
	"network_out": {APIFilterable: true},
	"transfer":    {APIFilterable: true},
	"vcpus":       {APIFilterable: true},
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
