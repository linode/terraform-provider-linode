---
page_title: "Linode: linode_lock"
description: |-
  Manages a Linode Lock to prevent accidental deletion of resources.
---

# linode\_lock

~> **Early Access** Locks are in Early Access and may not be available to all users.

~> **Important** Only unrestricted users can create and delete locks. Restricted users cannot manage locks even if they have read/write permissions for the resource.

Manages a Linode Lock which prevents accidental deletion and modification of resources. Locks protect against deletion, rebuild operations, and service transfers. The `cannot_delete_with_subresources` lock type also protects subresources such as disks, configs, interfaces, and IP addresses.

For more information, see the Linode APIv4 docs (TBD).

-> **Note** Only one lock can exist per resource at a time. You cannot have both `cannot_delete` and `cannot_delete_with_subresources` locks on the same resource simultaneously.

## Example Usage

Create a basic lock that prevents a Linode from being deleted:

```terraform
resource "linode_lock" "my-lock" {
  entity_id   = linode_instance.my-inst.id
  entity_type = "linode"
  lock_type   = "cannot_delete"
}

resource "linode_instance" "my-inst" {
  label  = "my-inst"
  region = "us-east"
  type   = "g6-nanode-1"
}
```

Create a lock that prevents a Linode and its subresources (disks, configs, etc.) from being deleted:

```terraform
resource "linode_lock" "my-lock" {
  entity_id   = linode_instance.my-inst.id
  entity_type = "linode"
  lock_type   = "cannot_delete_with_subresources"
}

resource "linode_instance" "my-inst" {
  label  = "my-inst"
  region = "us-east"
  type   = "g6-nanode-1"
}
```

## Argument Reference

The following arguments are supported:

* `entity_id` - (Required) The ID of the entity to lock.

* `entity_type` - (Required) The type of the entity to lock. Currently only `linode` is supported. Note: Linodes that are part of an LKE cluster cannot be locked.

* `lock_type` - (Required) The type of lock to apply. Only one lock type can exist per resource at a time. Valid values are:
  * `cannot_delete` - Prevents the resource from being deleted, rebuilt, or transferred to another account.
  * `cannot_delete_with_subresources` - Prevents the resource from being deleted, rebuilt, or transferred, and also prevents deletion of subresources (disks, configs, interfaces, and IP addresses).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the Lock.

* `entity_label` - The label of the locked entity.

* `entity_url` - The URL of the locked entity.

## Import

Locks can be imported using the Lock's ID, e.g.

```sh
terraform import linode_lock.my-lock 1234567
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
