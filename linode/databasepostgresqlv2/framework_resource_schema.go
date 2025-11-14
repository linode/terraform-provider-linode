package databasepostgresqlv2

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/databaseshared"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/stringplanmodifiers"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the PostgreSQL Database.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"engine_id": schema.StringAttribute{
			Required:    true,
			Description: "The unique ID of the database engine and version to use. (e.g. postgresql/16)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "A unique, user-defined string referring to the Managed Database.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The Region ID for the Managed Database.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"type": schema.StringAttribute{
			Required:    true,
			Description: "The Linode Instance type used by the Managed Database for its nodes.\n\n",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"ca_cert": schema.StringAttribute{
			Description:   "The base64-encoded SSL CA certificate for the Managed Database.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifiers.UseStateForUnknownIfNotNull()},
		},
		"cluster_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of Linode instance nodes deployed to the Managed Database.",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
			Default: int64default.StaticInt64(1),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"fork_restore_time": schema.StringAttribute{
			Description: "The database timestamp from which it was restored.",
			Optional:    true,
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplaceIf(
					func(
						ctx context.Context,
						sr planmodifier.StringRequest,
						rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
					) {
						rrifr.RequiresReplace = !helper.CompareRFC3339TimeStrings(
							sr.PlanValue.ValueString(),
							sr.StateValue.ValueString(),
						)
					},
					"Triggers replacement when `fork_restore_time` changes",
					"Changing `fork_restore_time` forces a new resource.",
				),
			},
		},
		"fork_source": schema.Int64Attribute{
			Description: "The ID of the database that was forked from.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
				int64planmodifier.RequiresReplace(),
			},
		},
		"suspended": schema.BoolAttribute{
			Description:   "Whether this database is suspended.",
			Computed:      true,
			Optional:      true,
			Default:       booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"updates": databaseshared.ResourceAttributeUpdates,
		"created": schema.StringAttribute{
			Description: "When this Managed Database was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"encrypted": schema.BoolAttribute{
			Description: "Whether the Managed Databases is encrypted.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine in engine/version format.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"host_primary": schema.StringAttribute{
			Description: "The primary host for the Managed Database.",
			Computed:    true,
		},
		"host_secondary": schema.StringAttribute{
			Description: "The secondary/private host for the Managed Database.",
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
		"pending_updates": databaseshared.ResourceAttributePendingUpdates,
		"platform": schema.StringAttribute{
			Computed:      true,
			Description:   "The back-end platform for relational databases used by the service.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"port": schema.Int64Attribute{
			Description:   "The access port for this Managed Database.",
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"private_network": databaseshared.ResourceAttributePrivateNetwork,
		"root_password": schema.StringAttribute{
			Description:   "The randomly generated root password for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifiers.UseStateForUnknownIfNotNull()},
		},
		"root_username": schema.StringAttribute{
			Description:   "The root username for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifiers.UseStateForUnknownIfNotNull()},
		},
		"ssl_connection": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The operating status of the Managed Database.",
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
			Optional:    true,
			Computed:    true,
			Description: "Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)",
			PlanModifiers: []planmodifier.Float64{
				float64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Float64{
				float64validator.Between(0.0, 1.0),
			},
		},
		"engine_config_pg_autovacuum_analyze_threshold": schema.Int32Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.",
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int32{
				int32validator.Between(0, 2147483647),
			},
		},
		"engine_config_pg_autovacuum_max_workers": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 20),
			},
		},
		"engine_config_pg_autovacuum_naptime": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 86400),
			},
		},
		"engine_config_pg_autovacuum_vacuum_cost_delay": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(-1, 100),
			},
		},
		"engine_config_pg_autovacuum_vacuum_cost_limit": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(-1, 10000),
			},
		},
		"engine_config_pg_autovacuum_vacuum_scale_factor": schema.Float64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)",
			PlanModifiers: []planmodifier.Float64{
				float64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Float64{
				float64validator.Between(0.0, 1.0),
			},
		},
		"engine_config_pg_autovacuum_vacuum_threshold": schema.Int32Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples",
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int32{
				int32validator.Between(0, 2147483647),
			},
		},
		"engine_config_pg_bgwriter_delay": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(10, 10000),
			},
		},
		"engine_config_pg_bgwriter_flush_after": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 2048),
			},
		},
		"engine_config_pg_bgwriter_lru_maxpages": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 1073741823),
			},
		},
		"engine_config_pg_bgwriter_lru_multiplier": schema.Float64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.",
			PlanModifiers: []planmodifier.Float64{
				float64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Float64{
				float64validator.Between(0.0, 10.0),
			},
		},
		"engine_config_pg_deadlock_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(500, 1800000),
			},
		},
		"engine_config_pg_default_toast_compression": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the default TOAST compression method for values of compressible columns (the default is lz4).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("lz4", "pglz"),
			},
		},
		"engine_config_pg_idle_in_transaction_session_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Time out sessions with open transactions after this number of milliseconds",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 604800000),
			},
		},
		"engine_config_pg_jit": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Controls system-wide use of Just-in-Time Compilation (JIT).",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"engine_config_pg_max_files_per_process": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum number of files that can be open per process",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1000, 4096),
			},
		},
		"engine_config_pg_max_locks_per_transaction": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum locks per transaction",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(64, 6400),
			},
		},
		"engine_config_pg_max_logical_replication_workers": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers)",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(4, 64),
			},
		},
		"engine_config_pg_max_parallel_workers": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sets the maximum number of workers that the system can support for parallel queries",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 96),
			},
		},
		"engine_config_pg_max_parallel_workers_per_gather": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sets the maximum number of workers that can be started by a single Gather or Gather Merge node",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 96),
			},
		},
		"engine_config_pg_max_pred_locks_per_transaction": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum predicate locks per transaction",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(64, 5120),
			},
		},
		"engine_config_pg_max_replication_slots": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum replication slots",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(8, 64),
			},
		},
		"engine_config_pg_max_slot_wal_keep_size": schema.Int32Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.",
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int32{
				int32validator.Between(-1, 2147483647),
			},
		},
		"engine_config_pg_max_stack_depth": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Maximum depth of the stack in bytes",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(2097152, 6291456),
			},
		},
		"engine_config_pg_max_standby_archive_delay": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Max standby archive delay in milliseconds",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 43200000),
			},
		},
		"engine_config_pg_max_standby_streaming_delay": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Max standby streaming delay in milliseconds",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 43200000),
			},
		},
		"engine_config_pg_max_wal_senders": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL maximum WAL senders",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(20, 64),
			},
		},
		"engine_config_pg_max_worker_processes": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sets the maximum number of background processes that the system can support",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(8, 96),
			},
		},
		"engine_config_pg_password_encryption": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Chooses the algorithm for encrypting passwords.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("md5", "scram-sha-256"),
			},
		},
		"engine_config_pg_pg_partman_bgw_interval": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sets the time interval to run pg_partman's scheduled tasks",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(3600, 604800),
			},
		},
		"engine_config_pg_pg_partman_bgw_role": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Controls which role to use for pg_partman's scheduled background tasks.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(regexp.MustCompile(`^[_A-Za-z0-9][-._A-Za-z0-9]{0,63}$`), "must be in the format 'myrolename'"),
			},
		},
		"engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Enables or disables query plan monitoring",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"engine_config_pg_pg_stat_monitor_pgsm_max_buckets": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sets the maximum number of buckets",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1, 10),
			},
		},
		"engine_config_pg_pg_stat_statements_track": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("all", "top", "none"),
			},
		},
		"engine_config_pg_temp_file_limit": schema.Int32Attribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL temporary file limit in KiB, -1 for unlimited",
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int32{
				int32validator.Between(-1, 2147483647),
			},
		},
		"engine_config_pg_timezone": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "PostgreSQL service timezone",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(regexp.MustCompile(`^[\w/]*$`), "must be in the format 'Area/Location' (e.g., 'Europe/Helsinki')"),
			},
		},
		"engine_config_pg_track_activity_query_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies the number of bytes reserved to track the currently executing command for each active session.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(1024, 10240),
			},
		},
		"engine_config_pg_track_commit_timestamp": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Record commit time of transactions.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("off", "on"),
			},
		},
		"engine_config_pg_track_functions": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Enables tracking of function call counts and time used.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("all", "pl", "none"),
			},
		},
		"engine_config_pg_track_io_timing": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("off", "on"),
			},
		},
		"engine_config_pg_wal_sender_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"engine_config_pg_wal_writer_delay": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(10, 200),
			},
		},
		"engine_config_pg_stat_monitor_enable": schema.BoolAttribute{
			Optional:      true,
			Computed:      true,
			Description:   "Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted. When this extension is enabled, pg_stat_statements results for utility commands are unreliable.",
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"engine_config_pglookout_max_failover_replication_time_lag": schema.Int64Attribute{
			Optional:      true,
			Computed:      true,
			Description:   "Number of seconds of master unavailability before triggering database failover to standby.",
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			Validators: []validator.Int64{
				int64validator.Between(10, 999999),
			},
		},
		"engine_config_shared_buffers_percentage": schema.Float64Attribute{
			Optional:      true,
			Computed:      true,
			Description:   "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.",
			PlanModifiers: []planmodifier.Float64{float64planmodifier.UseStateForUnknown()},
			Validators: []validator.Float64{
				float64validator.Between(20.0, 60.0),
			},
		},
		"engine_config_work_mem": schema.Int64Attribute{
			Optional:      true,
			Computed:      true,
			Description:   "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).",
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			Validators: []validator.Int64{
				int64validator.Between(1, 1024),
			},
		},
	},
}
