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

Creating a simple PostgreSQL database:

```hcl
resource "linode_database_postgresql_v2" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/16"
  region = "us-mia"
  type = "g6-nanode-1"
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

  updates {
    duration = 4
    frequency = "weekly"
    hour_of_day = 22
    week_of_month = 2
  }
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

* [`updates`](#updates) - (Optional) Configuration settings for automated patch update maintenance for the Managed Database.

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

* `frequency` - (Required) Whether maintenance occurs on a weekly or monthly basis. (`weekly`, `monthly`)

* `hour_of_day` - (Required) The hour to begin maintenance based in UTC time. (`0`..`23`)

* `week_of_month` - (Optional) The week of the month to perform monthly frequency updates. Required for `monthly` frequency updates. (`1`..`4`)

## Import

Linode PostgreSQL Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_postgresql_v2.foobar 1234567
```
