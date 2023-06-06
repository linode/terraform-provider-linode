package databaseengines

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"engine":  {APIFilterable: true},
	"version": {APIFilterable: true},
	"id":      {APIFilterable: false},
}

var engineSchema = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine type.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The Managed Database engine ID in engine/version format.",
			Computed:    true,
		},
		"version": schema.StringAttribute{
			Description: "The Managed Database engine version.",
			Computed:    true,
		},
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"latest": schema.BoolAttribute{
			Description: "If true, only the latest engine version will be returned.",
			Optional:    true,
		},
		"order_by": filterConfig.OrderBySchema(),
		"order":    filterConfig.OrderSchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"engines": schema.ListNestedBlock{
			Description:  "The returned list of engines.",
			NestedObject: engineSchema,
		},
	},
}
