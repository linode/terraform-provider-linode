package databasemysql

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dataSourceSchema = map[string]*schema.Schema{
	"database_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The ID of the MySQL database.",
	},

	"engine_id": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine in engine/version format. (e.g. mysql/8.0.26)",
		Computed:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "A unique, user-defined string referring to the Managed Database.",
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region to use for the Managed Database.",
		Computed:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The Linode Instance type used by the Managed Database for its nodes.",
		Computed:    true,
	},
	"allow_list": {
		Type: schema.TypeSet,
		Description: "A list of IP addresses that can access the Managed Database. " +
			"Each item can be a single IP address or a range in CIDR format.",
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"cluster_size": {
		Type:        schema.TypeInt,
		Description: "The number of Linode Instance nodes deployed to the Managed Database. Defaults to 1.",
		Computed:    true,
	},
	"encrypted": {
		Type:        schema.TypeBool,
		Description: "Whether the Managed Databases is encrypted.",
		Computed:    true,
	},
	"replication_type": {
		Type:        schema.TypeString,
		Description: "The replication method used for the Managed Database.",
		Computed:    true,
	},
	"ssl_connection": {
		Type:        schema.TypeBool,
		Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
		Computed:    true,
	},
	"ca_cert": {
		Type:        schema.TypeString,
		Description: "The base64-encoded SSL CA certificate for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this Managed Database was created.",
		Computed:    true,
	},
	"engine": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine.",
		Computed:    true,
	},
	"host_primary": {
		Type:        schema.TypeString,
		Description: "The primary host for the Managed Database.",
		Computed:    true,
	},
	"host_secondary": {
		Type:        schema.TypeString,
		Description: "The secondary host for the Managed Database.",
		Computed:    true,
	},
	"root_password": {
		Type:        schema.TypeString,
		Description: "The randomly-generated root password for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The operating status of the Managed Database.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this Managed Database was last updated.",
		Computed:    true,
	},
	"root_username": {
		Type:        schema.TypeString,
		Description: "The root username for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"version": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine version.",
		Computed:    true,
	},
}
