package monitorlogsstreams

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"type":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"status": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Description: "Provides a list of Linode Monitor Logs Streams.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"streams": schema.ListNestedAttribute{
			Description: "The returned list of logs streams.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The unique ID of this logs stream.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The label for this logs stream.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type of this logs stream.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of this logs stream.",
						Computed:    true,
					},
					"created_by": schema.StringAttribute{
						Description: "The user who created this logs stream.",
						Computed:    true,
					},
					"updated_by": schema.StringAttribute{
						Description: "The user who last updated this logs stream.",
						Computed:    true,
					},
					"created": schema.StringAttribute{
						Description: "When this logs stream was created.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"updated": schema.StringAttribute{
						Description: "When this logs stream was last updated.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"version": schema.Int64Attribute{
						Description: "The version of this logs stream.",
						Computed:    true,
					},
				},
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
	},
}
