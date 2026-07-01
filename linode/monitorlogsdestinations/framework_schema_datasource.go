package monitorlogsdestinations

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":     {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
	"label":  {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"type":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"status": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
	Description: "Provides a list of Linode Monitor Logs Destinations.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
		"destinations": schema.ListNestedAttribute{
			Description: "The returned list of logs destinations.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The unique ID of this logs destination.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The label for this logs destination.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type of this logs destination.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of this logs destination.",
						Computed:    true,
					},
					"created_by": schema.StringAttribute{
						Description: "The user who created this logs destination.",
						Computed:    true,
					},
					"updated_by": schema.StringAttribute{
						Description: "The user who last updated this logs destination.",
						Computed:    true,
					},
					"created": schema.StringAttribute{
						Description: "When this logs destination was created.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"updated": schema.StringAttribute{
						Description: "When this logs destination was last updated.",
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"version": schema.Int64Attribute{
						Description: "The version of this logs destination.",
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
