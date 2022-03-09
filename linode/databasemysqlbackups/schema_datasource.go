package databasemysqlbackups

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"created": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"type":    {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"id":    {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"label": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest backup will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"backups": {
		Type:        schema.TypeList,
		Description: "The returned list of backups.",
		Computed:    true,
		Elem:        DataSource(),
	},
}
