package objquota

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FrameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Object Storage quota.",
			Required:    true,
		},
		"quota_name": schema.StringAttribute{
			Description: "The name of the Object Storage quota.",
			Computed:    true,
		},
		"endpoint_type": schema.StringAttribute{
			Description: "The type of the S3 endpoint of the Object Storage.",
			Computed:    true,
		},
		"s3_endpoint": schema.StringAttribute{
			Description: "The S3 endpoint URL of the Object Storage, based on the `endpoint_type` and `region`.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "The description of the Object Storage quota.",
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
		"quota_usage": schema.ObjectAttribute{
			Description:    "The usage data for a specific Object Storage related quota on your account.",
			Computed:       true,
			AttributeTypes: quotaUsageAttributes,
		},
	},
}

var quotaUsageAttributes = map[string]attr.Type{
	"quota_limit": types.Int64Type,
	"usage":       types.Int64Type,
}
