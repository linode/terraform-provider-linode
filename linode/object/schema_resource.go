package object

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceSchema = map[string]*schema.Schema{
	"bucket": {
		Type:        schema.TypeString,
		Description: "The target bucket to put this object in.",
		Required:    true,
	},
	"cluster": {
		Type:        schema.TypeString,
		Description: "The target cluster that the bucket is in.",
		Required:    true,
	},
	"key": {
		Type:        schema.TypeString,
		Description: "The name of the uploaded object.",
		Required:    true,
	},
	"secret_key": {
		Type:        schema.TypeString,
		Description: "The S3 secret key with access to the target bucket.",
		Required:    true,
	},
	"access_key": {
		Type:        schema.TypeString,
		Description: "The S3 access key with access to the target bucket.",
		Required:    true,
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
		Type:        schema.TypeString,
		Description: "The ACL config given to this object.",
		Default:     s3.ObjectCannedACLPrivate,
		Optional:    true,
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
