package lkeversions

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"versions": {
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: elem},
		Description: "The Kubernetes version numbers available for deployment to a Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version.",
		Computed:    true,
	},
}

var elem = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The Kubernetes version.",
		Computed:    true,
	},
}
