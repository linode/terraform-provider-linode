package domains

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/domain"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"tags":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"domain": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},

	"type":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"description": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"master_ips":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"axfr_ips":    {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"ttl_sec":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"retry_sec":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"expire_sec":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"refresh_sec": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"soa_email":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"domains": schema.ListNestedAttribute{
			Description: "The returned list of Domains.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: domain.DomainAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
