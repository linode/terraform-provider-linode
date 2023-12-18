package nbconfigs

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig"
)

var filterConfig = frameworkfilter.Config{
	"algorithm":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"check":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"nodebalancer_id": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"protocol":        {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"proxy_protocol":  {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"stickiness":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
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
		"nodebalancer_configs": schema.ListNestedBlock{
			Description: "The returned list of NodeBalancer Configs.",
			NestedObject: schema.NestedBlockObject{
				Attributes: nbconfig.NBConfigAttributes,
			},
		},
	},
}
