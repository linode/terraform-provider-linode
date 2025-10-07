---
page_title: "Linode: linode_database_mysql_v2"
description: |-
  Provides information about a Linode MySQL Database.
---

# Data Source: linode\_database\_mysql\_v2

Provides information about a Linode MySQL Database.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-databases-mysql-instance).

## Example Usage

Get information about a MySQL database:

```hcl
data "linode_database_mysql_v2" "my-db" {
  id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `id` - The ID of the MySQL database.

## Attributes Reference

The `linode_database_mysql_v2` data source exports the following attributes:

* `allow_list` - A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format. Use `linode_database_access_controls` to manage your allow list separately.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database.

* `cluster_size` - The number of Linode Instance nodes deployed to the Managed Database. (default `1`)

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Databases is encrypted.

* `engine` - The Managed Database engine. (e.g. `mysql`)

* `engine_id` - The Managed Database engine in engine/version format. (e.g. `mysql`)

* `fork_restore_time` - The database timestamp from which it was restored.

* `fork_source` - The ID of the database that was forked from.

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private host for the managed database.

* `label` - A unique, user-defined string referring to the Managed Database.

* [`pending_updates`](#pending_updates) - A set of pending updates.

* `platform` - The back-end platform for relational databases used by the service.

* `port` - The access port for this Managed Database.

* [`private_network`](#private_network) - Restricts access to this database using a virtual private cloud (VPC).

* `region` - The region to use for the Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `ssl_connection` - Whether to require SSL credentials to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `suspended` - Whether this Managed Database is suspended.

* `type` - The Linode Instance type used for the nodes of the Managed Database.

* `updated` - When this Managed Database was last updated.

* [`updates`](#updates) - Configuration settings for automated patch update maintenance for the Managed Database.

* `version` - The Managed Database engine version. (e.g. `13.2`)

* `engine_config_binlog_retention_period` - The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default, for example if using the MySQL Debezium Kafka connector.

* `engine_config_mysql_connect_timeout` - The number of seconds that the mysqld server waits for a connect packet before responding with "Bad handshake".

* `engine_config_mysql_default_time_zone` - Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or `SYSTEM` to use the MySQL server default.

* `engine_config_mysql_group_concat_max_len` - The maximum permitted result length in bytes for the `GROUP_CONCAT()` function.

* `engine_config_mysql_information_schema_stats_expiry` - The time, in seconds, before cached statistics expire.

* `engine_config_mysql_innodb_change_buffer_max_size` - Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25.

* `engine_config_mysql_innodb_flush_neighbors` - Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent.

* `engine_config_mysql_innodb_ft_min_token_size` - Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_ft_server_stopword_table` - This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.

* `engine_config_mysql_innodb_lock_wait_timeout` - The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.

* `engine_config_mysql_innodb_log_buffer_size` - The size in bytes of the buffer that InnoDB uses to write to the log files on disk.

* `engine_config_mysql_innodb_online_alter_log_max_size` - The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.

* `engine_config_mysql_innodb_read_io_threads` - The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_rollback_on_timeout` - When enabled, a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_thread_concurrency` - Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit).

* `engine_config_mysql_innodb_write_io_threads` - The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_interactive_timeout` - The number of seconds the server waits for activity on an interactive connection before closing it.

* `engine_config_mysql_internal_tmp_mem_storage_engine` - The storage engine for in-memory internal temporary tables.

* `engine_config_mysql_max_allowed_packet` - Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M).

* `engine_config_mysql_max_heap_table_size` - Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M).

* `engine_config_mysql_net_buffer_length` - Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_net_read_timeout` - The number of seconds to wait for more data from a connection before aborting the read.

* `engine_config_mysql_net_write_timeout` - The number of seconds to wait for a block to be written to a connection before aborting the write.

* `engine_config_mysql_sort_buffer_size` - Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K).

* `engine_config_mysql_sql_mode` - Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.

* `engine_config_mysql_sql_require_primary_key` - Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.

* `engine_config_mysql_tmp_table_size` - Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M).

* `engine_config_mysql_wait_timeout` - The number of seconds the server waits for activity on a noninteractive connection before closing it.

## pending_updates

The following arguments are exposed by each entry in the `pending_updates` attribute:

* `deadline` - The time when a mandatory update needs to be applied.

* `description` - A description of the update.

* `planned_for` - The date and time a maintenance update will be applied.

## updates

The following arguments are supported in the `updates` specification block:

* `day_of_week` - The day to perform maintenance. (`monday`, `tuesday`, ...)

* `duration` - The maximum maintenance window time in hours. (`1`..`3`)

* `frequency` - The frequency at which maintenance occurs. (`weekly`)

* `hour_of_day` - The hour to begin maintenance based in UTC time. (`0`..`23`)

## private_network

The following arguments are exposed by the `private_network` attribute:

* `vpc_id` - The ID of the virtual private cloud (VPC) to restrict access to this database using.

* `subnet_id` - The ID of the VPC subnet to restrict access to this database using.

* `public_access` - If true, clients outside the VPC can connect to the database using a public IP address.
