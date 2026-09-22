---
page_title: "Linode: linode_vpc"
description: |-
  Manages a Linode VPC.
---

# linode\_vpc

Manages a Linode VPC.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-vpc).

Please refer to [linode_vpc_subnet](vpc_subnet.html.markdown) to manage the subnets under a Linode VPC.

## Example Usage

Create a VPC:

```terraform
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    description = "My first VPC."
}
```

Create a VPC with a `/52` IPv6 range prefix:

```terraform
# NOTE: IPv6 VPCs may not currently be available to all users.
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    
    ipv6 = [
      {
        range = "/52"
      }
    ]
}
```

Create a VPC with a custom IPv4 range:

```terraform
# NOTE: Custom VPC IPv4 Ranges may not currently be available to all users.
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    
    ipv4 = [
      {
        range = "10.0.0.0/8"
      }
    ]
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - (Required) The label of the VPC. This field can only contain ASCII letters, digits and dashes.

* `region` - (Required) The region of the VPC.

* `description` - (Optional) The user-defined description of this VPC.

* `vpc_type` - (Optional) The type of the VPC. Can be either `regular` or `rdma`. Defaults to `regular`. The `rdma` type creates an RDMA VPC and may not be available to all users. Changing this value forces the creation of a new VPC.

* [`ipv6`](#ipv6) - (Optional, Nested Attribute List) A list of IPv6 allocations under this VPC.

* [`ipv4`](#ipv4) - (Optional, Nested Attribute List) A list of IPv4 ranges under this VPC.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

Configures a single IPv6 range under this VPC.

* `range` - (Optional) An existing IPv6 prefix owned by the current account or a forward slash (/) followed by a valid prefix length. If unspecified, a range with the default prefix will be allocated for this VPC.

* `allocation_class` - (Optional) Indicates the labeled IPv6 Inventory that the VPC Prefix should be allocated from.

* `allocated_range` - (Read-Only) The value of range computed by the API. This is necessary when needing to access the range for an implicit allocation.

## IPv4

-> **Limited Availability** Custom VPC IPv4 Ranges may not currently be available to all users.

Configures a single IPv4 range under this VPC. Unlike IPv6, IPv4 ranges can be updated in-place without requiring resource replacement.

* `range` - (Required) The IPv4 range in CIDR format to assign to this VPC (e.g. `10.0.0.0/8`).

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was last updated.

* [`subnets`](#subnets) - A list of subnets under this VPC.

## Subnets

The following attributes are exported under each entry of the `subnets` field:

* `id` - The id of the VPC Subnet.

* `label` - The label of the VPC Subnet.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* `ipv6` - The IPv6 ranges of this subnet.

  * `range` - An IPv6 range allocated to this subnet.

* `linodes` - A list of Linodes assigned to this subnet.

  * `id` - ID of the Linode

  * `interfaces` - A list of networking interfaces objects.

    * `id` - ID of the interface.

    * `config_id` - ID of Linode Config that the interface is associated with. `null` for a Linode Interface.

    * `active` - Whether the Interface is actively in use.

* `databases` - A list of Managed Databases assigned to this subnet.

  * `id` - ID of a managed database assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the database.

  * `ipv6_ranges` - A list of IPv6 ranges assigned to the database.

    * `range` - An IPv6 address range in CIDR notation.

* `nodebalancers` - A list of NodeBalancers assigned to this subnet.

  * `id` - ID of a NodeBalancer assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the NodeBalancer.

  * `ipv6_ranges` - A list of IPv6 ranges assigned to the NodeBalancer.

    * `range` - An IPv6 address range in CIDR notation.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was last updated.
