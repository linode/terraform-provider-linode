package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

var filterConfig = helper.FilterConfig{
	"class":       {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"disk":        {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"gpus":        {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"label":       {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"memory":      {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"network_out": {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"transfer":    {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"vcpus":       {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
}

var dataSourceSchema = map[string]*schema.Schema{
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"types": {
		Type:        schema.TypeList,
		Description: "The returned list of Types.",
		Computed:    true,
		Elem:        instancetype.DataSource(),
	},
}
