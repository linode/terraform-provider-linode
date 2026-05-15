package monitorlogsdestination

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var frameworkResourceSchema = schema.Schema{
	Description: "Manages a Linode Monitor Logs Destination.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this logs destination.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The label for this logs destination.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of this logs destination. One of: akamai_object_storage, custom_https.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("akamai_object_storage", "custom_https"),
			},
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
		"akamai_object_storage_details": schema.SingleNestedAttribute{
			Description: "Details for an akamai_object_storage logs destination.",
			Optional:    true,
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRoot("custom_https_details"),
				),
			},
			Attributes: map[string]schema.Attribute{
				"access_key_id": schema.StringAttribute{
					Description: "The access key ID for the object storage bucket.",
					Required:    true,
				},
				"access_key_secret": schema.StringAttribute{
					Description: "The access key secret for the object storage bucket. " +
						"This value is write-only and will not be returned by the API.",
					Required:  true,
					Sensitive: true,
				},
				"bucket_name": schema.StringAttribute{
					Description: "The name of the object storage bucket.",
					Required:    true,
				},
				"host": schema.StringAttribute{
					Description: "The hostname of the object storage endpoint.",
					Required:    true,
				},
				"path": schema.StringAttribute{
					Description: "The path within the bucket where logs will be stored.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
		"custom_https_details": schema.SingleNestedAttribute{
			Description: "Details for a custom_https logs destination.",
			Optional:    true,
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRoot("akamai_object_storage_details"),
				),
			},
			Attributes: map[string]schema.Attribute{
				"endpoint_url": schema.StringAttribute{
					Description: "The HTTPS endpoint URL to send logs to.",
					Required:    true,
				},
				"content_type": schema.StringAttribute{
					Description: "The content type of the log data. One of: application/json, text/plain.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("application/json", "text/plain"),
					},
				},
				"data_compression": schema.StringAttribute{
					Description: "The compression format for log data. One of: none, gzip.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("none", "gzip"),
					},
				},
				"authentication": schema.SingleNestedAttribute{
					Description: "Authentication configuration for the HTTPS endpoint.",
					Required:    true,
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "The authentication type. One of: basic, none.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("basic", "none"),
							},
						},
						"username": schema.StringAttribute{
							Description: "The username for basic authentication. " +
								"This value is write-only and will not be returned by the API.",
							Optional:  true,
							Sensitive: true,
						},
						"password": schema.StringAttribute{
							Description: "The password for basic authentication. " +
								"This value is write-only and will not be returned by the API.",
							Optional:  true,
							Sensitive: true,
						},
					},
				},
				"client_certificate_details": schema.SingleNestedAttribute{
					Description: "TLS client certificate configuration.",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"tls_hostname": schema.StringAttribute{
							Description: "The TLS hostname for certificate verification.",
							Required:    true,
						},
						"client_ca_certificate": schema.StringAttribute{
							Description: "The client CA certificate. " +
								"This value is write-only and will not be returned by the API.",
							Required:  true,
							Sensitive: true,
						},
						"client_certificate": schema.StringAttribute{
							Description: "The client certificate. " +
								"This value is write-only and will not be returned by the API.",
							Required:  true,
							Sensitive: true,
						},
						"client_private_key": schema.StringAttribute{
							Description: "The client private key. " +
								"This value is write-only and will not be returned by the API.",
							Required:  true,
							Sensitive: true,
						},
					},
				},
				"custom_headers": schema.ListNestedAttribute{
					Description: "Custom HTTP headers to include in log delivery requests.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "The name of the HTTP header.",
								Required:    true,
							},
							"value": schema.StringAttribute{
								Description: "The value of the HTTP header. " +
									"This value is write-only and will not be returned by the API.",
								Required:  true,
								Sensitive: true,
							},
						},
					},
				},
			},
		},
	},
}
