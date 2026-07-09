package monitorlogsstream

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the logs stream.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the logs stream.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the logs stream.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the logs stream.",
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of the logs stream.",
			Computed:    true,
		},
		"destinations": schema.ListAttribute{
			Description: "The list of logs destination IDs attached to this stream.",
			Computed:    true,
			ElementType: types.Int64Type,
		},
		"details": schema.SingleNestedAttribute{
			Description: "Additional configuration details. Only applies to lke_audit_logs streams.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"cluster_ids": schema.ListAttribute{
					Description: "The list of LKE cluster IDs included in this stream.",
					Computed:    true,
					ElementType: types.Int64Type,
				},
				"is_auto_add_all_clusters_enabled": schema.BoolAttribute{
					Description: "When true, all LKE clusters are automatically added to this stream.",
					Computed:    true,
				},
			},
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the logs stream was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the logs stream was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"created_by": schema.StringAttribute{
			Description: "The user who created the logs stream.",
			Computed:    true,
		},
		"updated_by": schema.StringAttribute{
			Description: "The user who last updated the logs stream.",
			Computed:    true,
		},
	},
}
