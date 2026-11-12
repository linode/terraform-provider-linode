---
page_title: "Linode: linode_nfs_space"
description: |-
  Manages a Linode NFS Space.
---

# linode\_nfs\_space

Manages a Linode NFS Space. An NFS Space is the top-level logical isolation boundary for Akamai Cloud NFSaaS resources. NFS Filesystems are provisioned under an NFS Space.

~> **Beta** This resource is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

Create an NFS Space:

```terraform
resource "linode_nfs_space" "test" {
    label       = "test-nfs-space"
    description = "My first NFS Space."
    tags        = ["prod", "ml"]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The customer-provided label of the NFS Space. Must be unique within the account. 3-63 characters, lowercase alphanumeric and hyphens only.

* `description` - (Optional) An optional free-text description of this NFS Space.

* `tags` - (Optional) A set of customer-defined metadata tags for grouping and filtering.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The unique numeric ID of the NFS Space.

* `status` - The current lifecycle status of the NFS Space. (`creating`, `active`, `deleting`, `error`)

* `max_capacity_bytes` - The maximum total capacity available to this NFS Space across all of its filesystems, in bytes. Akamai-managed and read-only.

* `used_capacity_bytes` - The current used capacity across the NFS Space, in bytes. Akamai-managed and read-only.

* `created` - When this NFS Space was created.

* `updated` - When this NFS Space was last updated.

## Import

NFS Spaces can be imported using their numeric ID, e.g.

```sh
terraform import linode_nfs_space.test 12345
```
