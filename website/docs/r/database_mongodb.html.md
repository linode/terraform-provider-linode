---
layout: "linode"
page_title: "Linode: linode_database_mongodb"
sidebar_current: "docs-linode-resource-database-mongodb"
description: |-
  Manages a Linode MongoDB Database.
---

# linode\_database\_mongodb

Provides a Linode Mongo Database resource. This can be used to create, modify, and delete Linode MongoDB Databases.
For more information, see the [Linode APIv4 docs](https://www.linode.com/docs/api/databases/).

Please keep in mind that Managed Databases can take up to an hour to provision.

## Example Usage

Creating a simple MongoDB database instance:

```hcl
resource "linode_database_mongodb" "foobar" {
  label = "mydatabase"
  engine_id = "mongodb/4.4.10"
  region = "us-southeast"
  type = "g6-nanode-1"
}
```

Creating a complex MongoDB database instance:

```hcl
resource "linode_database_mongodb" "foobar" {
  label = "mydatabase"
  engine_id = "mongodb/4.4.10"
  region = "us-southeast"
  type = "g6-nanode-1"

  allow_list = ["0.0.0.0/0"]
  cluster_size = 3
  compression_type = "zlib"
  storage_engine = "wiredtiger"
  encrypted = true
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

* `engine_id` - (Required) The Managed Database engine in engine/version format. (e.g. `mongo/4.4.10`)

* `label` - (Required) A unique, user-defined string referring to the Managed Database.

* `region` - (Required) The region to use for the Managed Database.

* `type` - (Required) The Linode Instance type used for the nodes of the  Managed Database instance.

- - -

* `allow_list` - (Optional) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format. Use `linode_database_access_controls` to manage your allow list separately.

* `cluster_size` - (Optional) The number of Linode Instance nodes deployed to the Managed Database. (default `1`)

* `compression_type` - (Optional) The type of data compression for this Database. (`none`, `snappy`, `zlib`; default `none`)

* `encrypted` - (Optional) Whether the Managed Databases is encrypted. (default `false`)

* `ssl_connection` - (Optional) Whether to require SSL credentials to establish a connection to the Managed Database. (default `false`)

* `storage_engine` - (Optional) The type of storage engine for this Database. (`mmapv1`, `wiredtiger`; default `wiredtiger`)

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

* `engine` - The Managed Database engine. (e.g. `mongodb`)

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private network host for the Managed Database.

* `peers` - A set of peer addresses for this Database.

* `port` - The access port for this Managed Database.

* `replica_set` - Label for configuring a MongoDB replica set. Choose the same label on multiple Databases to include them in the same replica set.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `status` - The operating status of the Managed Database.

* `updated` - When this Managed Database was last updated.

* `version` - The Managed Database engine version. (e.g. `v8.0.26`)

## Import

Linode MongoDB Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_mongodb.foobar 1234567
```
