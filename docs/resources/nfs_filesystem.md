---
page_title: "Linode: linode_nfs_filesystem"
description: |-
  Manages a Linode NFS Filesystem.
---

# linode\_nfs\_filesystem

Manages a Linode NFS Filesystem. An NFS Filesystem is a managed NFS export provisioned under an [NFS Space](nfs_space.html.markdown). Each filesystem is region-specific.

~> **Beta** This resource is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

Create an NFS Filesystem under an NFS Space:

```terraform
resource "linode_nfs_space" "test" {
    label = "test-nfs-space"
}

resource "linode_nfs_filesystem" "test" {
    space_id           = linode_nfs_space.test.id
    label              = "test-filesystem"
    region             = "us-rno"
    max_capacity_bytes = 1099511627776
    tags               = ["training", "ml"]
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space this filesystem belongs to. Changing this value forces the creation of a new NFS Filesystem.

* `label` - (Required) The customer-provided label of the NFS Filesystem. Unique within the NFS Space. 3-63 characters, lowercase alphanumeric and hyphens only.

* `region` - (Required) The region where this filesystem is provisioned. Changing this value forces the creation of a new NFS Filesystem.

* `max_capacity_bytes` - (Required) The provisioned capacity of the filesystem, in bytes. Minimum `1073741824` (1 GiB).

* `protocol_versions` - (Optional) The NFS protocol versions supported by this filesystem. Currently only `nfsv4` is supported. Changing this value forces the creation of a new NFS Filesystem.

* `tags` - (Optional) A set of customer-defined metadata tags.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The unique numeric ID of the NFS Filesystem.

* `status` - The current lifecycle status of the filesystem. (`creating`, `active`, `updating`, `deleting`, `error`)

* `mount_target_ips` - A list of NFS mount strings using NFS endpoint IPv6 addresses.

* `mount_target_fqdn` - The full NFS mount string using the DNS-resolvable FQDN.

* `snapshot_usage_bytes` - The total unique snapshot capacity in bytes for this filesystem.

* `ldap_config_id` - The ID of the LDAP Config associated with this filesystem's parent NFS Space, if any.

* `source_snapshot_id` - The source snapshot ID if this filesystem was created via snapshot clone.

* `created` - When this filesystem was created.

* `updated` - When this filesystem was last updated.

* [`stats`](#stats) - Live usage statistics for this filesystem.

## stats

The following attributes are exported under the `stats` attribute:

* `used_capacity_bytes` - Capacity currently in use, in bytes.

* `used_file_count` - Number of files and directories currently in use.

* `collected_at` - When these statistics were collected.

## Import

NFS Filesystems can be imported using their numeric ID, e.g.

```sh
terraform import linode_nfs_filesystem.test 67890
```
