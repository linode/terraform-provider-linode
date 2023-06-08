package stackscripts

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/stackscript"
)

var filterConfig = frameworkfilter.Config{
	"deployments_total": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeInt},
	"description":       {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"is_public":         {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"label":             {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},

	"rev_note":           {TypeFunc: frameworkfilter.FilterTypeString},
	"mine":               {TypeFunc: frameworkfilter.FilterTypeBool},
	"deployments_active": {TypeFunc: frameworkfilter.FilterTypeInt},
	"images":             {TypeFunc: frameworkfilter.FilterTypeString},
	"username":           {TypeFunc: frameworkfilter.FilterTypeString},
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
