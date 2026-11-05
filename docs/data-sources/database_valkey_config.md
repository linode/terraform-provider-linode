---
page_title: "Linode: linode_database_valkey_config"
description: |-
  Provides information about available Linode Valkey Database configuration options.
---

# Data Source: linode\_database\_valkey\_config

Provides information about available Linode Valkey Database configuration options, including accepted types, ranges, enums, defaults or examples, and restart requirements.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-databases-valkey-config).

## Example Usage

Get information about Valkey configuration options:

```hcl
data "linode_database_valkey_config" "example" {}

output "backup_hour_types" {
  value = data.linode_database_valkey_config.example.backup_hour[0].type
}
```

## Argument Reference

This data source does not accept any arguments.

## Attributes Reference

The `linode_database_valkey_config` data source exports the following attributes:

* `id` - A content-derived identifier for the configuration catalog.

Each configuration option below is a computed list containing one metadata object. Access metadata using an index, for example `valkey_timeout[0].maximum`. These are catalog values, not the current configuration of a database cluster; use `linode_database_valkey_v2` to read a cluster's settings.

Every metadata object exposes `description`, `requires_restart`, and `type`. The `type` attribute is a list of strings, including `null` when the API supports a nullable value (for example `["integer", "null"]`). Additional fields vary by option as listed below; absent metadata is not synthesized.

* `backup_hour` - The UTC hour when a backup starts. Exposes `example`, `minimum`, and `maximum`.

* `backup_minute` - The minute within the hour when a backup starts. Exposes `example`, `minimum`, and `maximum`.

* `frequent_snapshots` - Whether frequent local RDB snapshots are enabled. Exposes a boolean `default`.

* `valkey_acl_channels_default` - The default pub/sub channel ACL for new users. Exposes an `enum` list of accepted strings.

* `valkey_active_expire_effort` - The effort used to reclaim expired keys. Exposes integer `default`, `minimum`, and `maximum` values.

* `valkey_activedefrag` - Whether active memory defragmentation is enabled. Exposes a boolean `default`.

* `valkey_lfu_decay_time` - The LFU counter decay time in minutes. Exposes integer `default`, `minimum`, and `maximum` values.

* `valkey_lfu_log_factor` - The logarithm factor for LFU counters. Exposes integer `default`, `minimum`, and `maximum` values.

* `valkey_maxmemory_policy` - The memory eviction policy. Exposes a string `default` and an `enum` list of accepted strings.

* `valkey_number_of_databases` - The number of logical Valkey databases. Exposes `example`, `minimum`, and `maximum`.

* `valkey_persistence` - The persistence mode. Exposes an `enum` list of accepted strings. Persistence mode `off` disables backups and forking.

* `valkey_pubsub_client_output_buffer_limit` - The hard output-buffer limit for pub/sub clients in MB. Exposes `example`, `minimum`, and `maximum`.

* `valkey_timeout` - The idle connection timeout in seconds. Exposes integer `default`, `minimum`, and `maximum` values.

## Metadata Reference

* `description` - A description of the configuration option.

* `requires_restart` - Whether changing this option requires a database service restart.

* `type` - A list of accepted API types. A nullable type is represented by multiple entries, not a joined string.

* `default` - The default value, with its native integer, boolean, or string type. Present only for options that provide a default in the SDK catalog model.

* `example` - An integer example value. Present only for options that provide an example in the SDK catalog model.

* `minimum` - The minimum accepted integer value, when provided.

* `maximum` - The maximum accepted integer value, when provided.

* `enum` - A list of accepted strings, when provided.
