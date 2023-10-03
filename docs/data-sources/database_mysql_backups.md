---
page_title: "Linode: linode_database_mysql_backups"
description: |-
  Provides information about Linode MySQL Database Backups that match a set of filters.
---

# Data Source: linode\_database\_mysql\_backups

~> **NOTICE:** This data source has been deprecated in favor of `linode_database_backups`.

Provides information about Linode MySQL Database Backups that match a set of filters.

## Example Usage

Get information about all backups for a MySQL database:

```hcl
data "linode_database_mysql_backups" "all-backups" {
  database_id = 12345
}
```

Get information about all automatic MySQL Database Backups:

```hcl
data "linode_database_mysql_backups" "auto-backups" {
  database_id = 12345
  
  filter {
    name = "type"
    values = ["auto"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_id` - (Required) The ID of the database to retrieve backups for.

* `latest` - (Optional) If true, only the latest backup will be returned.

* [`filter`](#filter) - (Optional) A set of filters used to select database backups that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. (`created`)

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each backup will be stored in the `backups` attribute and will export the following attributes:

* `created` - A time value given in a combined date and time format that represents when the database backup was created.

* `id` - The ID of the database backup object.

* `label` - The database backup’s label, for display purposes only.

* `type` - The type of database backup, determined by how the backup was created.
