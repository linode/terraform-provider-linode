---
page_title: "Linode: linode_database_valkey_v2"
description: |-
  Manages a Linode Valkey Managed Database.
---

# linode_database_valkey_v2

Creates and manages a Valkey Managed Database cluster. The database may take
several minutes to provision. Valkey availability is limited; confirm that
your account and region can create Valkey instances before applying a
configuration.

## Example Usage

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

- `label` (Required) A unique, user-defined database label.
- `engine_id` (Required) The Valkey engine and version, such as `valkey/8.1`.
- `region` (Required) The deployment region. Changing the region replaces
the resource.
- `type` (Required) The Managed Database plan type.
- `allow_list` (Optional) IP addresses or CIDR ranges allowed to connect.
- `cluster_size` (Optional) The cluster size. Defaults to `1`.
- `fork_source` (Optional) The ID of a Valkey database to restore from.
- `fork_restore_time` (Optional) A snapshot time from the source database's
  `available_restore_times`. If omitted, the API uses its latest snapshot.
  This is a discrete snapshot restore, not point-in-time recovery.
- `suspended` (Optional) Whether the database should be suspended.
- `updates` (Optional) Automated maintenance schedule.
- `private_network` (Optional) VPC and subnet settings for database access.

The engine configuration fields are optional and computed:

- `engine_config_backup_hour` (0–23) UTC hour to start a backup.
- `engine_config_backup_minute` (0–59) Minute within the hour to start a
  backup.
- `engine_config_frequent_snapshots` Whether frequent local RDB snapshots
  are enabled.
- `engine_config_valkey_acl_channels_default` Pub/sub channel ACL for new
  users: `allchannels` or `resetchannels`.
- `engine_config_valkey_active_expire_effort` Expiration cleanup effort
  (1–10).
- `engine_config_valkey_activedefrag` Whether active memory defragmentation
  is enabled.
- `engine_config_valkey_lfu_decay_time` LFU counter decay time in minutes
  (1–120).
- `engine_config_valkey_lfu_log_factor` LFU counter logarithm factor
  (0–100).
- `engine_config_valkey_maxmemory_policy` Memory eviction policy. Supported
  values are `noeviction`, `allkeys-lru`, `volatile-lru`, `allkeys-random`,
  `volatile-random`, `volatile-ttl`, `volatile-lfu`, and `allkeys-lfu`.
- `engine_config_valkey_number_of_databases` Number of logical databases
  (1–128). Changing this value restarts Valkey.
- `engine_config_valkey_persistence` Persistence mode: `off` or `rdb`.
  Disabling persistence disables backups and forking.
- `engine_config_valkey_pubsub_client_output_buffer_limit` Pub/sub output
  buffer hard limit in MB (32–262144).
- `engine_config_valkey_timeout` Idle connection timeout in seconds
  (0–2073600).

## Attribute Reference

The resource exports `id`, `engine`, `version`, `status`, `created`, `updated`,
`encrypted`, `platform`, `port`, `ssl_connection`, `host_primary`,
`host_standby`, `members`, `pending_updates`, `available_restore_times`,
`fork_source`, and `fork_restore_time`.

`available_restore_times` is a computed list of snapshot timestamps. Snapshot
availability changes over time, so the API validates a requested restore time
when it creates the fork. Do not assume an arbitrary timestamp is valid.

`root_username`, `root_password`, and `ca_cert` are marked sensitive in
Terraform output. Terraform still stores them in state; protect and restrict
access to the state file and its backups.

## Import

Import a database by its numeric ID:

```shell
terraform import linode_database_valkey_v2.cache 12345
```

Changing `region`, `fork_source`, or `fork_restore_time` replaces the
resource. Valkey resume operations may take time; use a suitable resource
update timeout and verify the database reaches its expected status.
