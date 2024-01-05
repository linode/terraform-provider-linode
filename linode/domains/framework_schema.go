package domains

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/domain"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"group": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"tags":  {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"domain":      {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"type":        {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"status":      {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"description": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"master_ips":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"axfr_ips":    {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"ttl_sec":     {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"retry_sec":   {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"expire_sec":  {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"refresh_sec": {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"soa_email":   {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
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
		"domains": schema.ListNestedBlock{
			Description: "The returned list of Domains.",
			NestedObject: schema.NestedBlockObject{
				Attributes: domain.DomainAttributes,
			},
		},
	},
}
