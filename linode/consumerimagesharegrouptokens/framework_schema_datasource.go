package consumerimagesharegrouptokens

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegrouptoken"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"token_uuid":                {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"label":                     {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"status":                    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"valid_for_sharegroup_uuid": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"sharegroup_uuid":           {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"sharegroup_label":          {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
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
		"tokens": schema.ListNestedBlock{
			Description: "The returned list of Image Share Group Tokens.",
			NestedObject: schema.NestedBlockObject{
				Attributes: consumerimagesharegrouptoken.Attributes,
			},
		},
	},
}
