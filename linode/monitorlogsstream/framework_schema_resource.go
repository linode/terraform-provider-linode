package monitorlogsstream

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the logs stream.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label of the logs stream.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the logs stream (audit_logs or lke_audit_logs).",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("audit_logs", "lke_audit_logs"),
			},
		},
		"status": schema.StringAttribute{
			Description: "The status of the logs stream (active, inactive, provisioning, deactivating).",
			Optional:    true,
			Computed:    true,
		},
		"version": schema.Int64Attribute{
			Description: "The version of the logs stream.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"destinations": schema.ListAttribute{
			Description: "The list of logs destination IDs attached to this stream.",
			Required:    true,
			ElementType: types.Int64Type,
		},
		"details": schema.SingleNestedAttribute{
			Description: "Additional configuration details. Only applies to lke_audit_logs streams.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"cluster_ids": schema.ListAttribute{
					Description: "The list of LKE cluster IDs to include in this stream.",
					Optional:    true,
					Computed:    true,
					ElementType: types.Int64Type,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				"is_auto_add_all_clusters_enabled": schema.BoolAttribute{
					Description: "When true, all LKE clusters are automatically added to this stream.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"created": schema.StringAttribute{
			Description: "The date and time when the logs stream was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "The date and time when the logs stream was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"created_by": schema.StringAttribute{
			Description: "The user who created the logs stream.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_by": schema.StringAttribute{
			Description: "The user who last updated the logs stream.",
			Computed:    true,
		},
	},
}

// ResourceSchema is the exported schema for use in tests and external packages.
var ResourceSchema = frameworkResourceSchema
