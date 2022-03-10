---
layout: "linode"
page_title: "Linode: linode_database_engines"
sidebar_current: "docs-linode-datasource-database-engines"
description: |-
Provides information about Linode Managed Database engines.
---

# Data Source: linode\_database\_engines

Provides information about Linode Managed Database engines that match a set of filters.

## Example Usage

Get information about all Linode Managed Database engines:

```hcl
data "linode_database_engines" "all" {}
```

Get information about all Linode MySQL Database engines:

```hcl
data "linode_database_engines" "mysql" {
  filter {
    name = "engine"
    values = ["mysql"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `latest` - (Optional) If true, only the latest engine version will be returned.

* [`filter`](#filter) - (Optional) A set of filters used to select engines that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. (`version`)

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes

Each engine will be stored in the `engines` attribute and will export the following attributes:

* `engine` - The Managed Database engine type.

* `id` - The Managed Database engine ID in engine/version format.

* `version` - The Managed Database engine version.
