---
page_title: "Linode: linode_vpc_subnet"
description: |-
  Manages a Linode VPC subnet.
---

# linode\_vpc\_subnet

Manages a Linode VPC subnet.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-vpc-subnet).

## Example Usage

Create a VPC subnet:

```terraform
resource "linode_vpc_subnet" "test" {
    vpc_id = 123
    label = "test-subnet"
    ipv4 = "10.0.0.0/24"
}
```

Create a VPC subnet with an implicitly determined IPv6 range:

```terraform
# NOTE: IPv6 VPCs may not currently be available to all users.
resource "linode_vpc_subnet" "test" {
    vpc_id = linode_vpc.test.id
    label = "test-subnet"
    ipv4 = "10.0.0.0/24"
  
    ipv6 = [
      {
        range = "auto"
      }
    ]
}

resource "linode_vpc" "test" {
  label = "test-vpc"
  region = "us-mia"
  
  ipv6 = [
    {
      range = "/52"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `vpc_id` - (Required) The id of the parent VPC for this VPC subnet.

* `label` - (Required) The label of the VPC. Only contains ASCII letters, digits and dashes.

* `ipv4` - (Required) The IPv4 range of this subnet in CIDR format.

* [`ipv6`](#ipv6) - (Optional, Nested Attribute List) A list of IPv6 ranges under this VPC subnet. NOTE: IPv6 VPCs may not currently be available to all users.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

The following arguments can be configured for each entry under the `ipv6` field:

* `range` - (Optional) An existing IPv6 prefix owned by the current account or a forward slash (/) followed by a valid prefix length. If `auto`, a range with the default prefix will be allocated for this VPC.

* `allocated_range` - (Read-Only) The value of range computed by the API. This is necessary when needing to access the range for an implicit allocation.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC Subnet.

* `vpc_type` - The type of the parent VPC (`regular` or `rdma`).

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

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was last updated.

## Import

Linode Virtual Private Cloud (VPC) Subnet can be imported using the `vpc_id` followed by the subnet `id` separated by a comma, e.g.

```sh
terraform import linode_vpc_subnet.my_subnet_duplicated 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
