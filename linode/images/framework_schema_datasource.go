package images

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/image"
)

var filterConfig = frameworkfilter.Config{
	"deprecated": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"is_public":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"label":      {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"size":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"type":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"vendor":     {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},

	"created_by":  {TypeFunc: frameworkfilter.FilterTypeString},
	"id":          {TypeFunc: frameworkfilter.FilterTypeString},
	"status":      {TypeFunc: frameworkfilter.FilterTypeString},
	"description": {TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"latest": schema.BoolAttribute{
			Description: "If true, only the latest image will be returned.",
			Optional:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"images": schema.ListNestedBlock{
			Description: "The returned list of Images.",
			NestedObject: schema.NestedBlockObject{
				Attributes: image.ImageAttributes,
			},
		},
	},
}
