---
page_title: "Linode: linode_placement_group_assignment"
description: |-
  Manages a single assignment between a Linode and a Placement Group.
---

# linode\_placement\_group\_assignment

Manages a single assignment between a Linode and a Placement Group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-group-add-linode).

**NOTE: Placement Groups may not currently be available to all users.**

To prevent update conflicts, Linodes managed through the `linode_instance` resource should specify `placement_group_externally_managed`:

```terraform
resource "linode_instance" "my-instance" {
  # ...
  placement_group_externally_managed = true
}
```

## Example Usage

```terraform
resource "linode_placement_group_assignment" "my-assignment" {
  placement_group_id = linode_placement_group.my-pg.id
  linode_id = linode_instance.my-inst.id
}

resource "linode_placement_group" "my-pg" {
  label = "my-pg"
  region = "us-east"
  placement_group_type = "anti_affinity:local"
}

resource "linode_instance" "my-inst" {
  label = "my-inst"
  region = "us-east"
  type = "g6-nanode-1"
  placement_group_externally_managed = true
}
```

## Argument Reference

The following arguments are supported:

* `placement_group_id` - (Required) The unique ID of the target Placement Group.

* `linode_id` - (Required) The unique ID of the Linode to assign.

## Import

Placement Group assignments can be imported using the Placement Group's ID followed by the Linode's ID separated by a comma, e.g.

```sh
terraform import linode_placement_group_assignment.my-assignment 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
