---
page_title: "Linode: linode_nfs_filesystems"
description: |-
  Provides information about a list of Linode NFS Filesystems.
---

# Data Source: linode\_nfs\_filesystems

Provides information about a list of Linode NFS Filesystems under an [NFS Space](nfs_space.html.markdown) that match a set of filters.

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

Get information about all NFS Filesystems in an NFS Space with a given label:

```terraform
data "linode_nfs_filesystems" "filtered" {
    space_id = 12345

    filter {
        name   = "label"
        values = ["training-data"]
    }
}

output "nfs_filesystem_id" {
    value = data.linode_nfs_filesystems.filtered.nfs_filesystems.0.id
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space to list filesystems for.

* [`filter`](#filter) - (Optional) A set of filters used to select NFS Filesystems that meet certain requirements.

## Filterable Fields

* `id`

* `label`

* `region`

* `status`

* `max_capacity_bytes`

## Attributes Reference

Each NFS Filesystem will be stored in the `nfs_filesystems` attribute and will export the following attributes:

* `id` - The unique numeric ID of the NFS Filesystem.

* `space_id` - The ID of the NFS Space this filesystem belongs to.

* `label` - The customer-provided label of the NFS Filesystem.

* `region` - The region where this filesystem is provisioned.

* `max_capacity_bytes` - The provisioned capacity of the filesystem, in bytes.

* `protocol_versions` - The NFS protocol versions supported by this filesystem.

* `status` - The current lifecycle status of the filesystem.

* `mount_target_ips` - A list of NFS mount strings using NFS endpoint IPv6 addresses.

* `mount_target_fqdn` - The full NFS mount string using the DNS-resolvable FQDN.

* `snapshot_usage_bytes` - The total unique snapshot capacity in bytes for this filesystem.

* `ldap_config_id` - The ID of the LDAP Config associated with this filesystem's parent NFS Space, if any.

* `source_snapshot_id` - The source snapshot ID if this filesystem was created via snapshot clone.

* `created` - When this filesystem was created.

* `updated` - When this filesystem was last updated.

* `tags` - A set of customer-defined metadata tags.

* [`stats`](#stats) - Live usage statistics for this filesystem.

## stats

The following attributes are exported under each `stats` attribute:

* `used_capacity_bytes` - Capacity currently in use, in bytes.

* `used_file_count` - Number of files and directories currently in use.

* `collected_at` - When these statistics were collected.

## Filter

A `filter` block supports the following arguments:

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.
