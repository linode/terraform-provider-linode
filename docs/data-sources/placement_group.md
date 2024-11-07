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

* `label` - The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - The region of the Placement Group.

* `placement_group_type` - The placement group type to use when placing Linodes in this group.

* `placement_group_policy` - Whether Linodes must be able to become compliant during assignment. (Default `strict`)

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's placement group type.

* `members` - A set of Linodes currently assigned to this Placement Group.

### Members

Represents a single Linode assigned to a Placement Group.

* `linode_id` - The ID of the Linode.

* `is_compliant` - Whether this Linode is currently compliant with the group's placement group type.
