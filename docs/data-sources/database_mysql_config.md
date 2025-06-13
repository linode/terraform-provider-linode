---
page_title: "Linode: linode_database_mysql_config"
description: |-
  Provides information about a Linode MySQL Database's Configuration Options.
---

# Data Source: linode\_database\_mysql\_config

Provides information about a Linode MySQL Database's Configuration Options.
For more information, see the [Linode APIv4 docs](TODO).

## Example Usage

Get information about a MySQL database's configuration options:

```hcl
data "linode_database_mysql_config" "my-db-config" {}
```

## Attributes Reference

The `linode_database_mysql_config` data source exports the following attributes:

* [`binlog_retention_period`](#binlog_retention_period) - The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.

* [`mysql`](#mysql) - MySQL configuration settings.

## binlog_retention_period

The following arguments are supported in the `binlog_retention_period` specification block:

* `description` - The description of `binlog_retention_period`.

* `example` - An example of a valid value for `binlog_retention_period`.

* `maximum` - The maximum valid value of `binlog_retention_period`.

* `minimum` - The minimum valid value of `binlog_retention_period`.

* `requires_restart` - Whether changing the value `binlog_retention_period` requires the DB to restart.

* `type` - The type of the value of `binlog_retention_period`.

## mysql

The following arguments are supported in the `mysql` specification block:

* [`connect_timeout`](#connect_timeout) - The number of seconds that the mysqld server waits for a connect packet before responding with "Bad handshake".

* [`default_time_zone`](#default_time_zone) - Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or `SYSTEM` to use the MySQL server default.

* [`group_concat_max_len`](#group_concat_max_len) - The maximum permitted result length in bytes for the `GROUP_CONCAT()` function.

* [`information_schema_stats_expiry`](#information_schema_stats_expiry) - The time, in seconds, before cached statistics expire.

* [`innodb_change_buffer_max_size`](#innodb_change_buffer_max_size) - Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25.

* [`innodb_flush_neighbors`](#innodb_flush_neighbors) - Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent.

* [`innodb_ft_min_token_size`](#innodb_ft_min_token_size) - Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.

* [`innodb_ft_server_stopword_table`](#innodb_ft_server_stopword_table) - This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.

* [`innodb_lock_wait_timeout`](#innodb_lock_wait_timeout) - The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.

* [`innodb_log_buffer_size`](#innodb_log_buffer_size) - The size in bytes of the buffer that InnoDB uses to write to the log files on disk.

* [`innodb_online_alter_log_max_size`](#innodb_online_alter_log_max_size) - The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.

* [`innodb_read_io_threads`](#innodb_read_io_threads) - The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* [`innodb_rollback_on_timeout`](#innodb_rollback_on_timeout) - When enabled, a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.

* [`innodb_thread_concurrency`](#innodb_thread_concurrency) - Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit).

* [`innodb_write_io_threads`](#innodb_write_io_threads) - The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* [`interactive_timeout`](#interactive_timeout) - The number of seconds the server waits for activity on an interactive connection before closing it.

* [`internal_tmp_mem_storage_engine`](#internal_tmp_mem_storage_engine) - The storage engine for in-memory internal temporary tables.

* [`max_allowed_packet`](#max_allowed_packet) - Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M).

* [`max_heap_table_size`](#max_heap_table_size) - Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M).

* [`net_buffer_length`](#net_buffer_length) - Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.

* [`net_read_timeout`](#net_read_timeout) - The number of seconds to wait for more data from a connection before aborting the read.

* [`net_write_timeout`](#net_write_timeout) - The number of seconds to wait for a block to be written to a connection before aborting the write.

* [`sort_buffer_size`](#sort_buffer_size) - Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K).

* [`sql_mode`](#sql_mode) - Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.

* [`sql_require_primary_key`](#sql_require_primary_key) - Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.

* [`tmp_table_size`](#tmp_table_size) - Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M).

* [`wait_timeout`](#wait_timeout) - The number of seconds the server waits for activity on a noninteractive connection before closing it.

## connect_timeout

The following arguments are supported in the `connect_timeout` specification block:

* `description` - The description of `connect_timeout`.

* `example` - An example of a valid value for `connect_timeout`.

* `maximum` - The maximum valid value of  `connect_timeout`.

* `minimum` - The minimum valid value of  `connect_timeout`.

* `requires_restart` - Whether changing the value `connect_timeout` requires the DB to restart.

* `type` - The type of the value of `connect_timeout`.

## default_time_zone

The following arguments are supported in the `default_time_zone` specification block:

* `description` - The description of `default_time_zone`.

* `example` - An example of a valid value for `default_time_zone`.

* `maxLength` - The maximum length of the `default_time_zone` value.

* `minLength` - The minimum length of the `default_time_zone` value.

* `pattern` - A regular expression that the `default_time_zone` value must match.

* `requires_restart` - Whether changing the value `default_time_zone` requires the DB to restart.

* `type` - The type of the value of `default_time_zone`.

## group_concat_max_len

The following arguments are supported in the `group_concat_max_len` specification block:

* `description` - The description of `group_concat_max_len`.

* `example` - An example of a valid value for `group_concat_max_len`.

* `maximum` - The maximum valid value of `group_concat_max_len`.

* `minimum` - The minimum valid value of `group_concat_max_len`.

* `requires_restart` - Whether changing the value `group_concat_max_len` requires the DB to restart.

* `type` - The type of the value of `group_concat_max_len`.

## information_schema_stats_expiry

The following arguments are supported in the `information_schema_stats_expiry` specification block:

* `description` - The description of `information_schema_stats_expiry`.

* `example` - An example of a valid value for `information_schema_stats_expiry`.

* `maximum` - The maximum valid value of `information_schema_stats_expiry`.

* `minimum` - The minimum valid value of `information_schema_stats_expiry`.

* `requires_restart` - Whether changing the value `information_schema_stats_expiry` requires the DB to restart.

* `type` - The type of the value of `information_schema_stats_expiry`.

## innodb_change_buffer_max_size

The following arguments are supported in the `innodb_change_buffer_max_size` specification block:

* `description` - The description of `innodb_change_buffer_max_size`.

* `example` - An example of a valid value for `innodb_change_buffer_max_size`.

* `maximum` - The maximum valid value of `innodb_change_buffer_max_size`.

* `minimum` - The minimum valid value of `innodb_change_buffer_max_size`.

* `requires_restart` - Whether changing the value `innodb_change_buffer_max_size` requires the DB to restart.

* `type` - The type of the value of `innodb_change_buffer_max_size`.

## innodb_flush_neighbors

The following arguments are supported in the `innodb_flush_neighbors` specification block:

* `description` - The description of `innodb_flush_neighbors`.

* `example` - An example of a valid value for `innodb_flush_neighbors`.

* `maximum` - The maximum valid value of `innodb_flush_neighbors`.

* `minimum` - The minimum valid value of `innodb_flush_neighbors`.

* `requires_restart` - Whether changing the value `innodb_flush_neighbors` requires the DB to restart.

* `type` - The type of the value of `innodb_flush_neighbors`.

## innodb_ft_min_token_size

The following arguments are supported in the `innodb_ft_min_token_size` specification block:

* `description` - The description of `innodb_ft_min_token_size`.

* `example` - An example of a valid value for `innodb_ft_min_token_size`.

* `maximum` - The maximum valid value of `innodb_ft_min_token_size`.

* `minimum` - The minimum valid value of `innodb_ft_min_token_size`.

* `requires_restart` - Whether changing the value `innodb_ft_min_token_size` requires the DB to restart.

* `type` - The type of the value of `innodb_ft_min_token_size`.

## innodb_ft_server_stopword_table

The following arguments are supported in the `innodb_ft_server_stopword_table` specification block:

* `description` - The description of `innodb_ft_server_stopword_table`.

* `example` - An example of a valid value for `innodb_ft_server_stopword_table`.

* `maxLength` - The maximum length of the value for `innodb_ft_server_stopword_table`.

* `pattern` - A regex pattern that a value of `innodb_ft_server_stopword_table` must match.

* `requires_restart` - Whether changing the value `innodb_ft_server_stopword_table` requires the DB to restart.

* `type` - The type of the value of `innodb_ft_server_stopword_table`.

## innodb_lock_wait_timeout

The following arguments are supported in the `innodb_lock_wait_timeout` specification block:

* `description` - The description of `innodb_lock_wait_timeout`.

* `example` - An example of a valid value for `innodb_lock_wait_timeout`.

* `maximum` - The maximum valid value of `innodb_lock_wait_timeout`.

* `minimum` - The minimum valid value of `innodb_lock_wait_timeout`.

* `requires_restart` - Whether changing the value `innodb_lock_wait_timeout` requires the DB to restart.

* `type` - The type of the value of `innodb_lock_wait_timeout`.

## innodb_log_buffer_size

The following arguments are supported in the `innodb_log_buffer_size` specification block:

* `description` - The description of `innodb_log_buffer_size`.

* `example` - An example of a valid value for `innodb_log_buffer_size`.

* `maximum` - The maximum valid value of `innodb_log_buffer_size`.

* `minimum` - The minimum valid value of `innodb_log_buffer_size`.

* `requires_restart` - Whether changing the value `innodb_log_buffer_size` requires the DB to restart.

* `type` - The type of the value of `innodb_log_buffer_size`.

## innodb_online_alter_log_max_size

The following arguments are supported in the `innodb_online_alter_log_max_size` specification block:

* `description` - The description of `innodb_online_alter_log_max_size`.

* `example` - An example of a valid value for `innodb_online_alter_log_max_size`.

* `maximum` - The maximum valid value of `innodb_online_alter_log_max_size`.

* `minimum` - The minimum valid value of `innodb_online_alter_log_max_size`.

* `requires_restart` - Whether changing the value `innodb_online_alter_log_max_size` requires the DB to restart.

* `type` - The type of the value of `innodb_online_alter_log_max_size`.

## innodb_read_io_threads

The following arguments are supported in the `innodb_read_io_threads` specification block:

* `description` - The description of `innodb_read_io_threads`.

* `example` - An example of a valid value for `innodb_read_io_threads`.

* `maximum` - The maximum valid value of `innodb_read_io_threads`.

* `minimum` - The minimum valid value of `innodb_read_io_threads`.

* `requires_restart` - Whether changing the value `innodb_read_io_threads` requires the DB to restart.

* `type` - The type of the value of `innodb_read_io_threads`.

## innodb_rollback_on_timeout

The following arguments are supported in the `innodb_rollback_on_timeout` specification block:

* `description` - The description of `innodb_rollback_on_timeout`.

* `example` - An example of a valid value for `innodb_rollback_on_timeout`.

* `requires_restart` - Whether changing the value `innodb_rollback_on_timeout` requires the DB to restart.

* `type` - The type of the value of `innodb_rollback_on_timeout`.

## innodb_thread_concurrency

The following arguments are supported in the `innodb_thread_concurrency` specification block:

* `description` - The description of `innodb_thread_concurrency`.

* `example` - An example of a valid value for `innodb_thread_concurrency`.

* `maximum` - The maximum valid value of `innodb_thread_concurrency`.

* `minimum` - The minimum valid value of `innodb_thread_concurrency`.

* `requires_restart` - Whether changing the value `innodb_thread_concurrency` requires the DB to restart.

* `type` - The type of the value of `innodb_thread_concurrency`.

## innodb_write_io_threads

The following arguments are supported in the `innodb_write_io_threads` specification block:

* `description` - The description of `innodb_write_io_threads`.

* `example` - An example of a valid value for `innodb_write_io_threads`.

* `maximum` - The maximum valid value of `innodb_write_io_threads`.

* `minimum` - The minimum valid value of `innodb_write_io_threads`.

* `requires_restart` - Whether changing the value `innodb_write_io_threads` requires the DB to restart.

* `type` - The type of the value of `innodb_write_io_threads`.

## interactive_timeout

The following arguments are supported in the `interactive_timeout` specification block:

* `description` - The description of `interactive_timeout`.

* `example` - An example of a valid value for `interactive_timeout`.

* `maximum` - The maximum valid value of `interactive_timeout`.

* `minimum` - The minimum valid value of `interactive_timeout`.

* `requires_restart` - Whether changing the value `interactive_timeout` requires the DB to restart.

* `type` - The type of the value of `interactive_timeout`.

## internal_tmp_mem_storage_engine

The following arguments are supported in the `internal_tmp_mem_storage_engine` specification block:

* `description` - The description of `internal_tmp_mem_storage_engine`.

* `enum` - A list of valid enum values for `internal_tmp_mem_storage_engine`.

* `example` - An example of a valid value for `internal_tmp_mem_storage_engine`.

* `requires_restart` - Whether changing the value `internal_tmp_mem_storage_engine` requires the DB to restart.

* `type` - The type of the value of `internal_tmp_mem_storage_engine`.

## max_allowed_packet

The following arguments are supported in the `max_allowed_packet` specification block:

* `description` - The description of `max_allowed_packet`.

* `example` - An example of a valid value for `max_allowed_packet`.

* `maximum` - The maximum valid value of `max_allowed_packet`.

* `minimum` - The minimum valid value of `max_allowed_packet`.

* `requires_restart` - Whether changing the value `max_allowed_packet` requires the DB to restart.

* `type` - The type of the value of `max_allowed_packet`.

## max_heap_table_size

The following arguments are supported in the `max_heap_table_size` specification block:

* `description` - The description of `max_heap_table_size`.

* `example` - An example of a valid value for `max_heap_table_size`.

* `maximum` - The maximum valid value of `max_heap_table_size`.

* `minimum` - The minimum valid value of `max_heap_table_size`.

* `requires_restart` - Whether changing the value `max_heap_table_size` requires the DB to restart.

* `type` - The type of the value of `max_heap_table_size`.

## net_buffer_length

The following arguments are supported in the `net_buffer_length` specification block:

* `description` - The description of `net_buffer_length`.

* `example` - An example of a valid value for `net_buffer_length`.

* `maximum` - The maximum valid value of `net_buffer_length`.

* `minimum` - The minimum valid value of `net_buffer_length`.

* `requires_restart` - Whether changing the value `net_buffer_length` requires the DB to restart.

* `type` - The type of the value of `net_buffer_length`.

## net_read_timeout

The following arguments are supported in the `net_read_timeout` specification block:

* `description` - The description of `net_read_timeout`.

* `example` - An example of a valid value for `net_read_timeout`.

* `maximum` - The maximum valid value of `net_read_timeout`.

* `minimum` - The minimum valid value of `net_read_timeout`.

* `requires_restart` - Whether changing the value `net_read_timeout` requires the DB to restart.

* `type` - The type of the value of `net_read_timeout`.

## net_write_timeout

The following arguments are supported in the `net_write_timeout` specification block:

* `description` - The description of `net_write_timeout`.

* `example` - An example of a valid value for `net_write_timeout`.

* `maximum` - The maximum valid value of `net_write_timeout`.

* `minimum` - The minimum valid value of `net_write_timeout`.

* `requires_restart` - Whether changing the value `net_write_timeout` requires the DB to restart.

* `type` - The type of the value of `net_write_timeout`.

## sort_buffer_size

The following arguments are supported in the `sort_buffer_size` specification block:

* `description` - The description of `sort_buffer_size`.

* `example` - An example of a valid value for `sort_buffer_size`.

* `maximum` - The maximum valid value of `sort_buffer_size`.

* `minimum` - The minimum valid value of `sort_buffer_size`.

* `requires_restart` - Whether changing the value `sort_buffer_size` requires the DB to restart.

* `type` - The type of the value of `sort_buffer_size`.

## sql_mode

The following arguments are supported in the `sql_mode` specification block:

* `description` - The description of `sql_mode`.

* `example` - An example of a valid value for `sql_mode`.

* `maxLength` - The maximum valid length of `sql_mode`.

* `pattern` - The pattern to match for `sql_mode`.

* `requires_restart` - Whether changing the value `sql_mode` requires the DB to restart.

* `type` - The type of the value of `sql_mode`.

## sql_require_primary_key

The following arguments are supported in the `sql_require_primary_key` specification block:

* `description` - The description of `sql_require_primary_key`.

* `example` - An example of a valid value for `sql_require_primary_key`.

* `requires_restart` - Whether changing the value `sql_require_primary_key` requires the DB to restart.

* `type` - The type of the value of `sql_require_primary_key`.

## tmp_table_size

The following arguments are supported in the `tmp_table_size` specification block:

* `description` - The description of `tmp_table_size`.

* `example` - An example of a valid value for `tmp_table_size`.

* `maximum` - The maximum valid value of `tmp_table_size`.

* `minimum` - The minimum valid value of `tmp_table_size`.

* `requires_restart` - Whether changing the value `tmp_table_size` requires the DB to restart.

* `type` - The type of the value of `tmp_table_size`.

## wait_timeout

The following arguments are supported in the `wait_timeout` specification block:

* `description` - The description of `wait_timeout`.

* `example` - An example of a valid value for `wait_timeout`.

* `maximum` - The maximum valid value of `wait_timeout`.

* `minimum` - The minimum valid value of `wait_timeout`.

* `requires_restart` - Whether changing the value `wait_timeout` requires the DB to restart.

* `type` - The type of the value of `wait_timeout`.
