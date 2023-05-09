package databasepostgresql

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dataSourceSchema = map[string]*schema.Schema{
	"database_id": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The ID of the PostgreSQL database.",
	},

	"engine_id": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine in engine/version format. (e.g. postgresql/12.6)",
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
	"replication_commit_type": {
		Type:        schema.TypeString,
		Description: "The synchronization level of the replicating server.",
		Computed:    true,
	},
	"ssl_connection": {
		Type:        schema.TypeBool,
		Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
		Computed:    true,
	},
	"updates": {
		Type:        schema.TypeList,
		Description: "Configuration settings for automated patch update maintenance for the Managed Database.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"day_of_week": {
					Type:        schema.TypeString,
					Description: "The day to perform maintenance.",
					Computed:    true,
				},
				"duration": {
					Type:        schema.TypeInt,
					Description: "The maximum maintenance window time in hours.",
					Computed:    true,
				},
				"frequency": {
					Type:        schema.TypeString,
					Description: "Whether maintenance occurs on a weekly or monthly basis.",
					Computed:    true,
				},
				"hour_of_day": {
					Type:        schema.TypeInt,
					Description: "The hour to begin maintenance based in UTC time.",
					Computed:    true,
				},
				"week_of_month": {
					Type: schema.TypeInt,
					Description: "The week of the month to perform monthly frequency updates." +
						" Required for monthly frequency updates.",
					Computed: true,
				},
			},
		},
	},

	// Computed fields
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
	"port": {
		Type:        schema.TypeInt,
		Description: "The access port for this Managed Database.",
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
