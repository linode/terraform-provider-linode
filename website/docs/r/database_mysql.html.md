---
layout: "linode"
page_title: "Linode: linode_database_mysql"
sidebar_current: "docs-linode-resource-database-mysql"
description: |-
  Manages a Linode MySQL Database.
---

# linode\_database\_mysql

Provides a Linode MySQL Database resource. This can be used to create, modify, and delete Linode MySQL Databases.
For more information, see the [Linode APIv4 docs](https://www.linode.com/docs/api/databases/).

## Example Usage

Creating a simple MySQL database instance:

```hcl
resource "linode_database_mysql" "foobar" {
  label = "mydatabase"
  engine = "mysql/8.0.26"
  region = "us-southeast"
  type = "g6-nanode-1"
}
```

Creating a complex MySQL database instance:

```hcl
resource "linode_database_mysql" "foobar" {
  label = "mydatabase"
  engine = "mysql/8.0.26"
  region = "us-southeast"
  type = "g6-nanode-1"

  allow_list = ["0.0.0.0/0"]
  cluster_size = 3
  encrypted = true
  replication_type = "asynch"
  ssl_connection = true
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required) The Managed Database engine in engine/version format.

* `label` - (Required) A unique, user-defined string referring to the Managed Database.

* `region` - (Required) The region to use for the Managed Database.

* `type` - (Required) The Linode Instance type used for the nodes of the  Managed Database instance.

- - -

* `allow_list` - (Optional) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format.

* `cluster_size` - (Optional) The number of Linode Instance nodes deployed to the Managed Database. (default `1`)

* `encrypted` - (Optional) Whether the Managed Databases is encrypted. (default `false`)

* `replication_type` - (Optional) The replication method used for the Managed Database. (`none`, `asynch`, `semi_synch`; default `none`)

* `ssl_connection` - (Optional) Whether to require SSL credentials to establish a connection to the Managed Database. (default `false`)

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Managed Database.

* `ca_cert` - The base64-encoded SSL CA certificate for the Managed Database instance.

* `created` - When this Managed Database was created.

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private network host for the Managed Database.

* `root_password` - The randomly-generated root password for the Managed Database instance.

* `status` - The operating status of the Managed Database.

* `updated` - When this Managed Database was last updated.

* `root_username` - The root username for the Managed Database instance.

* `version` - The Managed Database engine version.

## Import

Linode MySQL Databases can be imported using the `id`, e.g.

```sh
terraform import linode_database_mysql.foobar 1234567
```
