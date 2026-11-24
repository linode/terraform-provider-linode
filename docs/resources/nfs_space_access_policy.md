---
page_title: "Linode: linode_nfs_space_access_policy"
description: |-
  Manages the access policy of a Linode NFS Space.
---

# linode\_nfs\_space\_access\_policy

Manages the access policy of a Linode [NFS Space](nfs_space.html.markdown). The access policy controls which VPCs can reach the NFS Space and the optional mTLS configuration, and applies to all filesystems in the NFS Space.
For more information, see the [Linode APIv4 docs](TBD).

~> **Beta** This resource is in beta and requires the provider `api_version` to be set to `v4beta`.

-> **Note:** The API only accepts access policy updates once the NFS Space has an active tenant, which is provisioned along with the first [NFS Filesystem](nfs_filesystem.html.markdown) in the space. Add a `depends_on` on a filesystem in the space, as in the example below. Updates are retried while the tenant is still being provisioned.

-> **Note:** Every NFS Space always has exactly one access policy. Destroying this resource only removes it from the Terraform state; the access policy remains configured on the NFS Space.

## Example Usage

```terraform
resource "linode_nfs_space" "test" {
    label = "test-nfs-space"
}

resource "linode_nfs_filesystem" "test" {
    space_id           = linode_nfs_space.test.id
    label              = "test-filesystem"
    region             = "us-rno"
    max_capacity_bytes = 1099511627776
}

resource "linode_nfs_space_access_policy" "test" {
    space_id  = linode_nfs_space.test.id
    enabled   = true
    mtls_mode = "disabled"

    vpcs = [
        {
            id      = 123
            subnets = [456]
        },
    ]

    depends_on = [linode_nfs_filesystem.test]
}
```

## Argument Reference

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

The following arguments are supported:

* `space_id` - (Required) The ID of the NFS Space this access policy belongs to. Changing this value forces the creation of a new resource.

* `enabled` - (Optional) Whether this access policy is active.

* [`vpcs`](#vpcs) - (Optional, Nested Attribute Set) The VPCs (and optionally specific subnets) permitted to access this NFS Space. Set `vpcs = []` to remove all VPCs; omitting the argument leaves the current allowlist unchanged.

* `mtls_ca_cert` - (Optional) The PEM-encoded CA certificate used to validate mTLS client certificates.

* `mtls_mode` - (Optional) Whether mTLS client certificate presentation is enforced. (`required`, `optional`, `disabled`)

### vpcs

The following arguments are supported in each `vpcs` entry:

* `id` - (Required) The ID of the VPC.

* `subnets` - (Optional) The IDs of the subnets within the VPC to permit. When unset, all subnets in the VPC are permitted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of this resource. Equal to the ID of the NFS Space.

* `label` - The descriptive label of this access policy.

* [`vpc_acl`](#vpc_acl) - (Read-Only Nested Attribute List) The resolved VPC ACL entries, one per permitted VPC.

* `status` - The current lifecycle status of the access policy. (`creating`, `active`, `updating`, `deleting`)

* `created` - When this access policy was created.

* `updated` - When this access policy was last updated.

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

## Import

NFS Space access policies can be imported using the numeric ID of the NFS Space, e.g.

```sh
terraform import linode_nfs_space_access_policy.test 12345
```
