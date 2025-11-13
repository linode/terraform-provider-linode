package producerimagesharegroupmembers

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmember"
)

var filterConfig = frameworkfilter.Config{
	"token_uuid": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"label":      {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"status":     {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"sharegroup_id": schema.Int64Attribute{
			Description: "The ID of the Image Share Group for which to list members.",
			Required:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"members": schema.ListNestedBlock{
			Description: "The returned list of Image Share Group Members.",
			NestedObject: schema.NestedBlockObject{
				Attributes: producerimagesharegroupmember.Attributes,
			},
		},
	},
}
