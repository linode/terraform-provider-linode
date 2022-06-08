---
layout: "linode"
page_title: "Linode: linode_database_mongodb"
sidebar_current: "docs-linode-datasource-database-mongodb"
description: |-
Provides information about a Linode MongoDB Database.
---

# Data Source: linode\_database\_mongodb

Provides information about a Linode MongoDB Database.

## Example Usage

Get information about a MongoDB database:

```hcl
data "linode_database_mongodb" "my-db" {
  database_id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `database_id` - The ID of the MongoDB database.

## Attributes

The `linode_database_mongodb` data source exports the following attributes:

* `allow_list` - A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database instance.

* `cluster_size` - The number of Linode Instance nodes deployed to the Managed Database.

* `compression_type` - The type of data compression for this Database. (`none`, `snappy`, `zlib`)

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Databases is encrypted.

* `engine` - The Managed Database engine. (e.g. `mongodb`)

* `engine_id` - The Managed Database engine in engine/version format. (e.g. `mongodb/4.4.10`)

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private network host for the Managed Database.

* `label` - A unique, user-defined string referring to the Managed Database.

* `peers` - A set of peer addresses for this Database.

* `port` - The access port for this Managed Database.

* `replica_set` - Label for configuring a MongoDB replica set. Choose the same label on multiple Databases to include them in the same replica set.

* `region` - The region that hosts this Linode Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `root_username` - The root username for the Managed Database instance.

* `ssl_connection` - Whether to require SSL credentials to establish a connection to the Managed Database.

* `storage_engine` - The type of storage engine for this Database. (`mmapv1`, `wiredtiger`)

* `status` - The operating status of the Managed Database.

* `type` - The Linode Instance type used for the nodes of the  Managed Database instance.

* `updated` - When this Managed Database was last updated.

* [`updates`](#updates) - (Optional) Configuration settings for automated patch update maintenance for the Managed Database.

* `version` - The Managed Database engine version. (e.g. `v8.0.26`)

## updates

The following arguments are exported by the `updates` specification block:

* `day_of_week` - The day to perform maintenance. (`monday`, `tuesday`, ...)

* `duration` - The maximum maintenance window time in hours. (`1`..`3`)

* `frequency` - Whether maintenance occurs on a weekly or monthly basis. (`weekly`, `monthly`)

* `hour_of_day` - The hour to begin maintenance based in UTC time. (`0`..`23`)

* `week_of_month` - The week of the month to perform monthly frequency updates. Required for `monthly` frequency updates. (`1`..`4`)
