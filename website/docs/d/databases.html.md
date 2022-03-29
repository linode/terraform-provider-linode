---
layout: "linode"
page_title: "Linode: linode_databases"
sidebar_current: "docs-linode-datasource-databases"
description: |-
Provides information about Linode Managed Databases.
---

# Data Source: linode\_databases

Provides information about Linode Managed Databases that match a set of filters.

## Example Usage

Get information about all Linode Managed Databases:

```hcl
data "linode_databases" "all" {}
```

Get information about all Linode MySQL Databases:

```hcl
data "linode_databases" "mysql" {
  filter {
    name = "engine"
    values = ["mysql"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `latest` - (Optional) If true, only the latest create database will be returned.

* [`filter`](#filter) - (Optional) A set of filters used to select databases that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. (`version`)

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes

Each engine will be stored in the `databases` attribute and will export the following attributes:

* `allow_list` - A list of IP addresses that can access the Managed Database.

* `cluster_size` - The number of Linode Instance nodes deployed to the Managed Database.

* `created` - When this Managed Database was created.

* `encrypted` - Whether the Managed Databases is encrypted.

* `engine` - The Managed Database engine.

* `host_primary` - The primary host for the Managed Database.

* `host_secondary` - The secondary/private network host for the Managed Database.

* `id` - The ID of the Managed Database.

* `label` - A unique, user-defined string referring to the Managed Database.

* `region` - The region to use for the Managed Database.

* `replication_type` - The replication method used for the Managed Database.

* `ssl_connection` - Whether to require SSL credentials to establish a connection to the Managed Database.

* `status` - The operating status of the Managed Database.

* `type` - The Linode Instance type used for the nodes of the  Managed Database instance.

* `updated` - When this Managed Database was last updated.

* `version` - The Managed Database engine version.
