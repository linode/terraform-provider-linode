---
page_title: "Linode: linode_nfs_quota"
description: |-
  Manages a Linode NFS Quota.
---

# linode\_nfs\_quota

Manages a Linode NFS Quota.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-nfs-quota). (**Note: v4beta only.**)

## Example Usage

Creating a basic NFS Quota:

```terraform
resource "linode_nfs_quota" "example" {
  space_id           = 123
  filesystem_id      = 456
  path               = "/exports/data"
  max_capacity_bytes = 2147483648
  max_file_count     = 2000
}
```

Creating an NFS Quota with user and group level limits:

```terraform
resource "linode_nfs_quota" "example" {
  space_id           = 123
  filesystem_id      = 456
  path               = "/exports/data"
  max_capacity_bytes = 4294967296
  max_file_count     = 4000

  user_group_config = {
    default_user_limit = {
      max_capacity_bytes = 1073741824
      max_file_count     = 1000
    }

    user_limits = [
      {
        identifier_type    = "uid"
        identifier         = "1001"
        max_capacity_bytes = 536870912
        max_file_count     = 500
      }
    ]

    user_group_limits = [
      {
        user = {
          identifier_type = "uid"
          identifier      = "1001"
        }
        group = {
          identifier_type = "gid"
          identifier      = "2001"
        }
        max_capacity_bytes = 268435456
        max_file_count     = 250
      }
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space this quota belongs to. Changing this value forces the creation of a new NFS Quota.

* `filesystem_id` - (Required) The ID of the NFS Filesystem this quota belongs to. Changing this value forces the creation of a new NFS Quota.

* `path` - (Required) The path within the filesystem that this quota applies to. Changing this value forces the creation of a new NFS Quota.

* `max_capacity_bytes` - (Required) The maximum capacity in bytes allowed by this quota.

* `max_file_count` - (Required) The maximum number of files allowed by this quota.

* [`user_group_config`](#user_group_config) - (Optional) User and group level quota configuration.

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

### user_group_config

The following arguments are supported in the `user_group_config` **Nested Attribute**:

* [`default_user_limit`](#capacity-limit) - (Optional) The default limit applied to each user.

* [`default_group_limit`](#capacity-limit) - (Optional) The default limit applied to each group.

* [`user_limits`](#identified-limit) - (Optional) A list of per-user quota limits.

* [`group_limits`](#identified-limit) - (Optional) A list of per-group quota limits.

* [`user_group_limits`](#user_group_limits) - (Optional) A list of combined user and group quota limits.

### capacity-limit

The `default_user_limit` and `default_group_limit` **Nested Attributes** support the following arguments:

* `max_capacity_bytes` - (Optional) The maximum capacity in bytes allowed by this limit.

* `max_file_count` - (Optional) The maximum number of files allowed by this limit.

### identified-limit

Each entry in `user_limits` and `group_limits` supports the following arguments:

* `identifier_type` - (Required) The identifier type. One of: `uid`, `username`, `gid`, `groupname`.

* `identifier` - (Required) The user or group identifier value.

* `max_capacity_bytes` - (Optional) The maximum capacity in bytes allowed for this identifier.

* `max_file_count` - (Optional) The maximum number of files allowed for this identifier.

### user_group_limits

Each entry in `user_group_limits` supports the following arguments:

* [`user`](#identity) - (Required) The user portion of the combined limit.

* [`group`](#identity) - (Required) The group portion of the combined limit.

* `max_capacity_bytes` - (Optional) The maximum capacity in bytes allowed for this combined limit.

* `max_file_count` - (Optional) The maximum number of files allowed for this combined limit.

### identity

The `user` and `group` **Nested Attributes** support the following arguments:

* `identifier_type` - (Required) The identifier type. One of: `uid`, `username`, `gid`, `groupname`.

* `identifier` - (Required) The user or group identifier value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of this NFS Quota.

* `is_root` - Whether this quota applies to the root of the filesystem.

* `used_capacity_bytes` - The capacity in bytes currently used under this quota.

* `quota_exceeded` - Whether the quota has been exceeded.

* `used_file_count` - The number of files currently used under this quota.

* `status` - The current status of this NFS Quota.

* `created` - When this NFS Quota was created.

* `updated` - When this NFS Quota was last updated.

* `collected_at` - When this NFS Quota's usage data was last collected.

## Import

NFS Quotas can be imported using the NFS Space ID, NFS Filesystem ID, and NFS Quota ID separated by commas, e.g.

```sh
terraform import linode_nfs_quota.example 123,456,789
```

-> **Note:** Each NFS Filesystem is created with a root quota on its mount path. The root quota cannot be created, but it can be imported to manage its `user_group_config`. Its `max_capacity_bytes` (determined by the filesystem size) and `max_file_count` are read-only and must match the values reported by the API. Destroying a root quota only removes it from Terraform state; the API removes it when the filesystem is deleted.

The root quota can be looked up with the [`linode_nfs_quotas`](../data-sources/nfs_quotas.md) data source and imported with an `import` block:

```terraform
data "linode_nfs_quotas" "root" {
  space_id      = linode_nfs_space.example.id
  filesystem_id = linode_nfs_filesystem.example.id

  filter = [
    {
      name   = "is_root"
      values = ["true"]
    }
  ]
}

locals {
  root_quota = data.linode_nfs_quotas.root.nfs_quotas[0]
}

import {
  to = linode_nfs_quota.root
  id = "${linode_nfs_space.example.id},${linode_nfs_filesystem.example.id},${local.root_quota.id}"
}

resource "linode_nfs_quota" "root" {
  space_id           = linode_nfs_space.example.id
  filesystem_id      = linode_nfs_filesystem.example.id
  path               = local.root_quota.path
  max_capacity_bytes = local.root_quota.max_capacity_bytes
  max_file_count     = local.root_quota.max_file_count

  user_group_config = {
    default_user_limit = {
      max_capacity_bytes = 1073741824
      max_file_count     = 1000
    }
  }
}
```

The filesystem must already exist before this configuration is planned, since the import ID must be known at plan time.
