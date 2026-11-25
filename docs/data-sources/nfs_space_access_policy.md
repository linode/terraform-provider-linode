---
page_title: "Linode: linode_nfs_space_access_policy"
description: |-
  Provides information about the access policy of a Linode NFS Space.
---

# Data Source: linode\_nfs\_space\_access\_policy

Provides information about the access policy of a Linode [NFS Space](nfs_space.html.markdown). The access policy controls which VPCs can reach the NFS Space and the optional mTLS configuration, and applies to all filesystems in the NFS Space. Every NFS Space has an access policy, which is created with the NFS Space and is disabled by default.
For more information, see the [Linode APIv4 docs](TBD).

~> **Beta** This data source is in beta and requires the provider `api_version` to be set to `v4beta`.

## Example Usage

The following example shows how one might use this data source to access information about the access policy of a Linode NFS Space.

```terraform
data "linode_nfs_space_access_policy" "test" {
    space_id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space this access policy belongs to.

## Attributes Reference

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

The Linode NFS Space access policy data source exports the following attributes:

* `id` - The unique ID of this access policy. Equal to the ID of the NFS Space.

* `label` - The descriptive label of this access policy.

* `enabled` - Whether this access policy is active. When `false`, the NFS Space is inaccessible.

* [`vpcs`](#vpcs) - (Read-Only Nested Attribute Set) The VPCs (and optionally specific subnets) permitted to access this NFS Space.

* `mtls_ca_cert` - The PEM-encoded CA certificate used to validate mTLS client certificates.

* `mtls_mode` - Whether mTLS client certificate presentation is enforced. (`required`, `optional`, `disabled`)

* [`vpc_acl`](#vpc_acl) - (Read-Only Nested Attribute List) The resolved VPC ACL entries, one per permitted VPC.

* `status` - The current lifecycle status of the access policy. (`creating`, `active`, `updating`, `deleting`)

* `created` - When this access policy was created.

* `updated` - When this access policy was last updated.

### vpcs

The following attributes are exported in each `vpcs` entry:

* `id` - The ID of the VPC.

* `subnets` - The IDs of the permitted subnets within the VPC. Null when all subnets in the VPC are permitted.

### vpc_acl

The following attributes are exported in each `vpc_acl` entry:

* `id` - The ID of the VPC.

* `label` - The label of the VPC.

* `url` - The relative URL of the VPC.

* `range` - The resolved CIDR range of the entire VPC. Only set when no specific subnets were requested.

* `subnets` - (Read-Only Nested Attribute List) The resolved subnet entries within the VPC. Only populated when specific subnets were requested.

  * `id` - The ID of the subnet.

  * `label` - The label of the subnet.

  * `url` - The relative URL of the subnet.

  * `range` - The resolved CIDR range of the subnet.
