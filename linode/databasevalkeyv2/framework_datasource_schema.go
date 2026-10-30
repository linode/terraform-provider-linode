package databasevalkeyv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{Required: true, Description: "The ID of the Managed Valkey Database."},
		"allow_list": schema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "IP addresses or CIDR ranges allowed to access the database.",
		},
		"ca_cert":      schema.StringAttribute{Computed: true, Sensitive: true, Description: "The base64-encoded SSL CA certificate."},
		"cluster_size": schema.Int64Attribute{Computed: true, Description: "The number of nodes in the database cluster."},
		"created":      schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}, Description: "When the database was created."},
		"encrypted":    schema.BoolAttribute{Computed: true, Description: "Whether the database is encrypted."},
		"engine":       schema.StringAttribute{Computed: true, Description: "The database engine."},
		"engine_id":    schema.StringAttribute{Computed: true, Description: "The engine and version identifier."},
		"fork_restore_time": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "The snapshot time used when this database was forked.",
		},
		"fork_source":  schema.Int64Attribute{Computed: true, Description: "The source database ID when this database was forked."},
		"host_primary": schema.StringAttribute{Computed: true, Description: "The primary database host."},
		"host_secondary": schema.StringAttribute{
			Computed:           true,
			DeprecationMessage: "Use host_standby instead.",
			Description:        "The secondary database host.",
		},
		"host_standby":    schema.StringAttribute{Computed: true, Description: "The standby database host."},
		"label":           schema.StringAttribute{Computed: true, Description: "The database label."},
		"members":         schema.MapAttribute{Computed: true, ElementType: types.StringType, Description: "Cluster members and their roles."},
		"platform":        schema.StringAttribute{Computed: true, Description: "The backend platform."},
		"port":            schema.Int64Attribute{Computed: true, Description: "The database connection port."},
		"private_network": databaseshared.DataSourceAttributePrivateNetwork,
		"region":          schema.StringAttribute{Computed: true, Description: "The database region."},
		"root_password": schema.StringAttribute{
			Computed:    true,
			Sensitive:   true,
			Description: "The generated root password. This value is stored in Terraform state.",
		},
		"root_username": schema.StringAttribute{
			Computed:    true,
			Sensitive:   true,
			Description: "The root username. This value is stored in Terraform state.",
		},
		"ssl_connection": schema.BoolAttribute{Computed: true, Description: "Whether SSL is required for connections."},
		"status":         schema.StringAttribute{Computed: true, Description: "The current database status."},
		"suspended":      schema.BoolAttribute{Computed: true, Description: "Whether the database is suspended."},
		"type":           schema.StringAttribute{Computed: true, Description: "The database plan type."},
		"updated": schema.StringAttribute{
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			Description: "When the database was last updated.",
		},
		"version": schema.StringAttribute{Computed: true, Description: "The Valkey version."},
		"available_restore_times": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Available discrete snapshot times for restoring this database.",
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
