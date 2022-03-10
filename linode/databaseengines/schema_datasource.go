package databaseengines

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"engine":  {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"version": {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"id": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
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
	"engines": {
		Type:        schema.TypeList,
		Description: "The returned list of engines.",
		Computed:    true,
		Elem:        engineSchema(),
	},
}

func engineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:        schema.TypeString,
				Description: "The Managed Database engine type.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The Managed Database engine ID in engine/version format.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The Managed Database engine version.",
				Computed:    true,
			},
		},
	}
}
