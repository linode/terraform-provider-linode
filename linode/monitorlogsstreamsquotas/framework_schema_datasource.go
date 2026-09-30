package monitorlogsstreamsquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDataSourceSchema = schema.Schema{
	Description: "Provides details about logs streams quotas.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"quotas": schema.ListNestedAttribute{
			Description: "The returned list of logs streams quotas.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"quota_id": schema.StringAttribute{
						Description: "The ID of the logs streams quota.",
						Computed:    true,
					},
					"quota_name": schema.StringAttribute{
						Description: "The name of the logs streams quota.",
						Computed:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description of the logs streams quota.",
						Computed:    true,
					},
					"quota_limit": schema.Int64Attribute{
						Description: "The maximum number allowed by this quota.",
						Computed:    true,
					},
					"quota_type": schema.StringAttribute{
						Description: "The type of the logs streams quota.",
						Computed:    true,
					},
				},
			},
		},
	},
}
