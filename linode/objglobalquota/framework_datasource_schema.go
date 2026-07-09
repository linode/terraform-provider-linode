package objglobalquota

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"quota_id": schema.StringAttribute{
			Description: "The ID of the Object Storage global quota.",
			Required:    true,
		},
		"quota_name": schema.StringAttribute{
			Description: "The name of the Object Storage global quota.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the Object Storage global quota.",
			Computed:    true,
		},
		"quota_limit": schema.Int64Attribute{
			Description: "The maximum quantity of the `resource_metric` allowed by the quota.",
			Computed:    true,
		},
		"resource_metric": schema.StringAttribute{
			Description: "The specific Object Storage resource for the quota.",
			Computed:    true,
		},
		"quota_type": schema.StringAttribute{
			Description: "The type of the Object Storage global quota.",
			Computed:    true,
		},
		"has_usage": schema.BoolAttribute{
			Description: "Whether usage data is available for this Object Storage global quota.",
			Computed:    true,
		},
		"quota_usage": schema.ObjectAttribute{
			Description:    "The usage data for a specific global Object Storage quota on your account.",
			Computed:       true,
			AttributeTypes: quotaUsageAttributes,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID of the Object Storage global quota data source.",
			Computed:    true,
		},
	},
}

var quotaUsageAttributes = map[string]attr.Type{
	"quota_limit": types.Int64Type,
	"usage":       types.Int64Type,
}
