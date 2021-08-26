package template

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode Template.",
		Optional:    true,
	},
	"status": {
		Type:        schema.TypeInt,
		Description: "The status of the template, indicating the current readiness state.",
		Computed:    true,
	},
}
