---
page_title: "Linode: linode_locks"
description: |-
  Provides information about Linode Locks that match a set of filters.
---

# Data Source: linode\_locks

~> **Early Access:** Lock functionality is in early access and may not be available to all users.

~> **Important** Only unrestricted users can view locks. Restricted users cannot access lock information even if they have permissions for the resources.

Provides information about Linode Locks that match a set of filters. Locks prevent accidental deletion, rebuild operations, and service transfers of resources.

For more information, see the Linode APIv4 docs (TBD).

## Example Usage

Get information about all Linode Locks with a certain entity type:

```hcl
data "linode_locks" "linode_locks" {
  filter {
    name = "entity_type"
    values = ["linode"]
  }
}

output "lock_ids" {
  value = data.linode_locks.linode_locks.locks.*.id
}
```

Get information about a specific lock by entity ID:

```hcl
data "linode_locks" "my_instance_locks" {
  filter {
    name = "entity_id"
    values = ["12345"]
  }
}
```

Get information about all locks:

```hcl
data "linode_locks" "all" {}

output "all_lock_ids" {
  value = data.linode_locks.all.locks.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Locks that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode lock will be stored in the `locks` attribute and will export the following attributes:

* `id` - The unique ID of the Lock.

* `entity_id` - The ID of the locked entity.

* `entity_type` - The type of the locked entity.

* `lock_type` - The type of lock.

* `entity_label` - The label of the locked entity.

* `entity_url` - The URL of the locked entity.

## Filterable Fields

* `id`

* `entity_id`

* `entity_type`

* `entity_label`

* `lock_type`
