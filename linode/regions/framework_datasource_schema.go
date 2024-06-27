package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/region"
)

var filterConfig = frameworkfilter.Config{
	"site_type":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"capabilities": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"country":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"status":       {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"regions": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: region.DataSourceSchema.Attributes,
				Blocks:     region.DataSourceSchema.Blocks,
			},
		},
	},
}
