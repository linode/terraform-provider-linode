package monitorlogsdestination

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var DataSourceSchema = schema.Schema{
	Description: "Provides details about a Linode Monitor Logs Destination.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of this logs destination.",
			Required:    true,
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
		"details": schema.SingleNestedAttribute{
			Description: "The details returned by the API. Write-only fields are not included.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"access_key_id": schema.StringAttribute{
					Description: "The access key ID (akamai_object_storage type).",
					Computed:    true,
				},
				"bucket_name": schema.StringAttribute{
					Description: "The bucket name (akamai_object_storage type).",
					Computed:    true,
				},
				"host": schema.StringAttribute{
					Description: "The storage endpoint hostname (akamai_object_storage type).",
					Computed:    true,
				},
				"path": schema.StringAttribute{
					Description: "The path within the bucket (akamai_object_storage type).",
					Computed:    true,
				},
				"endpoint_url": schema.StringAttribute{
					Description: "The HTTPS endpoint URL (custom_https type).",
					Computed:    true,
				},
				"content_type": schema.StringAttribute{
					Description: "The content type of log data (custom_https type).",
					Computed:    true,
				},
				"data_compression": schema.StringAttribute{
					Description: "The compression format (custom_https type).",
					Computed:    true,
				},
				"authentication_type": schema.StringAttribute{
					Description: "The authentication type (custom_https type).",
					Computed:    true,
				},
				"tls_hostname": schema.StringAttribute{
					Description: "The TLS hostname (custom_https type).",
					Computed:    true,
				},
			},
		},
	},
}
