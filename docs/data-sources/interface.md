---
page_title: "Linode: linode_interface"
description: |-
  Provides details about a Linode interface.
---

# Data Source: linode\_interface

Provides information about a Linode interface.

This data source is specifically for Linode interfaces. For legacy config interfaces, see the `linode_instance_config` resource.

For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-linode-interface).

## Example Usage

The following example shows how one might use this data source to access information about a Linode interface.

```hcl
data "linode_interface" "foo" {
  id        = 67890
  linode_id = 12345
}

output "interface" {
  value = data.linode_interface.foo
}
```

### VPC Interface NAT Gateway Example

The following example reads a VPC interface and outputs the NAT gateway serving its IPv4 configuration.

```hcl
data "linode_interface" "vpc" {
  id        = linode_interface.vpc.id
  linode_id = linode_instance.my-instance.id
}

output "natgateway_addresses" {
  value = data.linode_interface.vpc.vpc.ipv4.natgateway.addresses
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Linode interface.

* `linode_id` - (Required) The ID of the Linode the interface is attached to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `default_route` - (Nested Attribute) Indicates if the interface serves as the default route. Referenced directly (e.g. `default_route.ipv4`).

  * `ipv4` - Whether this interface is used for the IPv4 default route.

  * `ipv6` - Whether this interface is used for the IPv6 default route.

* `public` - (Nested Attribute) The public interface configuration. Referenced directly (e.g. `public.ipv4.addresses`).

  * `ipv4` - (Nested Attribute) IPv4 configuration.

    * `addresses` - (Nested Attribute Set) IPv4 addresses assigned to this interface. Each object supports:

      * `address` - The IPv4 address.

      * `primary` - Whether this is the primary IPv4 address.

    * `shared` - (Nested Attribute Set) IPv4 addresses shared with other Linodes. Each object supports:

      * `address` - The shared IPv4 address.

      * `linode_id` - The ID of the Linode that this shared address belongs to.

  * `ipv6` - (Nested Attribute) IPv6 configuration.

    * `ranges` - (Nested Attribute Set) IPv6 ranges assigned to this interface. Each object supports:

      * `range` - The IPv6 range.

      * `route_target` - The route target for this IPv6 range.

    * `shared` - (Nested Attribute Set) IPv6 ranges shared with other Linodes. Each object supports the same fields as `ranges`.

    * `slaac` - (Nested Attribute Set) IPv6 SLAAC configuration. Each object supports:

      * `address` - The IPv6 SLAAC address.

      * `prefix` - The prefix length for the IPv6 SLAAC address.

* `vlan` - (Nested Attribute) The VLAN interface configuration. Referenced directly (e.g. `vlan.vlan_label`).

  * `vlan_label` - The label of the VLAN.

  * `ipam_address` - The IPAM (IP Address Management) address of the VLAN interface.

* `vpc` - (Nested Attribute) The VPC interface configuration. Referenced directly (e.g. `vpc.subnet_id`).

  * `subnet_id` - The ID of the VPC subnet.

  * `ipv4` - (Nested Attribute) IPv4 configuration for the VPC interface.

    * `addresses` - (Nested Attribute Set) IPv4 addresses assigned to this VPC interface. Each object supports:

      * `address` - The VPC IPv4 address.

      * `primary` - Whether this is the primary VPC IPv4 address.

      * `nat_1_1_address` - The 1:1 NAT address for this VPC IPv4 address.

    * `ranges` - (Nested Attribute Set) IPv4 ranges assigned to this VPC interface. Each object supports:

      * `range` - The VPC IPv4 range.

    * `natgateway` - (Nested Attribute) The NAT gateway serving this VPC interface's IPv4 configuration. Null if none is assigned. Referenced directly (e.g. `vpc.ipv4.natgateway.id`).

      * `id` - The ID of the NAT gateway.

      * `label` - The label of the NAT gateway.

      * `addresses` - (List of Strings) The public IPv4 addresses of the NAT gateway used by this interface.

      * `portset_assignments` - The number of portsets currently assigned to the NAT gateway.

      * `portset_capacity` - The maximum number of portsets the NAT gateway can assign.

      * `portsets` - (Nested Attribute List) The portsets assigned to the NAT gateway. Each object supports:

        * `address` - The public IP address of the portset.

        * `ports` - (Nested Attribute List) The port ranges assigned to the portset. Each object supports:

          * `start` - The first port in the range.

          * `end` - The last port in the range.

  * `ipv6` - (Nested Attribute) IPv6 configuration for the VPC interface.

    * `is_public` - Indicates whether the IPv6 configuration on the interface is public.

    * `slaac` - (Nested Attribute Set) IPv6 SLAAC address ranges. Each object supports:

      * `range` - The IPv6 network range in CIDR notation.

      * `address` - The assigned IPv6 address within the range.

    * `ranges` - (Nested Attribute Set) IPv6 ranges assigned to this VPC interface. Each object supports:

      * `range` - The IPv6 network range in CIDR notation.

* `rdma_vpc` - (Nested Attribute) The RDMA VPC interface configuration. **NOTE: RDMA VPC interfaces may not currently be available to all users.**

  * `vpc_id` - The ID of the parent RDMA VPC.

  * `subnet_id` - The ID of the RDMA VPC subnet.

  * `ipv4` - (Nested Attribute) IPv4 configuration for the RDMA VPC interface.

    * `addresses` - (Nested Attribute List) IPv4 addresses for the RDMA VPC interface. Each object supports:

      * `address` - The IPv4 address on the RDMA VPC interface.

      * `primary` - Whether this is the primary IPv4 address.
