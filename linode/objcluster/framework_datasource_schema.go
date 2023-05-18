package objcluster

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this Cluster.",
			Required:    true,
		},
		"domain": schema.StringAttribute{
			Description: "The base URL for this cluster.",
			Optional:    true,
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "This cluster's status.",
			Optional:    true,
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region this cluster is located in.",
			Optional:    true,
			Computed:    true,
		},
		"static_site_domain": schema.StringAttribute{
			Description: "The base URL for this cluster used when hosting static sites.",
			Optional:    true,
			Computed:    true,
		},
	},
}
