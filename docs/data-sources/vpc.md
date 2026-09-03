---
page_title: "Linode: linode_vpc"
description: |-
  Provides details about a Linode VPC.
---

# Data Source: linode\_vpc

Provides information about a Linode VPC.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpc).

## Example Usage

The following example shows how one might use this data source to access information about a Linode VPC.

```hcl
data "linode_vpc" "foo" {
    id = 123
}

output "vpc" {
  value = data.linode_vpc.foo
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique id of this VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - The label of the VPC.

* `description` - The user-defined description of this VPC.

* `ipv6` - (Nested Attribute List) A list of IPv6 allocations under this VPC.

* `ipv4` - (Nested Attribute List) A list of IPv4 ranges under this VPC.

* `region` - The region where the VPC is deployed.

* `vpc_type` - The type of the VPC (`regular` or `rdma`). Omitted if the requesting account does not have access to the GPUDirect RDMA functionality.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was last updated.

* [`subnets`](#subnets) - A list of subnets under this VPC.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

Contains information about a single IPv6 allocation under this VPC.

* `range` - The allocated range in CIDR format.

## IPv4

-> **Limited Availability** Custom VPC IPv4 Ranges may not currently be available to all users.

Contains information about a single IPv4 range under this VPC.

* `range` - The IPv4 range in CIDR format.

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

### Subnets data source

The `subnets` list in this resource requires an additional refresh after the initial apply before newly created subnets appear because all subnets are created as resources after the vpc resource is created. To list all subnets under a VPC with immediate availability after apply, use the [linode_vpc_subnets](vpc_subnets.html.markdown) data source with Terraform [`depends_on`](https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on).
