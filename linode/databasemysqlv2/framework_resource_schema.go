package databasemysqlv2

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var (
	updatesAttributes = map[string]attr.Type{
		"day_of_week": types.Int64Type,
		"duration":    types.Int64Type,
		"frequency":   types.StringType,
		"hour_of_day": types.Int64Type,
	}

	pendingUpdateAttributes = map[string]attr.Type{
		"deadline":    timetypes.RFC3339Type{},
		"description": types.StringType,
		"planned_for": timetypes.RFC3339Type{},
	}
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the MySQL Database.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"engine_id": schema.StringAttribute{
			Required:    true,
			Description: "The unique ID of the database engine and version to use. (e.g. mysql/8)",
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
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
		"updates": schema.ObjectAttribute{
			Description:    "Configuration settings for automated patch update maintenance for the Managed Database.",
			AttributeTypes: updatesAttributes,
			Computed:       true,
			Optional:       true,
			PlanModifiers:  []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
		},
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
		"pending_updates": schema.SetAttribute{
			Description:   "A set of pending updates.",
			Computed:      true,
			ElementType:   types.ObjectType{AttrTypes: pendingUpdateAttributes},
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
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
		"root_password": schema.StringAttribute{
			Description:   "The randomly generated root password for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"root_username": schema.StringAttribute{
			Description:   "The root username for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
		"engine_config_binlog_retention_period": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_connect_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshake.",
			Validators: []validator.Int64{
				int64validator.Between(2, 3600),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_default_time_zone": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(2, 100),
				stringvalidator.RegexMatches(regexp.MustCompile(`^([-+][\d:]*|[\w/]*)$`), "must be a valid time zone offset, name, or 'SYSTEM'"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_group_concat_max_len": schema.Float64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The maximum permitted result length in bytes for the GROUP_CONCAT() function.",
			Validators: []validator.Float64{
				float64validator.Between(4, 1.8446744073709552e+19),
			},
			PlanModifiers: []planmodifier.Float64{float64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_information_schema_stats_expiry": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The time, in seconds, before cached statistics expire.",
			Validators: []validator.Int64{
				int64validator.Between(900, 31536000),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_change_buffer_max_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25.",
			Validators: []validator.Int64{
				int64validator.Between(0, 50),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_flush_neighbors": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent.",
			Validators: []validator.Int64{
				int64validator.Between(0, 2),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_ft_min_token_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.",
			Validators: []validator.Int64{
				int64validator.Between(0, 16),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_ft_server_stopword_table": schema.StringAttribute{
			Optional:    true,
			Description: "This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.",
			Validators: []validator.String{
				stringvalidator.LengthAtMost(1024),
				stringvalidator.RegexMatches(regexp.MustCompile(`^.+/.+$`), "must be in the format 'db_name/table_name'"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_lock_wait_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.",
			Validators: []validator.Int64{
				int64validator.Between(1, 3600),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_log_buffer_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The size in bytes of the buffer that InnoDB uses to write to the log files on disk.",
			Validators: []validator.Int64{
				int64validator.Between(1048576, 4294967295),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_online_alter_log_max_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.",
			Validators: []validator.Int64{
				int64validator.Between(65536, 1099511627776),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_read_io_threads": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
			Validators: []validator.Int64{
				int64validator.Between(1, 64),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_rollback_on_timeout": schema.BoolAttribute{
			Optional:      true,
			Computed:      true,
			Description:   "When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.",
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_thread_concurrency": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit).",
			Validators: []validator.Int64{
				int64validator.Between(0, 1000),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_innodb_write_io_threads": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
			Validators: []validator.Int64{
				int64validator.Between(1, 64),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_interactive_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of seconds the server waits for activity on an interactive connection before closing it.",
			Validators: []validator.Int64{
				int64validator.Between(30, 604800),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_internal_tmp_mem_storage_engine": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The storage engine for in-memory internal temporary tables.",
			Validators: []validator.String{
				stringvalidator.OneOf("TempTable", "MEMORY"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_max_allowed_packet": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M).",
			Validators: []validator.Int64{
				int64validator.Between(102400, 1073741824),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_max_heap_table_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M).",
			Validators: []validator.Int64{
				int64validator.Between(1048576, 1073741824),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_net_buffer_length": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.",
			Validators: []validator.Int64{
				int64validator.Between(1024, 1048576),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_net_read_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of seconds to wait for more data from a connection before aborting the read.",
			Validators: []validator.Int64{
				int64validator.Between(1, 3600),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_net_write_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of seconds to wait for a block to be written to a connection before aborting the write.",
			Validators: []validator.Int64{
				int64validator.Between(1, 3600),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_sort_buffer_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K).",
			Validators: []validator.Int64{
				int64validator.Between(32768, 1073741824),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_sql_mode": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.",
			Validators: []validator.String{
				stringvalidator.LengthAtMost(1024),
				stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Z_]*(,[A-Z_]+)*$`), "must be a valid SQL mode format"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_sql_require_primary_key": schema.BoolAttribute{
			Optional:      true,
			Computed:      true,
			Description:   "Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.",
			PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_tmp_table_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M).",
			Validators: []validator.Int64{
				int64validator.Between(1048576, 1073741824),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"engine_config_mysql_wait_timeout": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of seconds the server waits for activity on a noninteractive connection before closing it.",
			Validators: []validator.Int64{
				int64validator.Between(1, 2147483),
			},
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
	},
}
