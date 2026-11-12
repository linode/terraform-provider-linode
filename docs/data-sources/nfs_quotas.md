---
page_title: "Linode: linode_nfs_quotas"
description: |-
  Provides information about a list of Linode NFS Quotas.
---

# Data Source: linode\_nfs\_quotas

Provides information about a list of Linode NFS Quotas under an [NFS Filesystem](nfs_filesystem.md) that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-nfs-quotas).

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

Get the root quota of an NFS Filesystem:

```terraform
data "linode_nfs_quotas" "root" {
    space_id      = 123
    filesystem_id = 456

    filter = [
        {
            name   = "is_root"
            values = ["true"]
        }
    ]
}

output "root_quota_id" {
    value = data.linode_nfs_quotas.root.nfs_quotas.0.id
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space the filesystem belongs to.

* `filesystem_id` - (Required) The ID of the NFS Filesystem to list quotas for.

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Nested Attribute Set) A set of filters used to select NFS Quotas that meet certain requirements.

## Filterable Fields

* `id`

* `path`

* `is_root`

* `status`

* `quota_exceeded`

* `max_capacity_bytes`

* `max_file_count`

## Attributes Reference

Each NFS Quota will be stored in the `nfs_quotas` attribute and will export the following attributes:

* `id` - The unique ID of this NFS Quota.

* `space_id` - The ID of the NFS Space this quota belongs to.

* `filesystem_id` - The ID of the NFS Filesystem this quota belongs to.

* `path` - The path within the filesystem that this quota applies to.

* `max_capacity_bytes` - The maximum capacity in bytes allowed by this quota.

* `max_file_count` - The maximum number of files allowed by this quota.

* `is_root` - Whether this quota applies to the root of the filesystem.

* `used_capacity_bytes` - The capacity in bytes currently used under this quota.

* `quota_exceeded` - Whether the quota has been exceeded.

* `used_file_count` - The number of files currently used under this quota.

* `status` - The current status of this NFS Quota.

* `created` - When this NFS Quota was created.

* `updated` - When this NFS Quota was last updated.

* `collected_at` - When this NFS Quota's usage data was last collected.

* [`user_group_config`](#user_group_config) - User and group level quota configuration.

### user_group_config

The `user_group_config` **Nested Attribute** exports the following attributes:

* `default_user_limit` - The default limit applied to each user. See [capacity limit](#capacity-limit) for details.

* `default_group_limit` - The default limit applied to each group. See [capacity limit](#capacity-limit) for details.

* `user_limits` - A list of per-user quota limits. See [identified limit](#identified-limit) for details.

* `group_limits` - A list of per-group quota limits. See [identified limit](#identified-limit) for details.

* `user_group_limits` - A list of combined user and group quota limits. See [user group limits](#user_group_limits) for details.

### capacity-limit

* `max_capacity_bytes` - The maximum capacity in bytes allowed by this limit.

* `max_file_count` - The maximum number of files allowed by this limit.

### identified-limit

* `identifier_type` - The identifier type. One of: `uid`, `username`, `gid`, `groupname`.

* `identifier` - The user or group identifier value.

* `max_capacity_bytes` - The maximum capacity in bytes allowed for this identifier.

* `max_file_count` - The maximum number of files allowed for this identifier.

### user_group_limits

* `user` - The user portion of the combined limit. See [identity](#identity) for details.

* `group` - The group portion of the combined limit. See [identity](#identity) for details.

* `max_capacity_bytes` - The maximum capacity in bytes allowed for this combined limit.

* `max_file_count` - The maximum number of files allowed for this combined limit.

### identity

* `identifier_type` - The identifier type. One of: `uid`, `username`, `gid`, `groupname`.

* `identifier` - The user or group identifier value.

## Filter

Each object in the `filter` nested attribute set supports the following arguments:

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)
