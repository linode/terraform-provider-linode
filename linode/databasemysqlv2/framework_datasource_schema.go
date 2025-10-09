package databasemysqlv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the MySQL Database.",
			Required:    true,
		},
		"engine_id": schema.StringAttribute{
			Description: "The unique ID of the database engine and version to use. (e.g. mysql/8)",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "A unique, user-defined string referring to the Managed Database.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The Region ID for the Managed Database.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The Linode Instance type used by the Managed Database for its nodes.",
			Computed:    true,
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
		"updates": schema.ObjectAttribute{
			Description:    "Configuration settings for automated patch update maintenance for the Managed Database.",
			AttributeTypes: updatesAttributes,
			Computed:       true,
		},
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
			Description: "A set of pending updates.",
			Computed:    true,
			ElementType: types.ObjectType{AttrTypes: pendingUpdateAttributes},
		},
		"platform": schema.StringAttribute{
			Computed:    true,
			Description: "The back-end platform for relational databases used by the service.",
		},
		"port": schema.Int64Attribute{
			Description: "The access port for this Managed Database.",
			Computed:    true,
		},
		"private_network": schema.SingleNestedAttribute{
			Description: "Restricts access to this database using a virtual private cloud (VPC) " +
				"that you've configured in the region where the database will live.",
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"vpc_id": schema.Int64Attribute{
					Description: "The ID of the virtual private cloud (VPC) " +
						"to restrict access to this database using.",
					Computed: true,
				},
				"subnet_id": schema.Int64Attribute{
					Description: "The ID of the VPC subnet to restrict access to this database using.",
					Computed:    true,
				},
				"public_access": schema.BoolAttribute{
					Description: "If true, clients outside of the VPC can " +
						"connect to the database using a public IP address.",
					Computed: true,
				},
			},
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
		"engine_config_binlog_retention_period": schema.Int64Attribute{
			Computed:    true,
			Description: "The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.",
		},
		"engine_config_mysql_connect_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshake.",
		},
		"engine_config_mysql_default_time_zone": schema.StringAttribute{
			Computed:    true,
			Description: "Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default.",
		},
		"engine_config_mysql_group_concat_max_len": schema.Float64Attribute{
			Computed:    true,
			Description: "The maximum permitted result length in bytes for the GROUP_CONCAT() function.",
		},
		"engine_config_mysql_information_schema_stats_expiry": schema.Int64Attribute{
			Computed:    true,
			Description: "The time, in seconds, before cached statistics expire.",
		},
		"engine_config_mysql_innodb_change_buffer_max_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25.",
		},
		"engine_config_mysql_innodb_flush_neighbors": schema.Int64Attribute{
			Computed:    true,
			Description: "Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent.",
		},
		"engine_config_mysql_innodb_ft_min_token_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.",
		},
		"engine_config_mysql_innodb_ft_server_stopword_table": schema.StringAttribute{
			Computed:    true,
			Description: "This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.",
		},
		"engine_config_mysql_innodb_lock_wait_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.",
		},
		"engine_config_mysql_innodb_log_buffer_size": schema.Int64Attribute{
			Computed:    true,
			Description: "The size in bytes of the buffer that InnoDB uses to write to the log files on disk.",
		},
		"engine_config_mysql_innodb_online_alter_log_max_size": schema.Int64Attribute{
			Computed:    true,
			Description: "The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.",
		},
		"engine_config_mysql_innodb_read_io_threads": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
		},
		"engine_config_mysql_innodb_rollback_on_timeout": schema.BoolAttribute{
			Computed:    true,
			Description: "When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.",
		},
		"engine_config_mysql_innodb_thread_concurrency": schema.Int64Attribute{
			Computed:    true,
			Description: "Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit).",
		},
		"engine_config_mysql_innodb_write_io_threads": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.",
		},
		"engine_config_mysql_interactive_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of seconds the server waits for activity on an interactive connection before closing it.",
		},
		"engine_config_mysql_internal_tmp_mem_storage_engine": schema.StringAttribute{
			Computed:    true,
			Description: "The storage engine for in-memory internal temporary tables.",
		},
		"engine_config_mysql_max_allowed_packet": schema.Int64Attribute{
			Computed:    true,
			Description: "Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M).",
		},
		"engine_config_mysql_max_heap_table_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M).",
		},
		"engine_config_mysql_net_buffer_length": schema.Int64Attribute{
			Computed:    true,
			Description: "Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.",
		},
		"engine_config_mysql_net_read_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of seconds to wait for more data from a connection before aborting the read.",
		},
		"engine_config_mysql_net_write_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of seconds to wait for a block to be written to a connection before aborting the write.",
		},
		"engine_config_mysql_sort_buffer_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K).",
		},
		"engine_config_mysql_sql_mode": schema.StringAttribute{
			Computed:    true,
			Description: "Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.",
		},
		"engine_config_mysql_sql_require_primary_key": schema.BoolAttribute{
			Computed:    true,
			Description: "Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.",
		},
		"engine_config_mysql_tmp_table_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M).",
		},
		"engine_config_mysql_wait_timeout": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of seconds the server waits for activity on a noninteractive connection before closing it.",
		},
	},
}
