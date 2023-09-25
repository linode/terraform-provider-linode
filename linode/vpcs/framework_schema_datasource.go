package vpcs

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/vpc"
)

var filterConfig = frameworkfilter.Config{
	"id":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"description": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"region":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"vpcs": schema.ListNestedBlock{
			Description: "The returned list of VPCs.",
			NestedObject: schema.NestedBlockObject{
				Attributes: vpc.VPCAttrs,
			},
		},
	},
}
