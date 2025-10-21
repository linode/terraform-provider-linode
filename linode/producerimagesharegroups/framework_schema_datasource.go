package producerimagesharegroups

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroup"
)

var filterConfig = frameworkfilter.Config{
	"id":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":        {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"is_suspended": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
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
		"image_share_groups": schema.ListNestedBlock{
			Description: "The returned list of Image SHare Groups.",
			NestedObject: schema.NestedBlockObject{
				Attributes: producerimagesharegroup.Attributes,
			},
		},
	},
}
