---
page_title: "Linode: linode_database_mysql_v2"
description: |-
  Manages a Linode MySQL Database.
---

# linode\_database\_mysql\_v2

Provides a Linode MySQL Database resource. This can be used to create, modify, and delete Linode MySQL Databases.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-databases-mysql-instances).

Please keep in mind that Managed Databases can take up to half an hour to provision.

## Example Usage

Creating a simple MySQL database that does not allow connections:

```hcl
resource "linode_database_mysql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "mysql/8"
  region = "us-mia"
  type = "g6-nanode-1"
}
```

Creating a simple MySQL database that allows connections from all IPv4 addresses:

```hcl
resource "linode_database_mysql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "mysql/8"
  region = "us-mia"
  type = "g6-nanode-1"
  
  allow_list = ["0.0.0.0/0"]
}
```

Creating a complex MySQL database:

```hcl
resource "linode_database_mysql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "mysql/8"
  region = "us-mia"
  type = "g6-nanode-1"

  allow_list = ["10.0.0.3/32"]
  cluster_size = 3

  updates = {
    duration = 4
    frequency = "weekly"
    hour_of_day = 22
    day_of_week = 3
  }
}
```

Creating a MySQL database with engine config fields specified:

```hcl
resource "linode_database_mysql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "mysql/8"
  region = "us-mia"
  type = "g6-nanode-1"

  engine_config_binlog_retention_period                    = 3600
  engine_config_mysql_connect_timeout                      = 10
  engine_config_mysql_default_time_zone                    = "+00:00"
  engine_config_mysql_group_concat_max_len                 = 4096.0
  engine_config_mysql_information_schema_stats_expiry      = 3600
  engine_config_mysql_innodb_change_buffer_max_size        = 25
  engine_config_mysql_innodb_flush_neighbors               = 0
  engine_config_mysql_innodb_ft_min_token_size             = 7
  engine_config_mysql_innodb_ft_server_stopword_table      = "mysql/innodb_ft_default_stopword"
  engine_config_mysql_innodb_lock_wait_timeout             = 300
  engine_config_mysql_innodb_log_buffer_size               = 16777216
  engine_config_mysql_innodb_online_alter_log_max_size     = 268435456
  engine_config_mysql_innodb_read_io_threads               = 4
  engine_config_mysql_innodb_rollback_on_timeout           = true
  engine_config_mysql_innodb_thread_concurrency            = 8
  engine_config_mysql_innodb_write_io_threads              = 4
  engine_config_mysql_interactive_timeout                  = 300
  engine_config_mysql_internal_tmp_mem_storage_engine      = "TempTable"
  engine_config_mysql_max_allowed_packet                   = 67108864
  engine_config_mysql_max_heap_table_size                  = 16777216
  engine_config_mysql_net_buffer_length                    = 16384
  engine_config_mysql_net_read_timeout                     = 30
  engine_config_mysql_net_write_timeout                    = 30
  engine_config_mysql_sort_buffer_size                     = 262144
  engine_config_mysql_sql_mode                             = "TRADITIONAL,ANSI"
  engine_config_mysql_sql_require_primary_key              = false
  engine_config_mysql_tmp_table_size                       = 16777216
  engine_config_mysql_wait_timeout                         = 28800
}
```

Creating a forked MySQL database:

```hcl
resource "linode_database_mysql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "mysql/8"
  region = "us-mia"
  type = "g6-nanode-1"

  fork_source = 12345
}
```

Creating a PostgreSQL database hidden behind a VPC:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"

  private_network = {
    vpc_id = 123
    subnet_id = 456
    public_access = false
  }
}
```

> **_NOTE:_** The name of the default database in the returned database cluster is `defaultdb`.

## Argument Reference

The following arguments are supported:

* `engine_id` - (Required) The Managed Database engine in engine/version format. (e.g. `mysql`)

* `label` - (Required) A unique, user-defined string referring to the Managed Database.

* `region` - (Required) The region to use for the Managed Database.

* `type` - (Required) The Linode Instance type used for the nodes of the Managed Database.

- - -

* `allow_list` - (Optional) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format. Use `linode_database_access_controls` to manage your allow list separately.

* `cluster_size` - (Optional) The number of Linode Instance nodes deployed to the Managed Database. (default `1`)

* `fork_restore_time` - (Optional) The database timestamp from which it was restored.

* `fork_source` - (Optional) The ID of the database that was forked from.

* [`private_network`](#private_network) - (Optional) Restricts access to this database using a virtual private cloud (VPC) that you've configured in the region where the database will live.

* [`updates`](#updates) - (Optional) Configuration settings for automated patch update maintenance for the Managed Database.

* `engine_config_binlog_retention_period` - (Optional) The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default, for example if using the MySQL Debezium Kafka connector.

* `engine_config_mysql_connect_timeout` - (Optional) The number of seconds that the mysqld server waits for a connect packet before responding with "Bad handshake".

* `engine_config_mysql_default_time_zone` - (Optional) Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or `SYSTEM` to use the MySQL server default.

* `engine_config_mysql_group_concat_max_len` - (Optional) The maximum permitted result length in bytes for the `GROUP_CONCAT()` function.

* `engine_config_mysql_information_schema_stats_expiry` - (Optional) The time, in seconds, before cached statistics expire.

* `engine_config_mysql_innodb_change_buffer_max_size` - (Optional) Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25.

* `engine_config_mysql_innodb_flush_neighbors` - (Optional) Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent.

* `engine_config_mysql_innodb_ft_min_token_size` - (Optional) Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_ft_server_stopword_table` - (Optional) This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables. This field is nullable.

* `engine_config_mysql_innodb_lock_wait_timeout` - (Optional) The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.

* `engine_config_mysql_innodb_log_buffer_size` - (Optional) The size in bytes of the buffer that InnoDB uses to write to the log files on disk.

* `engine_config_mysql_innodb_online_alter_log_max_size` - (Optional) The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.

* `engine_config_mysql_innodb_read_io_threads` - (Optional) The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_rollback_on_timeout` - (Optional) When enabled, a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_innodb_thread_concurrency` - (Optional) Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit).

* `engine_config_mysql_innodb_write_io_threads` - (Optional) The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_interactive_timeout` - (Optional) The number of seconds the server waits for activity on an interactive connection before closing it.

* `engine_config_mysql_internal_tmp_mem_storage_engine` - (Optional) The storage engine for in-memory internal temporary tables.

* `engine_config_mysql_max_allowed_packet` - (Optional) Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M).

* `engine_config_mysql_max_heap_table_size` - (Optional) Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M).

* `engine_config_mysql_net_buffer_length` - (Optional) Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.

* `engine_config_mysql_net_read_timeout` - (Optional) The number of seconds to wait for more data from a connection before aborting the read.

* `engine_config_mysql_net_write_timeout` - (Optional) The number of seconds to wait for a block to be written to a connection before aborting the write.

* `engine_config_mysql_sort_buffer_size` - (Optional) Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K).

* `engine_config_mysql_sql_mode` - (Optional) Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned. (default `ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES`)

* `engine_config_mysql_sql_require_primary_key` - (Optional) Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them. (default `true`)

* `engine_config_mysql_tmp_table_size` - (Optional) Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M).

* `engine_config_mysql_wait_timeout` - (Optional) The number of seconds the server waits for activity on a noninteractive connection before closing it.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Managed Database.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database.

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Databases is encrypted.

* `engine` - The Managed Database engine. (e.g. `mysql`)

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private host for the managed database.

* `pending_updates` - A set of pending updates.

* `platform` - The back-end platform for relational databases used by the service.

* `port` - The access port for this Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `ssl_connection` - Whether to require SSL credentials to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `suspended` - Whether this Managed Database should be suspended.

* `updated` - When this Managed Database was last updated.

* `version` - The Managed Database engine version. (e.g. `13.2`)

## pending_updates

The following arguments are exposed by each entry in the `pending_updates` attribute:

* `deadline` - The time when a mandatory update needs to be applied.

* `description` - A description of the update.

* `planned_for` - The date and time a maintenance update will be applied.

## updates

The following arguments are supported in the `updates` specification block:

* `day_of_week` - (Required) The day to perform maintenance. (`monday`, `tuesday`, ...)

* `duration` - (Required) The maximum maintenance window time in hours. (`1`..`3`)

* `frequency` - (Required) The frequency at which maintenance occurs. (`weekly`)

* `hour_of_day` - (Required) The hour to begin maintenance based in UTC time. (`0`..`23`)

## private_network

The following arguments are supported in the `private_network` specification block:

* `vpc_id` - (Required) The ID of the virtual private cloud (VPC) to restrict access to this database using.

* `subnet_id` - (Required) The ID of the VPC subnet to restrict access to this database using.

* `public_access` - (Optional) Set to `true` to allow clients outside the VPC to connect to the database using a public IP address. (Default `false`)

## Import

Linode MySQL Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_mysql_v2.foobar 1234567
```
