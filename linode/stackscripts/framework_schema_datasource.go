package stackscripts

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscript"
)

var filterConfig = frameworkfilter.Config{
	"deployments_total": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"description":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"is_public":         {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"label":             {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"rev_note":          {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"mine":              {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},

	"deployments_active": {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"images":             {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"username":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"latest": schema.BoolAttribute{
			Description: "If true, only the latest StackScript will be returned.",
			Optional:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"stackscripts": schema.ListNestedBlock{
			Description: "The returned list of StackScripts.",
			NestedObject: schema.NestedBlockObject{
				Attributes: stackscript.StackscriptAttributes,
			},
		},
	},
}
