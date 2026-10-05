---
page_title: "Linode: linode_placement_group"
description: |-
  Provides details about a Linode placement group.
---

# Data Source: linode\_placement\_group

`linode_placement_group` provides details about a Linode placement group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-placement-group).

## Example Usage

The following example shows how the resource might be used to obtain additional information about a Linode placement group.

```hcl
data "linode_placement_group" "pg" {
  id = 12345
}
```

## Argument Reference

* `id` - (Required) The ID of the Placement Group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - The region of the Placement Group.

* `placement_group_type` - The placement group type to use when placing Linodes in this group.

* `placement_group_policy` - Whether Linodes must be able to become compliant during assignment. (Default `strict`)

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's placement group type.

* `members` - (Nested Attribute Set) A set of Linodes currently assigned to this Placement Group. Set elements can't be referenced by index; use a `for` expression or `tolist(...)` to access them.

* `migrations` - (Nested Attribute) Any Linodes that are being migrated to or from the placement group. Referenced directly (e.g. `migrations.inbound`).

  * `inbound` - (Read-Only Object List) A list of the Linodes the system is migrating into the placement group. Referenced with an index (e.g. `migrations.inbound.0.linode_id`).

    * `linode_id` - The unique identifier for the Linode being migrated into the placement group.

  * `outbound` - (Read-Only Object List) A list of the Linodes the system is migrating out of the placement group. Referenced with an index (e.g. `migrations.outbound.0.linode_id`).

    * `linode_id` - The unique identifier for the Linode being migrated out of the placement group.

### Members

Represents a single Linode assigned to a Placement Group.

* `linode_id` - The ID of the Linode.

* `is_compliant` - Whether this Linode is currently compliant with the group's placement group type.
