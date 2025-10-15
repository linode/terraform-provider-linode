---
page_title: "Linode: linode_database_postgresql_v2"
description: |-
  Manages a Linode PostgreSQL Database.
---

# linode\_database\_postgresql\_v2

Provides a Linode PostgreSQL Database resource. This can be used to create, modify, and delete Linode PostgreSQL Databases.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-databases-postgre-sql-instances).

Please keep in mind that Managed Databases can take up to half an hour to provision.

## Example Usage

Creating a simple PostgreSQL database that does not allow connections:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"
}
```

Creating a simple PostgreSQL database that allows connections from all IPv4 addresses:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"
  
  allow_list = ["0.0.0.0/0"]
}
```

Creating a complex PostgreSQL database:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"

  allow_list = ["10.0.0.3/32"]
  cluster_size = 3

  updates = {
    duration = 4
    frequency = "weekly"
    hour_of_day = 22
    day_of_week = 2
  }
}
```

Creating a PostgreSQL database with engine config fields specified:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"
  
  engine_config_pg_autovacuum_analyze_scale_factor          = 0.1
  engine_config_pg_autovacuum_analyze_threshold             = 50
  engine_config_pg_autovacuum_max_workers                   = 3
  engine_config_pg_autovacuum_naptime                       = 100
  engine_config_pg_autovacuum_vacuum_cost_delay             = 20
  engine_config_pg_autovacuum_vacuum_cost_limit             = 200
  engine_config_pg_autovacuum_vacuum_scale_factor           = 0.2
  engine_config_pg_autovacuum_vacuum_threshold              = 100
  engine_config_pg_bgwriter_delay                           = 1000
  engine_config_pg_bgwriter_flush_after                     = 512
  engine_config_pg_bgwriter_lru_maxpages                    = 100
  engine_config_pg_bgwriter_lru_multiplier                  = 2.0
  engine_config_pg_deadlock_timeout                         = 1000
  engine_config_pg_default_toast_compression                = "pglz"
  engine_config_pg_idle_in_transaction_session_timeout      = 60000
  engine_config_pg_jit                                      = true
  engine_config_pg_max_files_per_process                    = 1000
  engine_config_pg_max_locks_per_transaction                = 64
  engine_config_pg_max_logical_replication_workers          = 4
  engine_config_pg_max_parallel_workers                     = 8
  engine_config_pg_max_parallel_workers_per_gather          = 2
  engine_config_pg_max_pred_locks_per_transaction           = 128
  engine_config_pg_max_replication_slots                    = 8
  engine_config_pg_max_slot_wal_keep_size                   = 128
  engine_config_pg_max_stack_depth                          = 2097152
  engine_config_pg_max_standby_archive_delay                = 60000
  engine_config_pg_max_standby_streaming_delay              = 60000
  engine_config_pg_max_wal_senders                          = 20
  engine_config_pg_max_worker_processes                     = 8
  engine_config_pg_password_encryption                      = "scram-sha-256"
  engine_config_pg_pg_partman_bgw_interval                  = 3600
  engine_config_pg_pg_partman_bgw_role                      = "myrolename"
  engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan   = true
  engine_config_pg_pg_stat_monitor_pgsm_max_buckets         = 5
  engine_config_pg_pg_stat_statements_track                 = "all"
  engine_config_pg_temp_file_limit                          = 100
  engine_config_pg_timezone                                 = "Europe/Helsinki"
  engine_config_pg_track_activity_query_size                = 2048
  engine_config_pg_track_commit_timestamp                   = "on"
  engine_config_pg_track_functions                          = "all"
  engine_config_pg_track_io_timing                          = "on"
  engine_config_pg_wal_sender_timeout                       = 60000
  engine_config_pg_wal_writer_delay                         = 200
  engine_config_pg_stat_monitor_enable                      = true
  engine_config_pglookout_max_failover_replication_time_lag = 10000
  engine_config_shared_buffers_percentage                   = 25.0
  engine_config_work_mem                                    = 400
}
```

Creating a forked PostgreSQL database:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
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

* `engine_id` - (Required) The Managed Database engine in engine/version format. (e.g. `postgresql/16`)

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

* `engine_config_pg_autovacuum_analyze_scale_factor` - (Optional) Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)

* `engine_config_pg_autovacuum_analyze_threshold` - (Optional) Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.

* `engine_config_pg_autovacuum_max_workers` - (Optional) Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.

* `engine_config_pg_autovacuum_naptime` - (Optional) Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute

* `engine_config_pg_autovacuum_vacuum_cost_delay` - (Optional) Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds

* `engine_config_pg_autovacuum_vacuum_cost_limit` - (Optional) Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.

* `engine_config_pg_autovacuum_vacuum_scale_factor` - (Optional) Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)

* `engine_config_pg_autovacuum_vacuum_threshold` - (Optional) Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples.

* `engine_config_pg_bgwriter_delay` - (Optional) Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.

* `engine_config_pg_bgwriter_flush_after` - (Optional) Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.

* `engine_config_pg_bgwriter_lru_maxpages` - (Optional) In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.

* `engine_config_pg_bgwriter_lru_multiplier` - (Optional) The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.

* `engine_config_pg_deadlock_timeout` - (Optional) This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.

* `engine_config_pg_default_toast_compression` - (Optional) Specifies the default TOAST compression method for values of compressible columns (the default is lz4).

* `engine_config_pg_idle_in_transaction_session_timeout` - (Optional) Time out sessions with open transactions after this number of milliseconds.

* `engine_config_pg_jit` - (Optional) Controls system-wide use of Just-in-Time Compilation (JIT).

* `engine_config_pg_max_files_per_process` - (Optional) PostgreSQL maximum number of files that can be open per process.

* `engine_config_pg_max_locks_per_transaction` - (Optional) PostgreSQL maximum locks per transaction.

* `engine_config_pg_max_logical_replication_workers` - (Optional) PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers).

* `engine_config_pg_max_parallel_workers` - (Optional) Sets the maximum number of workers that the system can support for parallel queries.

* `engine_config_pg_max_parallel_workers_per_gather` - (Optional) Sets the maximum number of workers that can be started by a single Gather or Gather Merge node.

* `engine_config_pg_max_pred_locks_per_transaction` - (Optional) PostgreSQL maximum predicate locks per transaction.

* `engine_config_pg_max_replication_slots` - (Optional) PostgreSQL maximum replication slots.

* `engine_config_pg_max_slot_wal_keep_size` - (Optional) PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.

* `engine_config_pg_max_stack_depth` - (Optional) Maximum depth of the stack in bytes.

* `engine_config_pg_max_standby_archive_delay` - (Optional) Max standby archive delay in milliseconds.

* `engine_config_pg_max_standby_streaming_delay` - (Optional) Max standby streaming delay in milliseconds.

* `engine_config_pg_max_wal_senders` - (Optional) PostgreSQL maximum WAL senders.

* `engine_config_pg_max_worker_processes` - (Optional) Sets the maximum number of background processes that the system can support.

* `engine_config_pg_password_encryption` - (Optional) Chooses the algorithm for encrypting passwords. (default `md5`)

* `engine_config_pg_pg_partman_bgw_interval` - (Optional) Sets the time interval to run pg_partman's scheduled tasks.

* `engine_config_pg_pg_partman_bgw_role` - (Optional) Controls which role to use for pg_partman's scheduled background tasks.

* `engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan` - (Optional) Enables or disables query plan monitoring.

* `engine_config_pg_pg_stat_monitor_pgsm_max_buckets` - (Optional) Sets the maximum number of buckets.

* `engine_config_pg_pg_stat_statements_track` - (Optional) Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.

* `engine_config_pg_temp_file_limit` - (Optional) PostgreSQL temporary file limit in KiB, -1 for unlimited.

* `engine_config_pg_timezone` - (Optional) PostgreSQL service timezone.

* `engine_config_pg_track_activity_query_size` - (Optional) Specifies the number of bytes reserved to track the currently executing command for each active session.

* `engine_config_pg_track_commit_timestamp` - (Optional) Record commit time of transactions.

* `engine_config_pg_track_functions` - (Optional) Enables tracking of function call counts and time used.

* `engine_config_pg_track_io_timing` - (Optional) Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.

* `engine_config_pg_wal_sender_timeout` - (Optional) Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout.

* `engine_config_pg_wal_writer_delay` - (Optional) WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance.

* `engine_config_pg_stat_monitor_enable` - (Optional) Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted. When this extension is enabled, pg_stat_statements results for utility commands are unreliable. (default `false`)

* `engine_config_pglookout_max_failover_replication_time_lag` - (Optional) Number of seconds of master unavailability before triggering database failover to standby. (default `60`)

* `engine_config_shared_buffers_percentage` - (Optional) Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.

* `engine_config_work_mem` - (Optional) Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Managed Database.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database.

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Databases is encrypted.

* `engine` - The Managed Database engine. (e.g. `postgresql`)

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

Linode PostgreSQL Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_postgresql_v2.foobar 1234567
```
