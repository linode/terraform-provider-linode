---
page_title: "Linode: linode_nfs_filesystem_access_policy"
description: |-
  Manages the access policy of a Linode NFS Filesystem.
---

# linode\_nfs\_filesystem\_access\_policy

Manages the access policy of a Linode [NFS Filesystem](nfs_filesystem.html.markdown). The access policy controls which Linodes can mount the filesystem. Linodes must also be permitted by the parent [NFS Space access policy](nfs_space_access_policy.html.markdown).

~> **Beta** This resource is in beta and requires the provider `api_version` to be set to `v4beta`.

-> **Note:** A filesystem access policy can only be enabled while the parent NFS Space access policy is enabled. Reference the `linode_nfs_space_access_policy` resource (as in the example below) so that it is applied first.

-> **Note:** Every NFS Filesystem always has exactly one access policy. Destroying this resource only removes it from the Terraform state; the access policy remains configured on the NFS Filesystem.

## Example Usage

```terraform
resource "linode_nfs_space" "test" {
    label = "test-nfs-space"
}

resource "linode_nfs_filesystem" "test" {
    space_id           = linode_nfs_space.test.id
    label              = "test-filesystem"
    region             = "us-rno"
    max_capacity_bytes = 1099511627776
}

resource "linode_nfs_space_access_policy" "test" {
    space_id  = linode_nfs_space.test.id
    enabled   = true
    mtls_mode = "disabled"

    vpcs = [{ id = 123 }]

    depends_on = [linode_nfs_filesystem.test]
}

resource "linode_nfs_filesystem_access_policy" "test" {
    space_id      = linode_nfs_space_access_policy.test.space_id
    filesystem_id = linode_nfs_filesystem.test.id
    enabled       = true
    linode_ids    = [12345678]
    squash_policy = "root_squash"
    protocols     = ["nfsv4"]
}
```

## Argument Reference

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space the filesystem belongs to. Changing this value forces the creation of a new resource.

* `filesystem_id` - (Required) The ID of the NFS Filesystem this access policy belongs to. Changing this value forces the creation of a new resource.

* `enabled` - (Optional) Whether this access policy is active. When false, no Linodes can mount this filesystem.

* `inherit` - (Optional) Whether this filesystem falls back to the NFS Space access policy when no `linode_ids` are configured. When `false`, access is restricted to `linode_ids` only. Defaults to `true`.

* `linode_ids` - (Optional) The IDs of the Linodes permitted to mount this filesystem. Each Linode must belong to a VPC permitted by the NFS Space access policy.

* `squash_policy` - (Optional) The user-ID squash policy. (`none`, `root_squash`, `all_squash`)

* `protocols` - (Optional) The NFS protocol versions permitted. Defaults to the filesystem's `protocol_versions`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of this resource, in the format `space_id,filesystem_id`.

* `label` - The descriptive label of this access policy. Mirrors the label of the NFS Filesystem.

* [`linode_acl`](#linode_acl) - (Read-Only Nested Attribute List) The resolved Linode ACL entries permitted to mount this filesystem.

* `status` - The current lifecycle status of the access policy. (`creating`, `active`, `updating`, `deleting`)

* `created` - When this access policy was created.

* `updated` - When this access policy was last updated.

### linode_acl

The following attributes are exported in each `linode_acl` entry:

* `id` - The ID of the Linode.

* `label` - The label of the Linode.

* `url` - The relative URL of the Linode.

* `ipv6` - The resolved IPv6 address of the Linode.

## Import

NFS Filesystem access policies can be imported using the NFS Space ID and the NFS Filesystem ID separated by a comma, e.g.

```sh
terraform import linode_nfs_filesystem_access_policy.test 12345,67890
```
