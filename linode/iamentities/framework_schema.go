package iamentities

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"label": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"type":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var frameworkSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The datasource's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"entities": schema.ListNestedAttribute{
			Description: "The user entity level access.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the entity.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The entity label.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The entity category.",
						Computed:    true,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
