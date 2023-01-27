package lkeversion

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "A Kubernetes version number available for deployment to a Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version.",
		Required:    true,
	},
}
