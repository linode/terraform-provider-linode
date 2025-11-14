package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/region"
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
		"regions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: region.DataSourceSchema.Attributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
