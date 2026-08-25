---
page_title: "Linode: linode_database_postgresql_config"
description: |-
  Provides information about a Linode PostgreSQL Database's Configuration Options.
---

# Data Source: linode\_database\_postgresql\_config

Provides information about a Linode PostgreSQL Database's Configuration Options.
For more information, see the [Linode APIv4 docs](TODO).

## Example Usage

Get information about a PostgreSQL database's configuration options:

```hcl
data "linode_database_postgresql_config" "my-db-config" {}
```

## Attributes Reference

The `linode_database_postgresql_config` data source exports the following attributes:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`pg_stat_monitor_enable`](#pg_stat_monitor_enable) - (Read-Only Object List) Enable the pg_stat_monitor extension. Enabling this extension will cause the cluster to be restarted.When this extension is enabled, pg_stat_statements results for utility commands are unreliable. Referenced with an index (e.g. `pg_stat_monitor_enable.0.description`).

* [`pglookout`](#pglookout) - (Read-Only Object List) System-wide settings for pglookout. Referenced with an index (e.g. `pglookout.0.max_failover_replication_time_lag`).

* [`shared_buffers_percentage`](#shared_buffers_percentage) - (Read-Only Object List) Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value. Referenced with an index (e.g. `shared_buffers_percentage.0.description`).

* [`work_mem`](#work_mem) - (Read-Only Object List) Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB). Referenced with an index (e.g. `work_mem.0.description`).

* [`pg`](#pg) - (Read-Only Object List) PostgreSQL configuration settings. Referenced with an index (e.g. `pg.0.autovacuum_analyze_scale_factor`).

## pg_stat_monitor_enable

The following arguments are supported in the `pg_stat_monitor_enable` specification block:

* `description` - The description of `pg_stat_monitor_enable`.

* `requires_restart` - Whether changing the value `pg_stat_monitor_enable` requires the DB to restart.

* `type` - The type of the value of `pg_stat_monitor_enable`.

## pglookout

The following arguments are supported in the `pglookout` specification block:

* [`max_failover_replication_time_lag`](#max_failover_replication_time_lag) - (Read-Only Object) The maximum failover replication time lag for `pglookout`. Referenced directly (e.g. `max_failover_replication_time_lag.description`).

## max_failover_replication_time_lag

The following arguments are supported in the `max_failover_replication_time_lag` specification block:

* `description` - The description of `max_failover_replication_time_lag`.

* `maximum` - The maximum valid value for `max_failover_replication_time_lag`.

* `minimum` - The minimum valid value for `max_failover_replication_time_lag`.

* `requires_restart` - Whether changing the value of `max_failover_replication_time_lag` requires the DB to restart.

* `type` - The type of the value of `max_failover_replication_time_lag`.

## shared_buffers_percentage

The following arguments are supported in the `shared_buffers_percentage` specification block:

* `description` - The description of `shared_buffers_percentage`.

* `example` - An example of a valid value for `shared_buffers_percentage`.

* `maximum` - The maximum valid value for `shared_buffers_percentage`.

* `minimum` - The minimum valid value for `shared_buffers_percentage`.

* `requires_restart` - Whether changing the value of `shared_buffers_percentage` requires the DB to restart.

* `type` - The type of the value of `shared_buffers_percentage`.

## work_mem

The following arguments are supported in the `work_mem` specification block:

* `description` - The description of `work_mem`.

* `example` - An example of a valid value for `work_mem`.

* `maximum` - The maximum valid value for `work_mem`.

* `minimum` - The minimum valid value for `work_mem`.

* `requires_restart` - Whether changing the value of `work_mem` requires the DB to restart.

* `type` - The type of the value of `work_mem`.

## pg

The following arguments are supported in the `pg` specification block:

* [`autovacuum_analyze_scale_factor`](#autovacuum_analyze_scale_factor) - (Read-Only Object) Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size) Referenced directly (e.g. `autovacuum_analyze_scale_factor.description`).

* [`autovacuum_analyze_threshold`](#autovacuum_analyze_threshold) - (Read-Only Object) Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples. Referenced directly (e.g. `autovacuum_analyze_threshold.description`).

* [`autovacuum_max_workers`](#autovacuum_max_workers) - (Read-Only Object) Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start. Referenced directly (e.g. `autovacuum_max_workers.description`).

* [`autovacuum_naptime`](#autovacuum_naptime) - (Read-Only Object) Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute Referenced directly (e.g. `autovacuum_naptime.description`).

* [`autovacuum_vacuum_cost_delay`](#autovacuum_vacuum_cost_delay) - (Read-Only Object) Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds Referenced directly (e.g. `autovacuum_vacuum_cost_delay.description`).

* [`autovacuum_vacuum_cost_limit`](#autovacuum_vacuum_cost_limit) - (Read-Only Object) Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used. Referenced directly (e.g. `autovacuum_vacuum_cost_limit.description`).

* [`autovacuum_vacuum_scale_factor`](#autovacuum_vacuum_scale_factor) - (Read-Only Object) Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size) Referenced directly (e.g. `autovacuum_vacuum_scale_factor.description`).

* [`autovacuum_vacuum_threshold`](#autovacuum_vacuum_threshold) - (Read-Only Object) Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples. Referenced directly (e.g. `autovacuum_vacuum_threshold.description`).

* [`bgwriter_delay`](#bgwriter_delay) - (Read-Only Object) Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200. Referenced directly (e.g. `bgwriter_delay.description`).

* [`bgwriter_flush_after`](#bgwriter_flush_after) - (Read-Only Object) Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback. Referenced directly (e.g. `bgwriter_flush_after.description`).

* [`bgwriter_lru_maxpages`](#bgwriter_lru_maxpages) - (Read-Only Object) In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100. Referenced directly (e.g. `bgwriter_lru_maxpages.description`).

* [`bgwriter_lru_multiplier`](#bgwriter_lru_multiplier) - (Read-Only Object) The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0. Referenced directly (e.g. `bgwriter_lru_multiplier.description`).

* [`deadlock_timeout`](#deadlock_timeout) - (Read-Only Object) This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition. Referenced directly (e.g. `deadlock_timeout.description`).

* [`default_toast_compression`](#default_toast_compression) - (Read-Only Object) Specifies the default TOAST compression method for values of compressible columns (the default is lz4). Referenced directly (e.g. `default_toast_compression.description`).

* [`idle_in_transaction_session_timeout`](#idle_in_transaction_session_timeout) - (Read-Only Object) Time out sessions with open transactions after this number of milliseconds. Referenced directly (e.g. `idle_in_transaction_session_timeout.description`).

* [`jit`](#jit) - (Read-Only Object) Controls system-wide use of Just-in-Time Compilation (JIT). Referenced directly (e.g. `jit.description`).

* [`max_files_per_process`](#max_files_per_process) - (Read-Only Object) PostgreSQL maximum number of files that can be open per process. Referenced directly (e.g. `max_files_per_process.description`).

* [`max_locks_per_transaction`](#max_locks_per_transaction) - (Read-Only Object) PostgreSQL maximum locks per transaction. Referenced directly (e.g. `max_locks_per_transaction.description`).

* [`max_logical_replication_workers`](#max_logical_replication_workers) - (Read-Only Object) PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers). Referenced directly (e.g. `max_logical_replication_workers.description`).

* [`max_parallel_workers`](#max_parallel_workers) - (Read-Only Object) Sets the maximum number of workers that the system can support for parallel queries. Referenced directly (e.g. `max_parallel_workers.description`).

* [`max_parallel_workers_per_gather`](#max_parallel_workers_per_gather) - (Read-Only Object) Sets the maximum number of workers that can be started by a single Gather or Gather Merge node. Referenced directly (e.g. `max_parallel_workers_per_gather.description`).

* [`max_pred_locks_per_transaction`](#max_pred_locks_per_transaction) - (Read-Only Object) PostgreSQL maximum predicate locks per transaction. Referenced directly (e.g. `max_pred_locks_per_transaction.description`).

* [`max_replication_slots`](#max_replication_slots) - (Read-Only Object) PostgreSQL maximum replication slots. Referenced directly (e.g. `max_replication_slots.description`).

* [`max_slot_wal_keep_size`](#max_slot_wal_keep_size) - (Read-Only Object) PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this. Referenced directly (e.g. `max_slot_wal_keep_size.description`).

* [`max_stack_depth`](#max_stack_depth) - (Read-Only Object) Maximum depth of the stack in bytes. Referenced directly (e.g. `max_stack_depth.description`).

* [`max_standby_archive_delay`](#max_standby_archive_delay) - (Read-Only Object) Max standby archive delay in milliseconds. Referenced directly (e.g. `max_standby_archive_delay.description`).

* [`max_standby_streaming_delay`](#max_standby_streaming_delay) - (Read-Only Object) Max standby streaming delay in milliseconds. Referenced directly (e.g. `max_standby_streaming_delay.description`).

* [`max_wal_senders`](#max_wal_senders) - (Read-Only Object) PostgreSQL maximum WAL senders. Referenced directly (e.g. `max_wal_senders.description`).

* [`max_worker_processes`](#max_worker_processes) - (Read-Only Object) Sets the maximum number of background processes that the system can support. Referenced directly (e.g. `max_worker_processes.description`).

* [`password_encryption`](#password_encryption) - (Read-Only Object) Chooses the algorithm for encrypting passwords. Referenced directly (e.g. `password_encryption.description`).

* [`pg_partman_bgw.interval`](#pg_partman_bgw_interval) - (Read-Only Object) Sets the time interval to run pg_partman's scheduled tasks. Referenced directly (e.g. `pg_partman_bgw.interval.description`).

* [`pg_partman_bgw.role`](#pg_partman_bgw_role) - (Read-Only Object) Controls which role to use for pg_partman's scheduled background tasks. Referenced directly (e.g. `pg_partman_bgw.role.description`).

* [`pg_stat_monitor.pgsm_enable_query_plan`](#pg_stat_monitor_pgsm_enable_query_plan) - (Read-Only Object) Enables or disables query plan monitoring. Referenced directly (e.g. `pg_stat_monitor.pgsm_enable_query_plan.description`).

* [`pg_stat_monitor.pgsm_max_buckets`](#pg_stat_monitor_pgsm_max_buckets) - (Read-Only Object) Sets the maximum number of buckets. Referenced directly (e.g. `pg_stat_monitor.pgsm_max_buckets.description`).

* [`pg_stat_statements.track`](#pg_stat_statements_track) - (Read-Only Object) Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top. Referenced directly (e.g. `pg_stat_statements.track.description`).

* [`temp_file_limit`](#temp_file_limit) - (Read-Only Object) PostgreSQL temporary file limit in KiB, -1 for unlimited. Referenced directly (e.g. `temp_file_limit.description`).

* [`timezone`](#timezone) - (Read-Only Object) PostgreSQL service timezone. Referenced directly (e.g. `timezone.description`).

* [`track_activity_query_size`](#track_activity_query_size) - (Read-Only Object) Specifies the number of bytes reserved to track the currently executing command for each active session. Referenced directly (e.g. `track_activity_query_size.description`).

* [`track_commit_timestamp`](#track_commit_timestamp) - (Read-Only Object) Record commit time of transactions. Referenced directly (e.g. `track_commit_timestamp.description`).

* [`track_functions`](#track_functions) - (Read-Only Object) Enables tracking of function call counts and time used. Referenced directly (e.g. `track_functions.description`).

* [`track_io_timing`](#track_io_timing) - (Read-Only Object) Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms. Referenced directly (e.g. `track_io_timing.description`).

* [`wal_sender_timeout`](#wal_sender_timeout) - (Read-Only Object) Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout. Referenced directly (e.g. `wal_sender_timeout.description`).

* [`wal_writer_delay`](#wal_writer_delay) - (Read-Only Object) WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance. Referenced directly (e.g. `wal_writer_delay.description`).

## autovacuum_analyze_scale_factor

The following arguments are supported in the `autovacuum_analyze_scale_factor` specification block:

* `description` - The description of `autovacuum_analyze_scale_factor`.

* `maximum` - The maximum valid value for `autovacuum_analyze_scale_factor`.

* `minimum` - The minimum valid value for `autovacuum_analyze_scale_factor`.

* `requires_restart` - Whether changing the value of `autovacuum_analyze_scale_factor` requires the DB to restart.

* `type` - The type of the value of `autovacuum_analyze_scale_factor`.

## autovacuum_analyze_threshold

The following arguments are supported in the `autovacuum_analyze_threshold` specification block:

* `description` - The description of `autovacuum_analyze_threshold`.

* `maximum` - The maximum valid value for `autovacuum_analyze_threshold`.

* `minimum` - The minimum valid value for `autovacuum_analyze_threshold`.

* `requires_restart` - Whether changing the value of `autovacuum_analyze_threshold` requires the DB to restart.

* `type` - The type of the value of `autovacuum_analyze_threshold`.

## autovacuum_max_workers

The following arguments are supported in the `autovacuum_max_workers` specification block:

* `description` - The description of `autovacuum_max_workers`.

* `maximum` - The maximum valid value for `autovacuum_max_workers`.

* `minimum` - The minimum valid value for `autovacuum_max_workers`.

* `requires_restart` - Whether changing the value of `autovacuum_max_workers` requires the DB to restart.

* `type` - The type of the value of `autovacuum_max_workers`.

## autovacuum_naptime

The following arguments are supported in the `autovacuum_naptime` specification block:

* `description` - The description of `autovacuum_naptime`.

* `maximum` - The maximum valid value for `autovacuum_naptime`.

* `minimum` - The minimum valid value for `autovacuum_naptime`.

* `requires_restart` - Whether changing the value of `autovacuum_naptime` requires the DB to restart.

* `type` - The type of the value of `autovacuum_naptime`.

## autovacuum_vacuum_cost_delay

The following arguments are supported in the `autovacuum_vacuum_cost_delay` specification block:

* `description` - The description of `autovacuum_vacuum_cost_delay`.

* `maximum` - The maximum valid value for `autovacuum_vacuum_cost_delay`.

* `minimum` - The minimum valid value for `autovacuum_vacuum_cost_delay`.

* `requires_restart` - Whether changing the value of `autovacuum_vacuum_cost_delay` requires the DB to restart.

* `type` - The type of the value of `autovacuum_vacuum_cost_delay`.

## autovacuum_vacuum_cost_limit

The following arguments are supported in the `autovacuum_vacuum_cost_limit` specification block:

* `description` - The description of `autovacuum_vacuum_cost_limit`.

* `maximum` - The maximum valid value for `autovacuum_vacuum_cost_limit`.

* `minimum` - The minimum valid value for `autovacuum_vacuum_cost_limit`.

* `requires_restart` - Whether changing the value of `autovacuum_vacuum_cost_limit` requires the DB to restart.

* `type` - The type of the value of `autovacuum_vacuum_cost_limit`.

## autovacuum_vacuum_scale_factor

The following arguments are supported in the `autovacuum_vacuum_scale_factor` specification block:

* `description` - The description of `autovacuum_vacuum_scale_factor`.

* `maximum` - The maximum valid value for `autovacuum_vacuum_scale_factor`.

* `minimum` - The minimum valid value for `autovacuum_vacuum_scale_factor`.

* `requires_restart` - Whether changing the value of `autovacuum_vacuum_scale_factor` requires the DB to restart.

* `type` - The type of the value of `autovacuum_vacuum_scale_factor`.

## autovacuum_vacuum_threshold

The following arguments are supported in the `autovacuum_vacuum_threshold` specification block:

* `description` - The description of `autovacuum_vacuum_threshold`.

* `maximum` - The maximum valid value for `autovacuum_vacuum_threshold`.

* `minimum` - The minimum valid value for `autovacuum_vacuum_threshold`.

* `requires_restart` - Whether changing the value of `autovacuum_vacuum_threshold` requires the DB to restart.

* `type` - The type of the value of `autovacuum_vacuum_threshold`.

## bgwriter_delay

The following arguments are supported in the `bgwriter_delay` specification block:

* `description` - The description of `bgwriter_delay`.

* `example` - An example of a valid value for `bgwriter_delay`.

* `maximum` - The maximum valid value for `bgwriter_delay`.

* `minimum` - The minimum valid value for `bgwriter_delay`.

* `requires_restart` - Whether changing the value of `bgwriter_delay` requires the DB to restart.

* `type` - The type of the value of `bgwriter_delay`.

## bgwriter_flush_after

The following arguments are supported in the `bgwriter_flush_after` specification block:

* `description` - The description of `bgwriter_flush_after`.

* `example` - An example of a valid value for `bgwriter_flush_after`.

* `maximum` - The maximum valid value for `bgwriter_flush_after`.

* `minimum` - The minimum valid value for `bgwriter_flush_after`.

* `requires_restart` - Whether changing the value of `bgwriter_flush_after` requires the DB to restart.

* `type` - The type of the value of `bgwriter_flush_after`.

## bgwriter_lru_maxpages

The following arguments are supported in the `bgwriter_lru_maxpages` specification block:

* `description` - The description of `bgwriter_lru_maxpages`.

* `example` - An example of a valid value for `bgwriter_lru_maxpages`.

* `maximum` - The maximum valid value for `bgwriter_lru_maxpages`.

* `minimum` - The minimum valid value for `bgwriter_lru_maxpages`.

* `requires_restart` - Whether changing the value of `bgwriter_lru_maxpages` requires the DB to restart.

* `type` - The type of the value of `bgwriter_lru_maxpages`.

## bgwriter_lru_multiplier

The following arguments are supported in the `bgwriter_lru_multiplier` specification block:

* `description` - The description of `bgwriter_lru_multiplier`.

* `example` - An example of a valid value for `bgwriter_lru_multiplier`.

* `maximum` - The maximum valid value for `bgwriter_lru_multiplier`.

* `minimum` - The minimum valid value for `bgwriter_lru_multiplier`.

* `requires_restart` - Whether changing the value of `bgwriter_lru_multiplier` requires the DB to restart.

* `type` - The type of the value of `bgwriter_lru_multiplier`.

## deadlock_timeout

The following arguments are supported in the `deadlock_timeout` specification block:

* `description` - The description of `deadlock_timeout`.

* `example` - An example of a valid value for `deadlock_timeout`.

* `maximum` - The maximum valid value for `deadlock_timeout`.

* `minimum` - The minimum valid value for `deadlock_timeout`.

* `requires_restart` - Whether changing the value of `deadlock_timeout` requires the DB to restart.

* `type` - The type of the value of `deadlock_timeout`.

## default_toast_compression

The following arguments are supported in the `default_toast_compression` specification block:

* `description` - The description of `default_toast_compression`.

* `enum` - A list of valid compression methods for `default_toast_compression`.

* `example` - An example of a valid value for `default_toast_compression`.

* `requires_restart` - Whether changing the value of `default_toast_compression` requires the DB to restart.

* `type` - The type of the value of `default_toast_compression`.

## idle_in_transaction_session_timeout

The following arguments are supported in the `idle_in_transaction_session_timeout` specification block:

* `description` - The description of `idle_in_transaction_session_timeout`.

* `maximum` - The maximum valid value for `idle_in_transaction_session_timeout`.

* `minimum` - The minimum valid value for `idle_in_transaction_session_timeout`.

* `requires_restart` - Whether changing the value of `idle_in_transaction_session_timeout` requires the DB to restart.

* `type` - The type of the value of `idle_in_transaction_session_timeout`.

## jit

The following arguments are supported in the `jit` specification block:

* `description` - The description of `jit`.

* `example` - An example of a valid value for `jit`.

* `requires_restart` - Whether changing the value of `jit` requires the DB to restart.

* `type` - The type of the value of `jit`.

## max_files_per_process

The following arguments are supported in the `max_files_per_process` specification block:

* `description` - The description of `max_files_per_process`.

* `maximum` - The maximum valid value for `max_files_per_process`.

* `minimum` - The minimum valid value for `max_files_per_process`.

* `requires_restart` - Whether changing the value of `max_files_per_process` requires the DB to restart.

* `type` - The type of the value of `max_files_per_process`.

## max_locks_per_transaction

The following arguments are supported in the `max_locks_per_transaction` specification block:

* `description` - The description of `max_locks_per_transaction`.

* `maximum` - The maximum valid value for `max_locks_per_transaction`.

* `minimum` - The minimum valid value for `max_locks_per_transaction`.

* `requires_restart` - Whether changing the value of `max_locks_per_transaction` requires the DB to restart.

* `type` - The type of the value of `max_locks_per_transaction`.

## max_logical_replication_workers

The following arguments are supported in the `max_logical_replication_workers` specification block:

* `description` - The description of `max_logical_replication_workers`.

* `maximum` - The maximum valid value for `max_logical_replication_workers`.

* `minimum` - The minimum valid value for `max_logical_replication_workers`.

* `requires_restart` - Whether changing the value of `max_logical_replication_workers` requires the DB to restart.

* `type` - The type of the value of `max_logical_replication_workers`.

## max_parallel_workers

The following arguments are supported in the `max_parallel_workers` specification block:

* `description` - The description of `max_parallel_workers`.

* `maximum` - The maximum valid value for `max_parallel_workers`.

* `minimum` - The minimum valid value for `max_parallel_workers`.

* `requires_restart` - Whether changing the value of `max_parallel_workers` requires the DB to restart.

* `type` - The type of the value of `max_parallel_workers`.

## max_parallel_workers_per_gather

The following arguments are supported in the `max_parallel_workers_per_gather` specification block:

* `description` - The description of `max_parallel_workers_per_gather`.

* `maximum` - The maximum valid value for `max_parallel_workers_per_gather`.

* `minimum` - The minimum valid value for `max_parallel_workers_per_gather`.

* `requires_restart` - Whether changing the value of `max_parallel_workers_per_gather` requires the DB to restart.

* `type` - The type of the value of `max_parallel_workers_per_gather`.

## max_pred_locks_per_transaction

The following arguments are supported in the `max_pred_locks_per_transaction` specification block:

* `description` - The description of `max_pred_locks_per_transaction`.

* `maximum` - The maximum valid value for `max_pred_locks_per_transaction`.

* `minimum` - The minimum valid value for `max_pred_locks_per_transaction`.

* `requires_restart` - Whether changing the value of `max_pred_locks_per_transaction` requires the DB to restart.

* `type` - The type of the value of `max_pred_locks_per_transaction`.

## max_replication_slots

The following arguments are supported in the `max_replication_slots` specification block:

* `description` - The description of `max_replication_slots`.

* `maximum` - The maximum valid value for `max_replication_slots`.

* `minimum` - The minimum valid value for `max_replication_slots`.

* `requires_restart` - Whether changing the value of `max_replication_slots` requires the DB to restart.

* `type` - The type of the value of `max_replication_slots`.

## max_slot_wal_keep_size

The following arguments are supported in the `max_slot_wal_keep_size` specification block:

* `description` - The description of `max_slot_wal_keep_size`.

* `maximum` - The maximum valid value for `max_slot_wal_keep_size`.

* `minimum` - The minimum valid value for `max_slot_wal_keep_size`.

* `requires_restart` - Whether changing the value of `max_slot_wal_keep_size` requires the DB to restart.

* `type` - The type of the value of `max_slot_wal_keep_size`.

## max_stack_depth

The following arguments are supported in the `max_stack_depth` specification block:

* `description` - The description of `max_stack_depth`.

* `maximum` - The maximum valid value for `max_stack_depth`.

* `minimum` - The minimum valid value for `max_stack_depth`.

* `requires_restart` - Whether changing the value of `max_stack_depth` requires the DB to restart.

* `type` - The type of the value of `max_stack_depth`.

## max_standby_archive_delay

The following arguments are supported in the `max_standby_archive_delay` specification block:

* `description` - The description of `max_standby_archive_delay`.

* `maximum` - The maximum valid value for `max_standby_archive_delay`.

* `minimum` - The minimum valid value for `max_standby_archive_delay`.

* `requires_restart` - Whether changing the value of `max_standby_archive_delay` requires the DB to restart.

* `type` - The type of the value of `max_standby_archive_delay`.

## max_standby_streaming_delay

The following arguments are supported in the `max_standby_streaming_delay` specification block:

* `description` - The description of `max_standby_streaming_delay`.

* `maximum` - The maximum valid value for `max_standby_streaming_delay`.

* `minimum` - The minimum valid value for `max_standby_streaming_delay`.

* `requires_restart` - Whether changing the value of `max_standby_streaming_delay` requires the DB to restart.

* `type` - The type of the value of `max_standby_streaming_delay`.

## max_wal_senders

The following arguments are supported in the `max_wal_senders` specification block:

* `description` - The description of `max_wal_senders`.

* `maximum` - The maximum valid value for `max_wal_senders`.

* `minimum` - The minimum valid value for `max_wal_senders`.

* `requires_restart` - Whether changing the value of `max_wal_senders` requires the DB to restart.

* `type` - The type of the value of `max_wal_senders`.

## max_worker_processes

The following arguments are supported in the `max_worker_processes` specification block:

* `description` - The description of `max_worker_processes`.

* `maximum` - The maximum valid value for `max_worker_processes`.

* `minimum` - The minimum valid value for `max_worker_processes`.

* `requires_restart` - Whether changing the value of `max_worker_processes` requires the DB to restart.

* `type` - The type of the value of `max_worker_processes`.

## password_encryption

The following arguments are supported in the `password_encryption` specification block:

* `description` - The description of the `password_encryption` setting.

* `enum` - A list of valid values for the `password_encryption` setting.

* `example` - An example value for the `password_encryption` setting.

* `requires_restart` - Whether changing the value of `password_encryption` requires the DB to restart.

* `type` - A list of types for the `password_encryption` setting.

## pg_partman_bgw_interval

The following arguments are supported in the `pg_partman_bgw_interval` specification block:

* `description` - The description of the `pg_partman_bgw_interval` setting.

* `example` - An example value for the `pg_partman_bgw_interval` setting.

* `maximum` - The maximum allowed value for the `pg_partman_bgw_interval` setting.

* `minimum` - The minimum allowed value for the `pg_partman_bgw_interval` setting.

* `requires_restart` - Whether changing the value of `pg_partman_bgw_interval` requires the DB to restart.

* `type` - The type of the `pg_partman_bgw_interval` setting.

## pg_partman_bgw_role

The following arguments are supported in the `pg_partman_bgw_role` specification block:

* `description` - The description of the `pg_partman_bgw_role` setting.

* `example` - An example value for the `pg_partman_bgw_role` setting.

* `maxLength` - The maximum length for the `pg_partman_bgw_role` setting.

* `pattern` - The regular expression pattern for validating the `pg_partman_bgw_role` setting.

* `requires_restart` - Whether changing the value of `pg_partman_bgw_role` requires the DB to restart.

* `type` - The type of the `pg_partman_bgw_role` setting.

## pg_stat_monitor_pgsm_enable_query_plan

The following arguments are supported in the `pg_stat_monitor_pgsm_enable_query_plan` specification block:

* `description` - The description of the `pg_stat_monitor_pgsm_enable_query_plan` setting.

* `example` - An example value for the `pg_stat_monitor_pgsm_enable_query_plan` setting.

* `requires_restart` - Whether changing the value of `pg_stat_monitor_pgsm_enable_query_plan` requires the DB to restart.

* `type` - The type of the `pg_stat_monitor_pgsm_enable_query_plan` setting.

## pg_stat_monitor_pgsm_max_buckets

The following arguments are supported in the `pg_stat_monitor_pgsm_max_buckets` specification block:

* `description` - The description of the `pg_stat_monitor_pgsm_max_buckets` setting.

* `example` - An example value for the `pg_stat_monitor_pgsm_max_buckets` setting.

* `maximum` - The maximum allowed value for the `pg_stat_monitor_pgsm_max_buckets` setting.

* `minimum` - The minimum allowed value for the `pg_stat_monitor_pgsm_max_buckets` setting.

* `requires_restart` - Whether changing the value of `pg_stat_monitor_pgsm_max_buckets` requires the DB to restart.

* `type` - The type of the `pg_stat_monitor_pgsm_max_buckets` setting.

## pg_stat_statements_track

The following arguments are supported in the `pg_stat_statements_track` specification block:

* `description` - The description of the `pg_stat_statements_track` setting.

* `enum` - A list of valid values for the `pg_stat_statements_track` setting.

* `requires_restart` - Whether changing the value of `pg_stat_statements_track` requires the DB to restart.

* `type` - The type of the `pg_stat_statements_track` setting.

## temp_file_limit

The following arguments are supported in the `temp_file_limit` specification block:

* `description` - The description of the `temp_file_limit` setting.

* `example` - An example value for the `temp_file_limit` setting.

* `maximum` - The maximum allowed value for the `temp_file_limit` setting.

* `minimum` - The minimum allowed value for the `temp_file_limit` setting.

* `requires_restart` - Whether changing the value of `temp_file_limit` requires the DB to restart.

* `type` - The type of the `temp_file_limit` setting.

## timezone

The following arguments are supported in the `timezone` specification block:

* `description` - The description of the `timezone` setting.

* `example` - An example value for the `timezone` setting.

* `maxLength` - The maximum length for the `timezone` setting.

* `pattern` - The regular expression pattern for validating the `timezone` setting.

* `requires_restart` - Whether changing the value of `timezone` requires the DB to restart.

* `type` - The type of the `timezone` setting.

## track_activity_query_size

The following arguments are supported in the `track_activity_query_size` specification block:

* `description` - The description of the `track_activity_query_size` setting.

* `example` - An example value for the `track_activity_query_size` setting.

* `maximum` - The maximum allowed value for the `track_activity_query_size` setting.

* `minimum` - The minimum allowed value for the `track_activity_query_size` setting.

* `requires_restart` - Whether changing the value of `track_activity_query_size` requires the DB to restart.

* `type` - The type of the `track_activity_query_size` setting.

## track_commit_timestamp

The following arguments are supported in the `track_commit_timestamp` specification block:

* `description` - The description of the `track_commit_timestamp` setting.

* `enum` - A list of valid values for the `track_commit_timestamp` setting.

* `example` - An example value for the `track_commit_timestamp` setting.

* `requires_restart` - Whether changing the value of `track_commit_timestamp` requires the DB to restart.

* `type` - The type of the `track_commit_timestamp` setting.

## track_functions

The following arguments are supported in the `track_functions` specification block:

* `description` - The description of the `track_functions` setting.

* `enum` - A list of valid values for the `track_functions` setting.

* `requires_restart` - Whether changing the value of `track_functions` requires the DB to restart.

* `type` - The type of the `track_functions` setting.

## track_io_timing

The following arguments are supported in the `track_io_timing` specification block:

* `description` - The description of the `track_io_timing` setting.

* `enum` - A list of valid values for the `track_io_timing` setting.

* `example` - An example value for the `track_io_timing` setting.

* `requires_restart` - Whether changing the value of `track_io_timing` requires the DB to restart.

* `type` - The type of the `track_io_timing` setting.

## wal_sender_timeout

The following arguments are supported in the `wal_sender_timeout` specification block:

* `description` - The description of the `wal_sender_timeout` setting.

* `example` - An example value for the `wal_sender_timeout` setting.

* `requires_restart` - Whether changing the value of `wal_sender_timeout` requires the DB to restart.

* `type` - The type of the `wal_sender_timeout` setting.

## wal_writer_delay

The following arguments are supported in the `wal_writer_delay` specification block:

* `description` - The description of the `wal_writer_delay` setting.

* `example` - An example value for the `wal_writer_delay` setting.

* `maximum` - The maximum allowed value for the `wal_writer_delay` setting.

* `minimum` - The minimum allowed value for the `wal_writer_delay` setting.

* `requires_restart` - Whether changing the value of `wal_writer_delay` requires the DB to restart.

* `type` - The type of the `wal_writer_delay` setting.
