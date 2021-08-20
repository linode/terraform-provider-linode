package objectkey

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The label given to this key. For display purposes only.",
		Required:    true,
	},
	"access_key": {
		Type:        schema.TypeString,
		Description: "This keypair's access key. This is not secret.",
		Computed:    true,
	},
	"secret_key": {
		Type:        schema.TypeString,
		Description: "This keypair's secret key.",
		Sensitive:   true,
		Computed:    true,
	},
	"limited": {
		Type:        schema.TypeBool,
		Description: "Whether or not this key is a limited access key.",
		Computed:    true,
	},
	"bucket_access": {
		Type:        schema.TypeList,
		Description: "A list of permissions to grant this limited access key.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: resourceAccessSchema,
		},
		ForceNew: true,
	},
}

var resourceAccessSchema = map[string]*schema.Schema{
	"bucket_name": {
		Type:        schema.TypeString,
		Description: "The unique label of the bucket to which the key will grant limited access.",
		Required:    true,
	},
	"cluster": {
		Type:        schema.TypeString,
		Description: "The Object Storage cluster where a bucket to which the key is granting access is hosted.",
		Required:    true,
	},
	"permissions": {
		Type:        schema.TypeString,
		Description: "This Limited Access Key’s permissions for the selected bucket.",
		Required:    true,
	},
}
