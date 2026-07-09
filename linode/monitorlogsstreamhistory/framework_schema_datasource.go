package monitorlogsstreamhistory

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var streamHistoryAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the logs stream version.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the logs stream at this version.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of the logs stream (audit_logs or lke_audit_logs).",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The status of the logs stream at this version.",
		Computed:    true,
	},
	"version": schema.Int64Attribute{
		Description: "The version number of this history entry.",
		Computed:    true,
	},
	"destinations": schema.ListAttribute{
		Description: "The destination IDs configured at this version.",
		Computed:    true,
		ElementType: types.Int64Type,
	},
	"details": schema.SingleNestedAttribute{
		Description: "Additional configuration details at this version.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"cluster_ids": schema.ListAttribute{
				Description: "The LKE cluster IDs included in this stream version.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"is_auto_add_all_clusters_enabled": schema.BoolAttribute{
				Description: "Whether all clusters were auto-added at this version.",
				Computed:    true,
			},
		},
	},
	"created": schema.StringAttribute{
		Description: "The date and time when this version was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "The date and time when this version was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"created_by": schema.StringAttribute{
		Description: "The user who created this stream version.",
		Computed:    true,
	},
	"updated_by": schema.StringAttribute{
		Description: "The user who last updated this stream version.",
		Computed:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier for this data source (stream_id as string).",
			Computed:    true,
		},
		"stream_id": schema.Int64Attribute{
			Description: "The ID of the logs stream to retrieve history for.",
			Required:    true,
		},
		"streams": schema.ListNestedAttribute{
			Description:  "The historical versions of the logs stream.",
			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: streamHistoryAttributes},
		},
	},
}
