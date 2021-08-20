package objectcluster

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The unique ID of this Cluster.",
		Required:    true,
	},
	"domain": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The base URL for this cluster.",
		Computed:    true,
	},
	"status": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "This cluster's status.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The region this cluster is located in.",
		Computed:    true,
	},
	"static_site_domain": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The base URL for this cluster used when hosting static sites.",
		Computed:    true,
	},
}
