package databases

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"engine":  {TypeFunc: helper.FilterTypeString, APIFilterable: true},
	"label":   {TypeFunc: helper.FilterTypeString, APIFilterable: true},
	"region":  {TypeFunc: helper.FilterTypeString, APIFilterable: true},
	"status":  {TypeFunc: helper.FilterTypeString, APIFilterable: true},
	"type":    {TypeFunc: helper.FilterTypeString, APIFilterable: true},
	"version": {TypeFunc: helper.FilterTypeString, APIFilterable: true},

	"allow_list":     {TypeFunc: helper.FilterTypeString},
	"cluster_size":   {TypeFunc: helper.FilterTypeInt},
	"created":        {TypeFunc: helper.FilterTypeString},
	"encrypted":      {TypeFunc: helper.FilterTypeBool},
	"host_primary":   {TypeFunc: helper.FilterTypeString},
	"host_secondary": {TypeFunc: helper.FilterTypeString},
	"id":             {TypeFunc: helper.FilterTypeInt},
	"instance_uri":   {TypeFunc: helper.FilterTypeString},
	"updated":        {TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest engine will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"databases": {
		Type:        schema.TypeList,
		Description: "The returned list of databases.",
		Computed:    true,
		Elem:        DataSource(),
	},
}
