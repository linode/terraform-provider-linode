---
page_title: "Linode: linode_database_valkey_v2"
description: |-
  Manages a Linode Valkey Managed Database.
---

# linode\_database\_valkey\_v2

Provides a Linode Valkey Database resource. This can be used to create, modify, and delete Linode Valkey Databases.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-databases-valkey-instances).

Please keep in mind that Managed Databases can take up to half an hour to provision.

## Example Usage

Creating a Valkey database with RDB persistence:

```hcl
resource "linode_database_valkey_v2" "cache" {
  label     = "my-valkey"
  engine_id = "valkey/8.1"
  region    = "us-mia"
  type      = "g7-dedicated-4-2"

  engine_config_valkey_persistence = "rdb"
}
```

The `type` is a Managed Database plan. Use the `linode_database_engines`
data source and the Linode API to confirm engine IDs, plan types, and region
availability for your account.

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - (Required) A unique, user-defined database label.

* `engine_id` - (Required) The Valkey engine and version, such as `valkey/8.1`. The engine component cannot be changed to another database engine.

* `region` - (Required) The deployment region. Changing the region replaces the resource.

* `type` - (Required) The Managed Database plan type used for the cluster nodes.

* `allow_list` - (Optional) IP addresses or CIDR ranges allowed to connect. Set to `[]` to clear the list. If using `linode_database_access_controls`, do not also manage the same list through this resource.

* `cluster_size` - (Optional) The number of nodes in the database cluster. Defaults to `1`.

* `fork_source` - (Optional) The ID of a Valkey database to restore from. Changing this value replaces the resource.

* `fork_restore_time` - (Optional) An RFC3339 snapshot timestamp from the source database's `available_restore_times`. Requires `fork_source`. If omitted, the API uses its latest snapshot. Changing this value replaces the resource. This is a discrete snapshot restore, not arbitrary point-in-time recovery.

* `suspended` - (Optional) Whether the database should be suspended. Defaults to `false`. Resuming a suspended database can take time; configure the update timeout as needed.

* [`updates`](#updates) - (Optional, Nested Attribute) Automated patch maintenance settings, configured as `updates = { ... }`.

* [`private_network`](#private_network) - (Optional, Nested Attribute) VPC and subnet settings, configured as `private_network = { ... }`. The VPC must be in the database's region.

* [`timeouts`](#timeouts) - (Optional, Block) Operation timeout settings, configured as `timeouts { ... }`.

The following engine configuration fields are optional and computed. Omitted values are populated from the API. Terraform `null` does not request an explicit API-null reset of nullable settings. Use the `linode_database_valkey_config` data source to inspect the configuration catalog and restart requirements.

* `engine_config_backup_hour` - The UTC hour when a backup starts (`0`..`23`).

* `engine_config_backup_minute` - The minute within the hour when a backup starts (`0`..`59`).

* `engine_config_frequent_snapshots` - Whether frequent local RDB snapshots are enabled.

* `engine_config_valkey_acl_channels_default` - The default pub/sub channel ACL for new users: `allchannels` or `resetchannels`.

* `engine_config_valkey_active_expire_effort` - The effort used to reclaim expired keys (`1`..`10`).

* `engine_config_valkey_activedefrag` - Whether active memory defragmentation is enabled.

* `engine_config_valkey_lfu_decay_time` - The LFU counter decay time in minutes (`1`..`120`).

* `engine_config_valkey_lfu_log_factor` - The logarithm factor for LFU counters (`0`..`100`).

* `engine_config_valkey_maxmemory_policy` - The memory eviction policy: `noeviction`, `allkeys-lru`, `volatile-lru`, `allkeys-random`, `volatile-random`, `volatile-ttl`, `volatile-lfu`, or `allkeys-lfu`.

* `engine_config_valkey_number_of_databases` - The number of logical Valkey databases (`1`..`128`). Changing this value restarts Valkey.

* `engine_config_valkey_persistence` - The persistence mode: `off` or `rdb`. Disabling persistence disables backups and forking.

* `engine_config_valkey_pubsub_client_output_buffer_limit` - The hard output-buffer limit for pub/sub clients in MB (`32`..`262144`).

* `engine_config_valkey_timeout` - The idle connection timeout in seconds (`0`..`2073600`).

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The unique ID of the Managed Database.

* `available_restore_times` - Available discrete snapshot timestamps, formatted as RFC3339 strings. Snapshot availability changes over time, so the API validates a requested timestamp when it creates a fork. Do not assume an arbitrary timestamp is valid.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database. This attribute is sensitive.

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Database is encrypted.

* `engine` - The Managed Database engine (`valkey`).

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary host for the Managed Database. Deprecated; use `host_standby` instead.

* `host_standby` - The standby host for the Managed Database.

* `members` - A mapping of cluster members to their roles.

* [`pending_updates`](#pending_updates) - (Nested Attribute Set) A set of pending updates. Set elements cannot be referenced by index; use a `for` expression or `tolist(...)` to access them.

* `platform` - The back-end platform used by the service.

* `port` - The access port for the Managed Database.

* `root_password` - The generated root password for the Managed Database instance. This attribute is sensitive.

* `root_username` - The root username for the Managed Database instance. This attribute is sensitive.

* `ssl_connection` - Whether SSL is required to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `updated` - When this Managed Database was last updated.

* `version` - The Managed Database engine version, such as `8.1.0`.

Sensitive attributes are hidden in Terraform output, but Terraform still stores them in state. Protect and restrict access to the state file and its backups. While suspended, the provider skips SSL and credential API reads and preserves any known certificate and credentials in state.

## pending_updates

The following attributes are exposed by each entry in the `pending_updates` set:

* `deadline` - The time when a mandatory update needs to be applied.

* `description` - A description of the update.

* `planned_for` - The date and time a maintenance update will be applied.

## updates

The following arguments are supported by the `updates` nested attribute:

* `day_of_week` - (Optional) The numeric day of the week for maintenance: `1` is Monday, through to `7` for Sunday.

* `duration` - (Optional) The maximum maintenance window in hours. Currently, only `4` is supported.

* `frequency` - (Optional) The maintenance frequency. Currently, only `weekly` is supported; defaults to `weekly`.

* `hour_of_day` - (Optional) The UTC hour when maintenance begins (`0`..`23`).

For example:

```hcl
updates = {
  day_of_week = 2
  duration    = 4
  frequency   = "weekly"
  hour_of_day = 3
}
```

## private_network

The following arguments are supported by the `private_network` nested attribute:

* `vpc_id` - (Required) The ID of the VPC used to restrict access to this database.

* `subnet_id` - (Required) The ID of the subnet in that VPC.

* `public_access` - (Optional) Whether clients outside the VPC can connect using a public IP address. Defaults to `false`.

For example, within the database resource:

```hcl
private_network = {
  vpc_id        = linode_vpc.example.id
  subnet_id     = linode_vpc_subnet.example.id
  public_access = false
}
```

## Timeouts

The `timeouts` block supports the following arguments:

* `create` - (Optional) The timeout for database creation, including readiness waits and initial configuration. Defaults to `60m`.

* `update` - (Optional) The timeout for database updates, including suspension/resume and readiness waits. Defaults to `60m`.

* `delete` - (Optional) The timeout for the database delete operation. Defaults to `5m`.

For example, within the database resource:

```hcl
timeouts {
  create = "90m"
  update = "90m"
  delete = "10m"
}
```

## Import

Import a database by its numeric ID:

```shell
terraform import linode_database_valkey_v2.cache 12345
```

Import reads the existing database into state; it does not create a new database.
