---
page_title: "Linode: linode_nfs_filesystem_access_policy"
description: |-
  Provides information about the access policy of a Linode NFS Filesystem.
---

# Data Source: linode\_nfs\_filesystem\_access\_policy

Provides information about the access policy of a Linode [NFS Filesystem](nfs_filesystem.html.markdown). The access policy controls which Linodes can mount the filesystem. Linodes must also be permitted by the parent [NFS Space access policy](nfs_space_access_policy.html.markdown).
For more information, see the [Linode APIv4 docs](TBD).

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

-> **Note:** An NFS Filesystem has no access policy until one is first configured; until then, the filesystem follows the NFS Space access policy. In that case, this data source sets `configured` to `false` and all other policy attributes to null. If the same configuration manages the access policy with the `linode_nfs_filesystem_access_policy` resource, reference that resource (as in the example below) so that the data source is read after the policy is configured.

## Example Usage

The following example shows how one might use this data source to access information about the access policy of a Linode NFS Filesystem.

```terraform
data "linode_nfs_filesystem_access_policy" "test" {
    space_id      = 12345
    filesystem_id = 67890
}
```

The following example reads an access policy managed in the same configuration:

```terraform
data "linode_nfs_filesystem_access_policy" "test" {
    space_id      = linode_nfs_filesystem_access_policy.test.space_id
    filesystem_id = linode_nfs_filesystem_access_policy.test.filesystem_id
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space the filesystem belongs to.

* `filesystem_id` - (Required) The ID of the NFS Filesystem this access policy belongs to.

## Attributes Reference

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

The Linode NFS Filesystem access policy data source exports the following attributes:

* `id` - The unique ID of this access policy, in the format `space_id,filesystem_id`.

* `configured` - Whether the access policy has been configured. When `false`, the filesystem follows the NFS Space access policy and all of the following attributes are null.

* `label` - The descriptive label of this access policy. Mirrors the label of the NFS Filesystem.

* `enabled` - Whether this access policy is active. When `false`, no Linodes can mount this filesystem.

* `inherit` - Whether this filesystem falls back to the NFS Space access policy when no `linode_ids` are configured. When `false`, access is restricted to `linode_ids` only.

* `linode_ids` - The IDs of the Linodes permitted to mount this filesystem.

* `squash_policy` - The user-ID squash policy. (`none`, `root_squash`, `all_squash`)

* `protocols` - The NFS protocol versions permitted.

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
