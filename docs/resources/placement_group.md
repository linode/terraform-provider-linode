---
page_title: "Linode: linode_placement_group"
description: |-
  Manages a Linode Placement Group.
---

# linode\_placement\_group

**NOTE: Placement Groups may not currently be available to all users.**

Manages a Linode Placement Group.

## Example Usage

Create a Placement Group with the local anti-affinity policy:

```terraform
resource "linode_placement_group" "test" {
    label = "my-placement-group"
    region = "us-mia"
    affinity_type = "anti_affinity:local"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - (Required) The region of the Placement Group.

* `affinity_type` - (Required) The affinity policy to use when placing Linodes in this group.

* `is_strict` - (Optional) Whether Linodes must be able to become compliant during assignment. (Default `true`)

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the Placement Group.

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's affinity policy.

* [`members`](#members) - A set of Linodes currently assigned to this Placement Group.

### Members

Represents a single Linode assigned to a Placement Group.

* `linode_id` - The ID of the Linode.

* `is_compliant` - Whether this Linode is currently compliant with the group's affinity policy.

## Import

Placement Groups be imported using their unique `id`, e.g.

```sh
terraform import linode_placement_group.mygroup 1234567
```
