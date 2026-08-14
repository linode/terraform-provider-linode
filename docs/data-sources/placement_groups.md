---
page_title: "Linode: linode_placement_groups"
description: |-
  Lists Linode Placement Groups on your account.
---

# Data Source: linode\_placement\_groups

Provides information about a list of Linode Placement Groups that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-placement-groups).

## Example Usage

The following example shows how one might use this data source to list Placement Groups.

```hcl
data "linode_placement_groups" "all" {}

data "linode_placement_groups" "filtered" {
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-pgs" {
  value = data.linode_placement_groups.all.placement_groups
}

output "filtered-pgs" {
  value = data.linode_placement_groups.filtered.placement_groups
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode Placement Groups that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Placement Group will be stored in the `placement_groups` attribute and will export the following attributes:

* `placement_groups` - (Nested Attribute List) The Placement Groups returned by this data source.

* `label` - The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - The region of the Placement Group.

* `placement_group_type` - The placement group type to use when placing Linodes in this group.

* `placement_group_policy` - Whether Linodes must be able to become compliant during assignment. (Default `strict`)

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's type.

* `members` - (Nested Attribute Set) A set of Linodes currently assigned to this Placement Group.

  * `linode_id` - The ID of the Linode.

  * `is_compliant` - Whether this Linode is currently compliant with the group's placement group type.

* `migrations` - (Nested Attribute) Any Linodes that are being migrated to or from the placement group. Referenced directly (e.g. `migrations.inbound`).

  * `inbound` - (Read-Only Object List) A list of the Linodes the system is migrating into the placement group. Referenced with an index (e.g. `inbound.0.linode_id`).

    * `linode_id` - The unique identifier for the Linode being migrated into the placement group.

  * `outbound` - (Read-Only Object List) A list of the Linodes the system is migrating out of the placement group. Referenced with an index (e.g. `outbound.0.linode_id`).

    * `linode_id` - The unique identifier for the Linode being migrated out of the placement group.

## Filterable Fields

* `id`

* `label`

* `region`

* `placement_group_type`

* `placement_group_policy`

* `is_compliant`
