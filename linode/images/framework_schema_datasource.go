package images

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/image"
)

var filterConfig = frameworkfilter.Config{
	"deprecated":  {APIFilterable: true},
	"is_public":   {APIFilterable: true},
	"label":       {APIFilterable: true},
	"size":        {APIFilterable: true},
	"type":        {APIFilterable: true},
	"vendor":      {APIFilterable: true},
	"created_by":  {APIFilterable: false},
	"id":          {APIFilterable: false},
	"status":      {APIFilterable: false},
	"description": {APIFilterable: false},
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
			NestedObject: schema.NestedBlockObject{
				Attributes: image.ImageAttributes,
			},
		},
	},
}
