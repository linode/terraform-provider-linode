package nfsspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique numeric ID of the NFS Space.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The customer-provided label of the NFS Space. Must be unique within the " +
				"account. 3-63 characters, lowercase alphanumeric and hyphens only.",
			Required: true,
		},
		"description": schema.StringAttribute{
			Description: "An optional free-text description of this NFS Space.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			Description: "The current lifecycle status of the NFS Space.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"max_capacity_bytes": schema.Int64Attribute{
			Description: "The maximum total capacity available to this NFS Space across all of its " +
				"filesystems, in bytes. Akamai-managed and read-only.",
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"used_capacity_bytes": schema.Int64Attribute{
			Description: "The current used capacity across the NFS Space, in bytes. Akamai-managed " +
				"and read-only.",
			Computed: true,
		},
		"created": schema.StringAttribute{
			Description: "When this NFS Space was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this NFS Space was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"tags": schema.SetAttribute{
			Description: "Customer-defined metadata tags for grouping and filtering.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
