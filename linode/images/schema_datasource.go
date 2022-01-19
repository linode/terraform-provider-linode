package images

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"
)

var filterConfig = helper.FilterConfig{
	"deprecated": {APIFilterable: true, TypeFunc: helper.FilterTypeBool},
	"is_public":  {APIFilterable: true, TypeFunc: helper.FilterTypeBool},
	"label":      {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"size":       {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"type":       {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"vendor":     {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"created_by": {TypeFunc: helper.FilterTypeString},
	"id":         {TypeFunc: helper.FilterTypeString},
	"status": {
		TypeFunc: func(value string) (interface{}, error) {
			return linodego.ImageStatus(value), nil
		},
	},
	"description": {TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest image will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": helper.OrderBySchema(filterConfig),
	"order":    helper.OrderSchema(),
	"filter":   helper.FilterSchema(filterConfig),
	"images": {
		Type:        schema.TypeList,
		Description: "The returned list of Images.",
		Computed:    true,
		Elem:        image.DataSource(),
	},
}
