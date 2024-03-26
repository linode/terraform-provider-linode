package objbucket

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"secret_key": {
		Type: schema.TypeString,
		Description: "The S3 secret key to use for this resource. (Required for lifecycle_rule and versioning). " +
			"If not specified with the resource, the value will be read from provider-level obj_secret_key, " +
			"or, generated implicitly at apply-time if obj_use_temp_keys in provider configuration is set.",
		Optional:  true,
		Sensitive: true,
	},
	"access_key": {
		Type: schema.TypeString,
		Description: "The S3 access key to use for this resource. (Required for lifecycle_rule and versioning). " +
			"If not specified with the resource, the value will be read from provider-level obj_access_key, " +
			"or, generated implicitly at apply-time if obj_use_temp_keys in provider configuration is set.",
		Optional: true,
	},
	"cluster": {
		Type:        schema.TypeString,
		Description: "The cluster of the Linode Object Storage Bucket.",
		Required:    true,
		ForceNew:    true,
	},
	"endpoint": {
		Type:        schema.TypeString,
		Description: "The endpoint for the bucket used for s3 connections.",
		Computed:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode Object Storage Bucket.",
		Required:    true,
		ForceNew:    true,
	},
	"acl": {
		Type:        schema.TypeString,
		Description: "The Access Control Level of the bucket using a canned ACL string.",
		Optional:    true,
		Default:     "private",
	},
	"cors_enabled": {
		Type:        schema.TypeBool,
		Description: "If true, the bucket will be created with CORS enabled for all origins.",
		Optional:    true,
		Default:     true,
	},
	"lifecycle_rule": {
		Type:        schema.TypeList,
		Description: "Lifecycle rules to be applied to the bucket.",
		Optional:    true,
		Elem:        resourceLifeCycle(),
	},
	"hostname": {
		Type: schema.TypeString,
		Description: "The hostname where this bucket can be accessed. " +
			"This hostname can be accessed through a browser if the bucket is made public.",
		Computed: true,
	},
	"versioning": {
		Type:        schema.TypeBool,
		Description: "Whether to enable versioning.",
		Optional:    true,
		Computed:    true,
	},
	"cert": {
		Type:        schema.TypeList,
		Description: "The cert used by this Object Storage Bucket.",
		MaxItems:    1,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"certificate": {
					Type:        schema.TypeString,
					Description: "The Base64 encoded and PEM formatted SSL certificate.",
					Sensitive:   true,
					Required:    true,
				},
				"private_key": {
					Type:        schema.TypeString,
					Description: "The private key associated with the TLS/SSL certificate.",
					Sensitive:   true,
					Required:    true,
				},
			},
		},
	},
}

var resourceSchemaLifeCycle = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The unique identifier for the rule.",
		Optional:    true,
		Computed:    true,
	},
	"prefix": {
		Type:        schema.TypeString,
		Description: "The object key prefix identifying one or more objects to which the rule applies.",
		Optional:    true,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Description: "Specifies whether the lifecycle rule is active.",
		Required:    true,
	},
	"abort_incomplete_multipart_upload_days": {
		Type: schema.TypeInt,
		Description: "Specifies the number of days after initiating a multipart upload when the multipart " +
			"upload must be completed.",
		Optional: true,
	},
	"expiration": {
		Type:        schema.TypeList,
		Description: "Specifies a period in the object's expire.",
		Optional:    true,
		MaxItems:    1,
		Elem:        resourceLifecycleExpiration(),
	},
	"noncurrent_version_expiration": {
		Type:        schema.TypeList,
		Description: "Specifies when non-current object versions expire.",
		Optional:    true,
		MaxItems:    1,
		Elem:        resourceLifecycleNoncurrentExp(),
	},
}

var resourceSchemaExpiration = map[string]*schema.Schema{
	"date": {
		Type:        schema.TypeString,
		Description: "Specifies the date after which you want the corresponding action to take effect.",
		Optional:    true,
	},
	"days": {
		Type:        schema.TypeInt,
		Description: "Specifies the number of days after object creation when the specific rule action takes effect.",
		Optional:    true,
	},
	"expired_object_delete_marker": {
		Type:        schema.TypeBool,
		Description: "Directs Linode Object Storage to remove expired deleted markers.",
		Optional:    true,
	},
}

var resourceSchemaNonCurrentExp = map[string]*schema.Schema{
	"days": {
		Type:        schema.TypeInt,
		Description: "Specifies the number of days non-current object versions expire.",
		Required:    true,
	},
}
