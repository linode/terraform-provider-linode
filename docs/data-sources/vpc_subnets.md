---
page_title: "Linode: linode_vpc_subnets"
description: |-
  Lists all subnets under a Linode VPC.
---

# Data Source: linode\_vpc\_subnets

Provides information about a list of Linode VPC subnets that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpc-subnets).

## Example Usage

The following example shows how one might use this data source to list VPC subnets.

```hcl
data "linode_vpc_subnets" "filtered-subnets" {
    vpc_id = 123
    filter {
        name = "label"
        values = ["test"]
    }
}

output "vpc_subnets" {
  value = data.linode_vpc_subnets.filtered-subnets.vpc_subnets
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `vpc_id` - (Required) The id of the parent VPC for the list of VPCs.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode VPC subnets that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VPC subnet will be stored in the `vpc_subnets` attribute and will export the following attributes:

* `vpc_subnets` - (Nested Attribute List) The returned list of subnets under a VPC. Referenced by index (e.g. `vpc_subnets[0].id`).

* `id` - The unique id of the VPC subnet.

* `label` - The label of the VPC subnet.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* [`ipv6`](#ipv6) - (Nested Attribute List) A list of IPv6 ranges under this subnet.

* `linodes` - (Read-Only Object List) A list of Linodes added to this subnet. Referenced with an index (e.g. `linodes.0.id`).

  * `id` - ID of the Linode

  * `interfaces` - (Read-Only Object List) A list of networking interfaces objects. Referenced with an index (e.g. `interfaces.0.id`).

    * `id` - ID of the interface.

    * `config_id` - ID of Linode Config that the interface is associated with. `null` for a Linode Interface.

    * `active` - Whether the Interface is actively in use.

* `databases` - (Read-Only Object List) A list of Managed databases assigned to the VPC Subnet. Referenced with an index (e.g. `databases.0.id`).

  * `id` - ID of a managed database assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the database.

  * `ipv6_ranges` - (Read-Only Object List) A list of IPv6 ranges assigned to the database. Referenced with an index (e.g. `ipv6_ranges.0.range`).

    * `range` - An IPv6 address range in CIDR notation.

* `nodebalancers` - (Read-Only Object List) A list of NodeBalancers assigned to the VPC Subnet. Referenced with an index (e.g. `nodebalancers.0.id`).

  * `id` - ID of a NodeBalancer assigned to the VPC Subnet.

  * `ipv4_range` - IPv4 range assigned to the NodeBalancer.

  * `ipv6_ranges` - (Read-Only Object List) A list of IPv6 ranges assigned to the NodeBalancer. Referenced with an index (e.g. `ipv6_ranges.0.range`).

    * `range` - An IPv6 address range in CIDR notation.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was last updated.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

The following attributes are exported under each entry of the `ipv6` field:

* `range` - An IPv6 range allocated to this subnet in CIDR format.

## Filterable Fields

* `id`

* `label`

* `ipv4`
