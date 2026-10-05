package objbucket

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of this bucket.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The name of this bucket.",
			Required:    true,
		},
		"region": schema.StringAttribute{
			Description: "The ID of the region this bucket is in.",
			Optional:    true,
			Computed:    true,
		},
		"endpoint_type": schema.StringAttribute{
			Description: "The type of the S3 endpoint of the bucket.",
			Computed:    true,
		},
		"s3_endpoint": schema.StringAttribute{
			Description: "The S3 endpoint URL of the bucket, based on the `endpoint_type` and `region`.",
			Computed:    true,
		},
		"hostname": schema.StringAttribute{
			Description: "The hostname where this bucket can be accessed." +
				"This hostname can be accessed through a browser if the bucket is made public.",
			Computed: true,
		},
		"objects": schema.Int64Attribute{
			Description: "The number of objects stored in this bucket.",
			Computed:    true,
		},
		"size": schema.Int64Attribute{
			Description: "The size of the bucket in bytes.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this bucket was created.",
			CustomType:  timetypes.RFC3339Type{},
			Computed:    true,
		},
	},
}
