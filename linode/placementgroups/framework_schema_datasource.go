package placementgroups

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroup"
)

var filterConfig = frameworkfilter.Config{
	"id":                     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":                  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"region":                 {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"placement_group_type":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"is_compliant":           {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"placement_group_policy": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
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
		"placement_groups": schema.ListNestedBlock{
			Description: "The returned list of Placement Groups.",
			NestedObject: schema.NestedBlockObject{
				Attributes: placementgroup.DataSourceSchema.Attributes,
				Blocks:     placementgroup.DataSourceSchema.Blocks,
			},
		},
	},
}
