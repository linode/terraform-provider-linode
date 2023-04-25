package objbucket

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var bucketDataSourceSchema = map[string]*schema.Schema{
	"cluster": {
		Type:        schema.TypeString,
		Description: "The ID of the Object Storage Cluster this bucket is in.",
		Required:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this bucket was created.",
		Computed:    true,
	},
	"hostname": {
		Type: schema.TypeString,
		Description: "The hostname where this bucket can be accessed." +
			"This hostname can be accessed through a browser if the bucket is made public.",
		Computed: true,
	},
	"id": {
		Type:        schema.TypeString,
		Description: "The id of this bucket.",
		Computed:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The name of this bucket.",
		Required:    true,
	},
	"objects": {
		Type:        schema.TypeInt,
		Description: "The number of objects stored in this bucket.",
		Computed:    true,
	},
	"size": {
		Type:        schema.TypeInt,
		Description: "The size of the bucket in bytes.",
		Computed:    true,
	},
}
