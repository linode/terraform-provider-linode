package objcluster

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	DeprecationMessage: "This data source has been deprecated because it " +
		"relies on deprecated API endpoints. Going forward, `region` will " +
		"be the preferred way to designate where Object Storage resources " +
		"should be created.",

	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this Cluster.",
			Required:    true,
		},
		"domain": schema.StringAttribute{
			Description: "The base URL for this cluster.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "This cluster's status.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region this cluster is located in.",
			Computed:    true,
		},
		"static_site_domain": schema.StringAttribute{
			Description: "The base URL for this cluster used when hosting static sites.",
			Computed:    true,
		},
	},
}
