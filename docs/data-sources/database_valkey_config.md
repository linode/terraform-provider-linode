---
page_title: "Linode: linode_database_valkey_config"
description: |-
  Provides information about a Linode Valkey Database's Configuration Options.
---

# Data Source: linode\_database\_valkey\_config

Provides information about a Linode Valkey Database's Configuration Options.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-databases-valkey-config).

## Example Usage

Get information about a Valkey database's configuration options:

```hcl
data "linode_database_valkey_config" "example" {}
```

## Attributes Reference

The `linode_database_valkey_config` data source exports the following attributes:

* `id` - Unique identifier for this DataSource.

* [`backup_hour`](#backup_hour) - (Read-Only Object List) The UTC hour when a backup starts.

* [`backup_minute`](#backup_minute) - (Read-Only Object List) The minute within the hour when a backup starts.

* [`frequent_snapshots`](#frequent_snapshots) - (Read-Only Object List) Whether frequent local RDB snapshots are enabled.

* [`valkey_acl_channels_default`](#valkey_acl_channels_default) - (Read-Only Object List) The default pub/sub channel ACL for new users.

* [`valkey_active_expire_effort`](#valkey_active_expire_effort) - (Read-Only Object List) The effort used to reclaim expired keys.

* [`valkey_activedefrag`](#valkey_activedefrag) - (Read-Only Object List) Whether active memory defragmentation is enabled.

* [`valkey_lfu_decay_time`](#valkey_lfu_decay_time) - (Read-Only Object List) The LFU counter decay time in minutes.

* [`valkey_lfu_log_factor`](#valkey_lfu_log_factor) - (Read-Only Object List) The logarithm factor for LFU counters.

* [`valkey_maxmemory_policy`](#valkey_maxmemory_policy) - (Read-Only Object List) The memory eviction policy.

* [`valkey_number_of_databases`](#valkey_number_of_databases) - (Read-Only Object List) The number of logical Valkey databases.

* [`valkey_persistence`](#valkey_persistence) - (Read-Only Object List) The persistence mode.

* [`valkey_pubsub_client_output_buffer_limit`](#valkey_pubsub_client_output_buffer_limit) - (Read-Only Object List) The hard output-buffer limit for pub/sub clients in MB.

* [`valkey_timeout`](#valkey_timeout) - (Read-Only Object List) The idle connection timeout in seconds.

These are catalog values, not a database cluster's current configuration. Each option is referenced with an index, for example `backup_hour[0].description`. All option objects expose:

* `description` - The description of the configuration option.

* `requires_restart` - Whether changing the option requires the database to restart.

* `type` - A list of accepted API types, including `null` when supported (for example `["integer", "null"]`).

The additional fields for each option are described below.

## backup_hour

* `example` - An example of a valid value for `backup_hour`.
* `maximum` - The maximum valid value for `backup_hour`.
* `minimum` - The minimum valid value for `backup_hour`.

## backup_minute

* `example` - An example of a valid value for `backup_minute`.
* `maximum` - The maximum valid value for `backup_minute`.
* `minimum` - The minimum valid value for `backup_minute`.

## frequent_snapshots

* `default` - The default boolean value for `frequent_snapshots`.

## valkey_acl_channels_default

* `enum` - The accepted values for `valkey_acl_channels_default`.

## valkey_active_expire_effort

* `default` - The default value for `valkey_active_expire_effort`.
* `maximum` - The maximum valid value for `valkey_active_expire_effort`.
* `minimum` - The minimum valid value for `valkey_active_expire_effort`.

## valkey_activedefrag

* `default` - The default boolean value for `valkey_activedefrag`.

## valkey_lfu_decay_time

* `default` - The default value for `valkey_lfu_decay_time`.
* `maximum` - The maximum valid value for `valkey_lfu_decay_time`.
* `minimum` - The minimum valid value for `valkey_lfu_decay_time`.

## valkey_lfu_log_factor

* `default` - The default value for `valkey_lfu_log_factor`.
* `maximum` - The maximum valid value for `valkey_lfu_log_factor`.
* `minimum` - The minimum valid value for `valkey_lfu_log_factor`.

## valkey_maxmemory_policy

* `default` - The default value for `valkey_maxmemory_policy`.
* `enum` - The accepted values for `valkey_maxmemory_policy`.

## valkey_number_of_databases

* `example` - An example of a valid value for `valkey_number_of_databases`.
* `maximum` - The maximum valid value for `valkey_number_of_databases`.
* `minimum` - The minimum valid value for `valkey_number_of_databases`.

## valkey_persistence

* `enum` - The accepted values for `valkey_persistence`. Persistence mode `off` disables backups and forking.

## valkey_pubsub_client_output_buffer_limit

* `example` - An example of a valid value for `valkey_pubsub_client_output_buffer_limit`.
* `maximum` - The maximum valid value for `valkey_pubsub_client_output_buffer_limit`.
* `minimum` - The minimum valid value for `valkey_pubsub_client_output_buffer_limit`.

## valkey_timeout

* `default` - The default value for `valkey_timeout`.
* `maximum` - The maximum valid value for `valkey_timeout`.
* `minimum` - The minimum valid value for `valkey_timeout`.
