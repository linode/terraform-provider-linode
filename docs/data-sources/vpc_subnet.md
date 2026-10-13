---
page_title: "Linode: linode_vpc_subnet"
description: |-
  Provides details about a subnet under a Linode VPC.
---

# Data Source: linode\_vpc\_subnet

Provides information about a Linode VPC subnet.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpc-subnet).

## Example Usage

The following example shows how one might use this data source to access information about a Linode VPC subnet.

```hcl
data "linode_vpc_subnet" "foo" {
    vpc_id = 123
    id = 12345
}

output "vpc_subnet" {
  value = data.linode_vpc_subnet.foo
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the parent VPC for this VPC Subnet.

* `id` - (Required) The unique id of this VPC subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - The label of the VPC subnet.

* `vpc_type` - The type of the parent VPC (`regular` or `rdma`). Omitted if the requesting account does not have access to the GPUDirect RDMA functionality.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* [`ipv6`](#ipv6) - (Nested Attribute List) A list of IPv6 ranges under this subnet.

* `linodes` - (Read-Only Object List) A list of Linodes added to this subnet. Referenced with an index (e.g. `linodes.0.id`).

  * `id` - ID of the Linode

  * `interfaces` - (Read-Only Object List) A list of networking interfaces objects. Referenced with an index (e.g. `linodes.0.interfaces.0.id`).

    * `id` - ID of the interface.

    * `config_id` - ID of Linode Config that the interface is associated with. `null` for a Linode Interface.

    * `active` - Whether the Interface is actively in use.

* `databases` - (Read-Only Object List) A list of Managed databases assigned to the VPC Subnet. Referenced with an index (e.g. `databases.0.id`).

  * `id` - ID of a managed database assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the database.

  * `ipv6_ranges` - (Read-Only Object List) A list of IPv6 ranges assigned to the database. Referenced with an index (e.g. `databases.0.ipv6_ranges.0.range`).

    * `range` - An IPv6 address range in CIDR notation.

* `nodebalancers` - (Read-Only Object List) A list of NodeBalancers assigned to the VPC Subnet. Referenced with an index (e.g. `nodebalancers.0.id`).

  * `id` - ID of a NodeBalancer assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the NodeBalancer.
  
  * `ipv6_ranges` - (Read-Only Object List) A list of IPv6 ranges assigned to the NodeBalancer. Referenced with an index (e.g. `nodebalancers.0.ipv6_ranges.0.range`).

    * `range` - An IPv6 address range in CIDR notation.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was last updated.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

The following attributes are exported under each entry of the `ipv6` field:

* `range` - An IPv6 range allocated to this subnet in CIDR format.
