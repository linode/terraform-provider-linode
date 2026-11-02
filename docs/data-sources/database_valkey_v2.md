---
page_title: "Linode: linode_database_valkey_v2"
description: |-
  Provides information about a Linode Valkey Database.
---

# Data Source: linode\_database\_valkey\_v2

Provides information about a Linode Valkey Database.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-databases-valkey-instance).

## Example Usage

Get information about a Valkey database:

```hcl
data "linode_database_valkey_v2" "my-db" {
  id = "12345"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Valkey database.

## Attributes Reference

The `linode_database_valkey_v2` data source exports the following attributes:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `allow_list` - A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format. Use `linode_database_access_controls` to manage your allow list separately; do not also manage the same list through the database resource.

* `available_restore_times` - Available discrete snapshot timestamps that can be used to fork this database. These are not a continuous point-in-time recovery range; use a returned timestamp or omit the restore time to request the latest snapshot.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database. This attribute is sensitive.

* `cluster_size` - The number of Linode Instance nodes deployed to the Managed Database.

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Database is encrypted.

* `engine` - The Managed Database engine. (e.g. `valkey`)

* `engine_id` - The Managed Database engine in engine/version format. (e.g. `valkey/8.1`)

* `fork_restore_time` - The snapshot timestamp used when this database was forked.

* `fork_source` - The ID of the database that was forked from.

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary host for the Managed Database. Deprecated; use `host_standby` instead.

* `host_standby` - The standby host for the Managed Database.

* `label` - A unique, user-defined string referring to the Managed Database.

* `members` - A mapping of cluster members to their roles.

* [`pending_updates`](#pending_updates) - (Nested Attribute Set) A set of pending updates. Set elements can't be referenced by index; use a `for` expression or `tolist(...)` to access them.

* `platform` - The back-end platform used by the service.

* `port` - The access port for this Managed Database.

* [`private_network`](#private_network) - (Nested Attribute) Restricts access to this database using a virtual private cloud (VPC). Referenced directly (e.g. `private_network.public_access`).

* `region` - The region of the Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance. This attribute is sensitive.

* `root_username` - The root username for the Managed Database instance. This attribute is sensitive.

* `ssl_connection` - Whether SSL is required to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `suspended` - Whether this Managed Database is suspended.

* `type` - The Managed Database plan type used for the cluster nodes.

* `updated` - When this Managed Database was last updated.

* [`updates`](#updates) - (Nested Attribute) Configuration settings for automated patch update maintenance for the Managed Database. Referenced directly (e.g. `updates.day_of_week`).

* `version` - The Managed Database engine version. (e.g. `8.1.0`)

* `engine_config_backup_hour` - The UTC hour when a backup starts (`0`..`23`).

* `engine_config_backup_minute` - The minute within the hour when a backup starts (`0`..`59`).

* `engine_config_frequent_snapshots` - Whether frequent local RDB snapshots are enabled.

* `engine_config_valkey_acl_channels_default` - The default pub/sub channel ACL for new users: `allchannels` or `resetchannels`.

* `engine_config_valkey_active_expire_effort` - The effort used to reclaim expired keys (`1`..`10`).

* `engine_config_valkey_activedefrag` - Whether active memory defragmentation is enabled.

* `engine_config_valkey_lfu_decay_time` - The LFU counter decay time in minutes (`1`..`120`).

* `engine_config_valkey_lfu_log_factor` - The logarithm factor for LFU counters (`0`..`100`).

* `engine_config_valkey_maxmemory_policy` - The memory eviction policy: `noeviction`, `allkeys-lru`, `volatile-lru`, `allkeys-random`, `volatile-random`, `volatile-ttl`, `volatile-lfu`, or `allkeys-lfu`.

* `engine_config_valkey_number_of_databases` - The number of logical Valkey databases (`1`..`128`). Changing this setting restarts Valkey.

* `engine_config_valkey_persistence` - The persistence mode: `off` or `rdb`. Disabling persistence disables backups and forking.

* `engine_config_valkey_pubsub_client_output_buffer_limit` - The hard output-buffer limit for pub/sub clients in MB (`32`..`262144`).

* `engine_config_valkey_timeout` - The idle connection timeout in seconds (`0`..`2073600`).

Sensitive attributes are hidden in Terraform output, but Terraform still stores them in state. Protect and restrict access to the state file and its backups.

## pending_updates

The following arguments are exposed by each entry in the `pending_updates` attribute:

* `deadline` - The time when a mandatory update needs to be applied.

* `description` - A description of the update.

* `planned_for` - The date and time a maintenance update will be applied.

## updates

The following arguments are exposed by the `updates` attribute:

* `day_of_week` - The numeric reference for the day of the week to perform maintenance. 1 is Monday, 2 is Tuesday, through to 7 which is Sunday.

* `duration` - The maximum maintenance window time in hours. (`4`)

* `frequency` - The frequency at which maintenance occurs. (`weekly`)

* `hour_of_day` - The hour to begin maintenance based in UTC time. (`0`..`23`)

## private_network

The following arguments are exposed by the `private_network` attribute:

* `vpc_id` - The ID of the virtual private cloud (VPC) to restrict access to this database using.

* `subnet_id` - The ID of the VPC subnet to restrict access to this database using.

* `public_access` - If true, clients outside the VPC can connect to the database using a public IP address.
