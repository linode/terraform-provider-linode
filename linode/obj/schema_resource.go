package obj

import (
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	"bucket": {
		Type:        schema.TypeString,
		Description: "The target bucket to put this object in.",
		Required:    true,
		ForceNew:    true,
	},
	"cluster": {
		Type:        schema.TypeString,
		Description: "The target cluster that the bucket is in.",
		Deprecated: "The cluster attribute has been deprecated, please consider switching to the region attribute. " +
			"For example, a cluster value of `us-mia-1` can be translated to a region value of `us-mia`.",
		Optional:     true,
		ForceNew:     true,
		ExactlyOneOf: []string{"cluster", "region"},
	},
	"region": {
		Type:         schema.TypeString,
		Description:  "The target region that the bucket is in.",
		Optional:     true,
		ForceNew:     true,
		ExactlyOneOf: []string{"cluster", "region"},
	},
	"key": {
		Type:        schema.TypeString,
		Description: "The name of the uploaded object.",
		Required:    true,
		ForceNew:    true,
	},
	"secret_key": {
		Type: schema.TypeString,
		Description: "The REQUIRED S3 secret key with access to the target bucket. " +
			"If not specified with the resource, you must provide its value by configuring the obj_secret_key, " +
			"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
		Optional:  true,
		Sensitive: true,
	},
	"access_key": {
		Type: schema.TypeString,
		Description: "The REQUIRED S3 access key with access to the target bucket. " +
			"If not specified with the resource, you must provide its value by configuring the obj_access_key, " +
			"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
		Optional: true,
	},
	"content": {
		Type:         schema.TypeString,
		Description:  "The contents of the Object to upload.",
		Optional:     true,
		ExactlyOneOf: []string{"content", "content_base64", "source"},
	},
	"content_base64": {
		Type:        schema.TypeString,
		Description: "The base64 contents of the Object to upload.",
		Optional:    true,
	},
	"source": {
		Type:        schema.TypeString,
		Description: "The source file to upload.",
		Optional:    true,
	},
	"acl": {
		Type:             schema.TypeString,
		Description:      "The ACL config given to this object.",
		Default:          s3types.ObjectCannedACLPrivate,
		ValidateDiagFunc: helper.SDKv2ObjectCannedACLValidator,
		Optional:         true,
	},
	"cache_control": {
		Type:        schema.TypeString,
		Description: "This cache_control configuration of this object.",
		Optional:    true,
	},
	"content_disposition": {
		Type:        schema.TypeString,
		Description: "The content disposition configuration of this object.",
		Optional:    true,
	},
	"content_encoding": {
		Type:        schema.TypeString,
		Description: "The encoding of the content of this object.",
		Optional:    true,
	},
	"content_language": {
		Type:        schema.TypeString,
		Description: "The language metadata of this object.",
		Optional:    true,
	},
	"content_type": {
		Type:        schema.TypeString,
		Description: "The MIME type of the content.",
		Optional:    true,
		Computed:    true,
	},
	"endpoint": {
		Type:        schema.TypeString,
		Description: "The endpoint for the bucket used for s3 connections.",
		Computed:    true,
		Optional:    true,
	},
	"etag": {
		Type:        schema.TypeString,
		Description: "The specific version of this object.",
		Optional:    true,
		Computed:    true,
	},
	"force_destroy": {
		Type:        schema.TypeBool,
		Description: "Whether the object should bypass deletion restrictions.",
		Optional:    true,
		Default:     false,
	},
	"metadata": {
		Type:        schema.TypeMap,
		Description: "The metadata of this object",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"version_id": {
		Type:        schema.TypeString,
		Description: "The version ID of this object.",
		Computed:    true,
	},
	"website_redirect": {
		Type:        schema.TypeString,
		Description: "The website redirect location of this object.",
		Optional:    true,
	},
}
