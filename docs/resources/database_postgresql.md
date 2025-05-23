---
page_title: "Linode: linode_database_postgresql"
description: |-
  Manages a Linode PostgreSQL Database.
---

# linode\_database\_postgresql

~> **DEPRECATION NOTICE:** This resource has been deprecated. Please use [linode_database_postgresql_v2](database_postgresql_v2.html.markdown) for all future implementations.

Provides a Linode PostgreSQL Database resource. This can be used to create, modify, and delete Linode PostgreSQL Databases.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-databases-postgre-sql-instances).

Please keep in mind that Managed Databases can take up to an hour to provision.

## Example Usage

Creating a simple PostgreSQL database instance:

```hcl
resource "linode_database_postgresql" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/13.2"
  region = "us-southeast"
  type = "g6-nanode-1"
}
```

Creating a complex PostgreSQL database instance:

```hcl
resource "linode_database_postgresql" "foobar" {
  label = "mydatabase"
  engine_id = "postgresql/13.2"
  region = "us-southeast"
  type = "g6-nanode-1"

  allow_list = ["0.0.0.0/0"]
  cluster_size = 3
  encrypted = true
  replication_type = "semi_synch"
  replication_commit_type = "remote_write"
  ssl_connection = true

  updates {
    day_of_week = "saturday"
    duration = 1
    frequency = "monthly"
    hour_of_day = 22
    week_of_month = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_id` - (Required) The Managed Database engine in engine/version format. (e.g. `postgresql/13.2`)

* `label` - (Required) A unique, user-defined string referring to the Managed Database.

* `region` - (Required) The region to use for the Managed Database.

* `type` - (Required) The Linode Instance type used for the nodes of the  Managed Database instance.

- - -

* `allow_list` - (Optional) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format. Use `linode_database_access_controls` to manage your allow list separately.

* `cluster_size` - (Optional) The number of Linode Instance nodes deployed to the Managed Database. (default `1`)

* `encrypted` - (Optional) Whether the Managed Databases is encrypted. (default `false`)

* `replication_type` - (Optional) The replication method used for the Managed Database. (`none`, `asynch`, `semi_synch`; default `none`)

  * Must be `none` for a single node cluster.

  * Must be `asynch` or `semi_synch` for a high availability cluster.

* `replication_commit_type` - (Optional) The synchronization level of the replicating server. (`on`, `local`, `remote_write`, `remote_apply`, `off`; default `off`)

  * Must be `local` or `off` for the `asynch` replication type.

  * Must be `on`, `remote_write`, or `remote_apply` for the `semi_synch` replication type.

* `ssl_connection` - (Optional) Whether to require SSL credentials to establish a connection to the Managed Database. (default `false`)

* [`updates`](#updates) - (Optional) Configuration settings for automated patch update maintenance for the Managed Database.

## updates

The following arguments are supported in the `updates` specification block:

* `day_of_week` - (Required) The day to perform maintenance. (`monday`, `tuesday`, ...)

* `duration` - (Required) The maximum maintenance window time in hours. (`1`..`3`)

* `frequency` - (Required) Whether maintenance occurs on a weekly or monthly basis. (`weekly`, `monthly`)

* `hour_of_day` - (Required) The hour to begin maintenance based in UTC time. (`0`..`23`)

* `week_of_month` - (Optional) The week of the month to perform monthly frequency updates. Required for `monthly` frequency updates. (`1`..`4`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Managed Database.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database instance.

* `created` - When this Managed Database was created.

* `engine` - The Managed Database engine. (e.g. `postgresql`)

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private network host for the Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `status` - The operating status of the Managed Database.

* `updated` - When this Managed Database was last updated.

* `version` - The Managed Database engine version. (e.g. `13.2`)

## Import

Linode PostgreSQL Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_postgresql.foobar 1234567
```
