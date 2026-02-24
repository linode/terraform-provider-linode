package databasepostgresqlv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/databaseshared"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the PostgreSQL Database.",
			Required:    true,
		},

		"engine_id": schema.StringAttribute{
			Computed:    true,
			Description: "The unique ID of the database engine and version to use. (e.g. postgresql/16)",
		},
		"label": schema.StringAttribute{
			Computed:    true,
			Description: "A unique, user-defined string referring to the Managed Database.",
		},
		"region": schema.StringAttribute{
			Computed:    true,
			Description: "The Region ID for the Managed Database.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "The Linode Instance type used by the Managed Database for its nodes.\n\n",
		},

		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
		},
		"ca_cert": schema.StringAttribute{
			Description: "The base64-encoded SSL CA certificate for the Managed Database.",
			Computed:    true,
			Sensitive:   true,
		},
		"cluster_size": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of Linode instance nodes deployed to the Managed Database.",
		},
		"fork_restore_time": schema.StringAttribute{
			Description: "The database timestamp from which it was restored.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"fork_source": schema.Int64Attribute{
			Description: "The ID of the database that was forked from.",
			Computed:    true,
		},
		"private_network": databaseshared.DataSourceAttributePrivateNetwork,
		"updates":         databaseshared.ResourceAttributeUpdates,

		"created": schema.StringAttribute{
			Description: "When this Managed Database was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"encrypted": schema.BoolAttribute{
			Description: "Whether the Managed Databases is encrypted.",
			Computed:    true,
		},
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine in engine/version format.",
			Computed:    true,
		},
		"host_primary": schema.StringAttribute{
			Description: "The primary host for the Managed Database.",
			Computed:    true,
		},
		"host_secondary": schema.StringAttribute{
			Description:        "The secondary/private host for the Managed Database.",
			Computed:           true,
			DeprecationMessage: "Use host_standby instead.",
		},
		"host_standby": schema.StringAttribute{
			Description: "The standby host for the Managed Database.",
			Computed:    true,
		},
		"members": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "A mapping between IP addresses and strings designating them as primary or failover.",
		},
		"oldest_restore_time": schema.StringAttribute{
			Description: "The oldest time to which a database can be restored.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"pending_updates": databaseshared.DataSourceAttributePendingUpdates,
		"platform": schema.StringAttribute{
			Computed:    true,
			Description: "The back-end platform for relational databases used by the service.",
		},
		"port": schema.Int64Attribute{
			Description: "The access port for this Managed Database.",
			Computed:    true,
		},
		"root_password": schema.StringAttribute{
			Description: "The randomly generated root password for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"root_username": schema.StringAttribute{
			Description: "The root username for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
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
			Description: "Whether this database is suspended.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this Managed Database was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"version": schema.StringAttribute{
			Description: "The Managed Database engine version.",
			Computed:    true,
		},
		"engine_config_pg_autovacuum_analyze_scale_factor": schema.Float64Attribute{
			Computed:    true,
			Description: "Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)",
		},
		"engine_config_pg_autovacuum_analyze_threshold": schema.Int32Attribute{
			Computed:    true,
			Description: "Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.",
		},
		"engine_config_pg_autovacuum_max_workers": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.",
		},
		"engine_config_pg_autovacuum_naptime": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute",
		},
		"engine_config_pg_autovacuum_vacuum_cost_delay": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds",
		},
		"engine_config_pg_autovacuum_vacuum_cost_limit": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.",
		},
		"engine_config_pg_autovacuum_vacuum_scale_factor": schema.Float64Attribute{
			Computed:    true,
			Description: "Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)",
		},
		"engine_config_pg_autovacuum_vacuum_threshold": schema.Int32Attribute{
			Computed:    true,
			Description: "Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples",
		},
		"engine_config_pg_bgwriter_delay": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.",
		},
		"engine_config_pg_bgwriter_flush_after": schema.Int64Attribute{
			Computed:    true,
			Description: "Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.",
		},
		"engine_config_pg_bgwriter_lru_maxpages": schema.Int64Attribute{
			Computed:    true,
			Description: "In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.",
		},
		"engine_config_pg_bgwriter_lru_multiplier": schema.Float64Attribute{
			Computed:    true,
			Description: "The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.",
		},
		"engine_config_pg_deadlock_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.",
		},
		"engine_config_pg_default_toast_compression": schema.StringAttribute{
			Computed:    true,
			Description: "Specifies the default TOAST compression method for values of compressible columns (the default is lz4).",
		},
		"engine_config_pg_idle_in_transaction_session_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "Time out sessions with open transactions after this number of milliseconds",
		},
		"engine_config_pg_jit": schema.BoolAttribute{
			Computed:    true,
			Description: "Controls system-wide use of Just-in-Time Compilation (JIT).",
		},
		"engine_config_pg_max_files_per_process": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum number of files that can be open per process",
		},
		"engine_config_pg_max_locks_per_transaction": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum locks per transaction",
		},
		"engine_config_pg_max_logical_replication_workers": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers)",
		},
		"engine_config_pg_max_parallel_workers": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the maximum number of workers that the system can support for parallel queries",
		},
		"engine_config_pg_max_parallel_workers_per_gather": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the maximum number of workers that can be started by a single Gather or Gather Merge node",
		},
		"engine_config_pg_max_pred_locks_per_transaction": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum predicate locks per transaction",
		},
		"engine_config_pg_max_replication_slots": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum replication slots",
		},
		"engine_config_pg_max_slot_wal_keep_size": schema.Int32Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.",
		},
		"engine_config_pg_max_stack_depth": schema.Int64Attribute{
			Computed:    true,
			Description: "Maximum depth of the stack in bytes",
		},
		"engine_config_pg_max_standby_archive_delay": schema.Int64Attribute{
			Computed:    true,
			Description: "Max standby archive delay in milliseconds",
		},
		"engine_config_pg_max_standby_streaming_delay": schema.Int64Attribute{
			Computed:    true,
			Description: "Max standby streaming delay in milliseconds",
		},
		"engine_config_pg_max_wal_senders": schema.Int64Attribute{
			Computed:    true,
			Description: "PostgreSQL maximum WAL senders",
		},
		"engine_config_pg_max_worker_processes": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the maximum number of background processes that the system can support",
		},
		"engine_config_pg_password_encryption": schema.StringAttribute{
			Computed:    true,
			Description: "Chooses the algorithm for encrypting passwords.",
		},
		"engine_config_pg_pg_partman_bgw_interval": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the time interval to run pg_partman's scheduled tasks",
		},
		"engine_config_pg_pg_partman_bgw_role": schema.StringAttribute{
			Computed:    true,
			Description: "Controls which role to use for pg_partman's scheduled background tasks.",
		},
		"engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan": schema.BoolAttribute{
			Computed:    true,
			Description: "Enables or disables query plan monitoring",
		},
		"engine_config_pg_pg_stat_monitor_pgsm_max_buckets": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the maximum number of buckets",
		},
		"engine_config_pg_pg_stat_statements_track": schema.StringAttribute{
			Computed:    true,
			Description: "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.",
		},
		"engine_config_pg_temp_file_limit": schema.Int32Attribute{
			Computed:    true,
			Description: "PostgreSQL temporary file limit in KiB, -1 for unlimited",
		},
		"engine_config_pg_timezone": schema.StringAttribute{
			Computed:    true,
			Description: "PostgreSQL service timezone",
		},
		"engine_config_pg_track_activity_query_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies the number of bytes reserved to track the currently executing command for each active session.",
		},
		"engine_config_pg_track_commit_timestamp": schema.StringAttribute{
			Computed:    true,
			Description: "Record commit time of transactions.",
		},
		"engine_config_pg_track_functions": schema.StringAttribute{
			Computed:    true,
			Description: "Enables tracking of function call counts and time used.",
		},
		"engine_config_pg_track_io_timing": schema.StringAttribute{
			Computed:    true,
			Description: "Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.",
		},
		"engine_config_pg_wal_sender_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout.",
		},
		"engine_config_pg_wal_writer_delay": schema.Int64Attribute{
			Computed:    true,
			Description: "WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance.",
		},
		"engine_config_pg_stat_monitor_enable": schema.BoolAttribute{
			Computed:    true,
			Description: "Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted. When this extension is enabled, pg_stat_statements results for utility commands are unreliable.",
		},
		"engine_config_pglookout_max_failover_replication_time_lag": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of seconds of master unavailability before triggering database failover to standby.",
		},
		"engine_config_shared_buffers_percentage": schema.Float64Attribute{
			Computed:    true,
			Description: "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.",
		},
		"engine_config_work_mem": schema.Int64Attribute{
			Computed:    true,
			Description: "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).",
		},
	},
}
