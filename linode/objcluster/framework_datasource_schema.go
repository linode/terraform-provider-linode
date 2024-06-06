package objcluster

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	DeprecationMessage: "'linode_object_storage_cluster' data source " +
		"has been deprecated because it uses a deprecated View Cluster " +
		"API. In the future, 'region' will be preferred over 'cluster' " +
		"for the request body of many other Object Storage related API " +
		"endpoints because there may be multiple clusters per region.",
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
