---
page_title: "Linode: linode_nfs_filesystem"
description: |-
  Provides information about a Linode NFS Filesystem.
---

# Data Source: linode\_nfs\_filesystem

Provides information about a Linode NFS Filesystem.

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

The following example shows how one might use this data source to access information about a Linode NFS Filesystem.

```terraform
data "linode_nfs_filesystem" "test" {
    id = 67890
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique numeric ID of the NFS Filesystem.

## Attributes Reference

The Linode NFS Filesystem resource exports the following attributes:

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

The following attributes are exported under the `stats` attribute:

* `used_capacity_bytes` - Capacity currently in use, in bytes.

* `used_file_count` - Number of files and directories currently in use.

* `collected_at` - When these statistics were collected.
