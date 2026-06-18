package obj

import (
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the object.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"bucket": schema.StringAttribute{
			Description: "The target bucket to put this object in.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"region": schema.StringAttribute{
			Description: "The target region that the bucket is in.",
			Optional:    true,
		},
		"key": schema.StringAttribute{
			Description: "The name of the uploaded object.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"secret_key": schema.StringAttribute{
			Description: "The REQUIRED S3 secret key with access to the target bucket. " +
				"If not specified with the resource, you must provide its value by configuring the obj_secret_key, " +
				"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
			Optional:  true,
			Sensitive: true,
		},
		"access_key": schema.StringAttribute{
			Description: "The REQUIRED S3 access key with access to the target bucket. " +
				"If not specified with the resource, you must provide its value by configuring the obj_access_key, " +
				"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
			Optional: true,
		},
		"content": schema.StringAttribute{
			Description: "The contents of the Object to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content_base64"),
					path.MatchRoot("source"),
				),
			},
		},
		"content_base64": schema.StringAttribute{
			Description: "The base64 contents of the Object to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content"),
					path.MatchRoot("source"),
				),
			},
		},
		"source": schema.StringAttribute{
			Description: "The source file to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content"),
					path.MatchRoot("content_base64"),
				),
			},
		},
		"acl": schema.StringAttribute{
			Description: "The ACL config given to this object.",
			Optional:    true,
			Computed:    true,
			Default: stringdefault.StaticString(
				string(s3types.ObjectCannedACLPrivate),
			),
			Validators: []validator.String{
				stringvalidator.OneOf(
					helper.StringAliasSliceToStringSlice(
						s3types.ObjectCannedACLPrivate.Values(),
					)...,
				),
			},
		},
		"cache_control": schema.StringAttribute{
			Description: "This cache_control configuration of this object.",
			Optional:    true,
		},
		"content_disposition": schema.StringAttribute{
			Description: "The content disposition configuration of this object.",
			Optional:    true,
		},
		"content_encoding": schema.StringAttribute{
			Description: "The encoding of the content of this object.",
			Optional:    true,
		},
		"content_language": schema.StringAttribute{
			Description: "The language metadata of this object.",
			Optional:    true,
		},
		"content_type": schema.StringAttribute{
			Description: "The MIME type of the content.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"endpoint": schema.StringAttribute{
			Description: "The endpoint for the bucket used for s3 connections.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"etag": schema.StringAttribute{
			Description: "The specific version of this object.",
			Optional:    true,
			Computed:    true,
		},
		"force_destroy": schema.BoolAttribute{
			Description: "Whether the object should bypass deletion restrictions.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"metadata": schema.MapAttribute{
			Description: "The metadata of this object",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			Default:     helper.EmptyMapDefault(types.StringType),
		},
		"version_id": schema.StringAttribute{
			Description: "The version ID of this object.",
			Computed:    true,
		},
		"website_redirect": schema.StringAttribute{
			Description: "The website redirect location of this object.",
			Optional:    true,
		},
	},
}
