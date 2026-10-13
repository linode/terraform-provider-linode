---
page_title: "Linode: linode_nfs_space"
description: |-
  Provides information about a Linode NFS Space.
---

# Data Source: linode\_nfs\_space

Provides information about a Linode NFS Space.

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

The following example shows how one might use this data source to access information about a Linode NFS Space.

```terraform
data "linode_nfs_space" "test" {
    id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique numeric ID of the NFS Space.

## Attributes Reference

The Linode NFS Space resource exports the following attributes:

* `label` - The customer-provided label of the NFS Space.

* `description` - An optional free-text description of this NFS Space.

* `status` - The current lifecycle status of the NFS Space.

* `max_capacity_bytes` - The maximum total capacity available to this NFS Space across all of its filesystems, in bytes.

* `used_capacity_bytes` - The current used capacity across the NFS Space, in bytes.

* `created` - When this NFS Space was created.

* `updated` - When this NFS Space was last updated.

* `tags` - A set of customer-defined metadata tags.
