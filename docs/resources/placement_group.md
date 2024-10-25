---
page_title: "Linode: linode_placement_group"
description: |-
  Manages a Linode Placement Group.
---

# linode\_placement\_group

Manages a Linode Placement Group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-placement-group).

## Example Usage

Create a Placement Group with the local anti-affinity policy:

```terraform
resource "linode_placement_group" "test" {
    label = "my-placement-group"
    region = "us-mia"
    placement_group_type = "anti_affinity:local"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - (Required) The region of the Placement Group.

* `placement_group_type` - (Required) The placement group type to use when placing Linodes in this group.

* `placement_group_policy` - (Optional) Whether Linodes must be able to become compliant during assignment. (Default `strict`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the Placement Group.

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's placement group type.

* [`members`](#members) - A set of Linodes currently assigned to this Placement Group.

### Members

Represents a single Linode assigned to a Placement Group.

* `linode_id` - The ID of the Linode.

* `is_compliant` - Whether this Linode is currently compliant with the group's placement group type.

## Import

Placement Groups be imported using their unique `id`, e.g.

```sh
terraform import linode_placement_group.mygroup 1234567
```
