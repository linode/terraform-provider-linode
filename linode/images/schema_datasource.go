package images

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"
)

var filterConfig = helper.FilterConfig{
	"deprecated": helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeBool},
	"is_public":  helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeBool},
	"label":      helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"size":       helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"type":       helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"vendor":     helper.FilterAttribute{APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"created_by": helper.FilterAttribute{TypeFunc: helper.FilterTypeString},
	"id":         helper.FilterAttribute{TypeFunc: helper.FilterTypeString},
	"status": helper.FilterAttribute{
		TypeFunc: func(value string) (interface{}, error) {
			return linodego.ImageStatus(value), nil
		},
	},
	"description": helper.FilterAttribute{TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest image will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"images": {
		Type:        schema.TypeList,
		Description: "The returned list of Images.",
		Computed:    true,
		Elem:        image.DataSource(),
	},
}
