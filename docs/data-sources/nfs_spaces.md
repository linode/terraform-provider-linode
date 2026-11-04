---
page_title: "Linode: linode_nfs_spaces"
description: |-
  Provides information about a list of Linode NFS Spaces.
---

# Data Source: linode\_nfs\_spaces

Provides information about a list of Linode NFS Spaces that match a set of filters.

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

Get information about all NFS Spaces with a given label:

```terraform
data "linode_nfs_spaces" "filtered" {
    filter = [
        {
            name   = "label"
            values = ["prod-ml-space"]
        }
    ]
}

output "nfs_space_id" {
    value = data.linode_nfs_spaces.filtered.nfs_spaces.0.id
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Nested Attribute Set) A set of filters used to select NFS Spaces that meet certain requirements.

## Filterable Fields

* `id`

* `label`

* `description`

* `status`

* `max_capacity_bytes`

* `used_capacity_bytes`

## Attributes Reference

Each NFS Space will be stored in the `nfs_spaces` attribute and will export the following attributes:

* `id` - The unique numeric ID of the NFS Space.

* `label` - The customer-provided label of the NFS Space.

* `description` - An optional free-text description of this NFS Space.

* `status` - The current lifecycle status of the NFS Space.

* `max_capacity_bytes` - The maximum total capacity available to this NFS Space across all of its filesystems, in bytes.

* `used_capacity_bytes` - The current used capacity across the NFS Space, in bytes.

* `created` - When this NFS Space was created.

* `updated` - When this NFS Space was last updated.

* `tags` - A set of customer-defined metadata tags.

## Filter

Each object in the `filter` nested attribute set supports the following arguments:

* `name` - (Required) The name of the field to filter by.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)
