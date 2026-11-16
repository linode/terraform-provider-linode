package databasevalkeyv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{Required: true, Description: "The id of the Valkey Database."},
		"allow_list": schema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
		},
		"ca_cert": schema.StringAttribute{
			Computed:    true,
			Sensitive:   true,
			Description: "The base64-encoded SSL CA certificate for the Managed Database.",
		},
		"cluster_size": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of Linode instance nodes deployed to the Managed Database.",
		},
		"created": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "When this Managed Database was created.",
		},
		"encrypted": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the Managed Database is encrypted.",
		},
		"engine": schema.StringAttribute{
			Computed:    true,
			Description: "The Managed Database engine (e.g. valkey).",
		},
		"engine_id": schema.StringAttribute{
			Computed:    true,
			Description: "The unique ID of the database engine and version to use. (e.g. valkey/8.1)",
		},
		"fork_restore_time": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "The database timestamp from which it was restored.",
		},
		"fork_source": schema.Int64Attribute{
			Computed:    true,
			Description: "The ID of the database that was forked from.",
		},
		"host_primary": schema.StringAttribute{
			Computed:    true,
			Description: "The primary host for the Managed Database.",
		},
		"host_secondary": schema.StringAttribute{
			Computed:           true,
			DeprecationMessage: "Use host_standby instead.",
			Description:        "The secondary host for the Managed Database.",
		},
		"host_standby": schema.StringAttribute{
			Computed:    true,
			Description: "The standby host for the Managed Database.",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "A unique, user-defined string referring to the Managed Database.",
		},
		"members": schema.MapAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "A mapping between IP addresses and strings designating them as primary or failover.",
		},
		"platform": schema.StringAttribute{
			Computed:    true,
			Description: "The back-end platform used by the service.",
		},
		"port": schema.Int64Attribute{
			Computed:    true,
			Description: "The access port for this Managed Database.",
		},
		"private_network": databaseshared.DataSourceAttributePrivateNetwork,
		"region": schema.StringAttribute{
			Computed:    true,
			Description: "The Region ID for the Managed Database.",
		},
		"root_password": schema.StringAttribute{
			Computed:    true,
			Sensitive:   true,
			Description: "The randomly generated root password for the Managed Database instance.",
		},
		"root_username": schema.StringAttribute{
			Computed:    true,
			Description: "The Valkey username (akmadmin) for the Managed Database instance.",
		},
		"ssl_connection": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The operating status of the Managed Database.",
		},
		"suspended": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this database is suspended.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "The Linode Instance type used by the Managed Database for its nodes.",
		},
		"updated": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "When this Managed Database was last updated.",
		},
		"version": schema.StringAttribute{
			Computed:    true,
			Description: "The Managed Database engine version.",
		},
		"available_restore_times": schema.ListAttribute{
			Computed:    true,
			ElementType: timetypes.RFC3339Type{},
			Description: "Available discrete snapshot timestamps that can be used to fork this database. " +
				"These are not a continuous point-in-time recovery range.",
		},
		"updates":                   databaseshared.DataSourceAttributeUpdates,
		"pending_updates":           databaseshared.DataSourceAttributePendingUpdates,
		"engine_config_backup_hour": schema.Int64Attribute{Computed: true, Description: "The UTC hour when a backup starts."},
		"engine_config_backup_minute": schema.Int64Attribute{
			Computed:    true,
			Description: "The minute within the hour when a backup starts.",
		},
		"engine_config_frequent_snapshots":          schema.BoolAttribute{Computed: true, Description: "Whether frequent local RDB snapshots are enabled."},
		"engine_config_valkey_acl_channels_default": schema.StringAttribute{Computed: true, Description: "The default pub/sub channel ACL for new users."},
		"engine_config_valkey_active_expire_effort": schema.Int64Attribute{Computed: true, Description: "The effort used to reclaim expired keys."},
		"engine_config_valkey_activedefrag":         schema.BoolAttribute{Computed: true, Description: "Whether active memory defragmentation is enabled."},
		"engine_config_valkey_lfu_decay_time":       schema.Int64Attribute{Computed: true, Description: "The LFU counter decay time in minutes."},
		"engine_config_valkey_lfu_log_factor":       schema.Int64Attribute{Computed: true, Description: "The logarithm factor for LFU counters."},
		"engine_config_valkey_maxmemory_policy":     schema.StringAttribute{Computed: true, Description: "The policy used when memory is full."},
		"engine_config_valkey_number_of_databases":  schema.Int64Attribute{Computed: true, Description: "The number of logical Valkey databases."},
		"engine_config_valkey_persistence":          schema.StringAttribute{Computed: true, Description: "The persistence mode."},
		"engine_config_valkey_pubsub_client_output_buffer_limit": schema.Int64Attribute{
			Computed:    true,
			Description: "The hard output-buffer limit for pub/sub clients in MB.",
		},
		"engine_config_valkey_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The idle connection timeout in seconds.",
		},
	},
}
