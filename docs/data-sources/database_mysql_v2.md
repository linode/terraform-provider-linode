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

* `region` - The region to use for the Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `ssl_connection` - Whether to require SSL credentials to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `type` - The Linode Instance type used for the nodes of the Managed Database.

* `updated` - When this Managed Database was last updated.

* [`updates`](#updates) - Configuration settings for automated patch update maintenance for the Managed Database.

* `version` - The Managed Database engine version. (e.g. `13.2`)

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
