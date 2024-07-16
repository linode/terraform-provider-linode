---
page_title: "Linode: linode_instance_networking"
description: |-
  Provides details about the networking configuration of an Instance.
---

# Data Source: linode\_instance\_networking

Provides details about the networking configuration of an Instance.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-linode-config-interfaces).

## Example Usage

```terraform
data "linode_instance_networking" "example" {
    linode_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The Linode instance's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* [`ipv4`](#ipv4) - Information about this Linode’s IPv4 addresses.

* [`ipv6`](#ipv6) - Information about this Linode’s IPv6 addresses.

### IPv4

The following attributes are available for each Linode’s IPv4 addresses:

* [`private`] (#private) - A list of private IP Address objects belonging to this Linode.

* [`public`] (#public) - A list of public IP Address objects belonging to this Linode.

* [`reserved`] (#reserved) - A list of reserved IP Address objects belonging to this Linode.

* [`shared`] (#shared)- A list of shared IP Address objects assigned to this Linode.

* [`vpc`] (#vpc)- A list of VPC IP Address objects assigned to this Linode.

### IPv6

The following attributes are available for each Linode’s IPv6 addresses:

* [`global`] (#global) - An object representing an IPv6 pool.

* [`link_local`] (#link_local) - A link-local IPv6 address object that exists in Linode’s system.

* [`slaac`] (#slaac) - A SLAAC IPv6 address object that exists in Linode’s system.

### Private

A list of private IP Address objects belonging to this Linode.

* `address` - The private IPv4 address.

* `gateway` - The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to.

* `prefix` - The number of bits set in the subnet mask.

* `public` - Whether this is a public or private IP address.

* `rdns` - The reverse DNS assigned to this address.

* `region` - (Filterable) The Region this address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of address this is.

### Public

A list of public IP Address objects belonging to this Linode.

* `address` - The IP address.

* `gateway` - (Nullable) The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to. For IPv4 addresses, this is by default the Linode that this address was assigned to on creation, and these addresses my be moved using the [/networking/ipv4/assign](https://techdocs.akamai.com/linode-api/reference/post-assign-ips) endpoint. For SLAAC and link-local addresses, this value may not be changed.

* `prefix` - The number of bits set in the subnet mask.

* `public` - Whether this is a public or private IP address.

* `rdns` - (Nullable) The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `region` - (Filterable) The Region this IP address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of address this is.

* `vpc_nat_1_1` - IPv4 address configured as a 1:1 NAT for this Interface.
  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.
  * `subnet_id` - The `id` of the VPC Subnet for this Interface.
  * `vpc_id` - The `id` of the VPC configured for this Interface.

### Reserved

A list of reserved IP Address objects belonging to this Linode.

* `address` - The IP address.

* `gateway` - (Nullable) The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to. For IPv4 addresses, this is by default the Linode that this address was assigned to on creation, and these addresses my be moved using the [/networking/ipv4/assign](https://www.linode.com/docs/api/networking/#ips-to-linodes-assign) endpoint. For SLAAC and link-local addresses, this value may not be changed.

* `prefix` - The number of bits set in the subnet mask.

* `public` - Whether this is a public or private IP address.

* `rdns` - (Nullable) The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `region` - (Filterable) The Region this IP address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of address this is.

* `vpc_nat_1_1` - IPv4 address configured as a 1:1 NAT for this Interface.
  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.
  * `subnet_id` - The `id` of the VPC Subnet for this Interface.
  * `vpc_id` - The `id` of the VPC configured for this Interface.

### Shared

A list of shared IP Address objects assigned to this Linode.

* `address` - The IP address.

* `gateway` - (Nullable) The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to. For IPv4 addresses, this is by default the Linode that this address was assigned to on creation, and these addresses my be moved using the [/networking/ipv4/assign](https://www.linode.com/docs/api/networking/#ips-to-linodes-assign) endpoint. For SLAAC and link-local addresses, this value may not be changed.

* `prefix` - The number of bits set in the subnet mask.

* `public` - Whether this is a public or private IP address.

* `rdns` - (Nullable) The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `region` - (Filterable) The Region this IP address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of address this is.

* `vpc_nat_1_1` - IPv4 address configured as a 1:1 NAT for this Interface.
  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.
  * `subnet_id` - The `id` of the VPC Subnet for this Interface.
  * `vpc_id` - The `id` of the VPC configured for this Interface.

### VPC

A list of VPC IP Address objects assigned to this Linode.

* `address` - An IPv4 address configured for this VPC interface. It will be `null` if it's an `address_range`.

* `address_range` - A range of IPv4 addresses configured for this VPC interface. it will be `null` if it's a single `address`.

* `active` - Returns `true` if the VPC interface is in use, meaning that the Linode was powered on using the `config_id` to which the interface belongs. Otherwise returns `false`.

* `vpc_id` - The unique globally general API entity identifier for the VPC.

* `subnet_id` - The unique globally general API entity identifier for the VPC subnet.

* `config_id` - The globally general entity identifier for the Linode configuration profile where the VPC is included.

* `interface_id` - The globally general API entity identifier for the Linode interface.

* `gateway` - The default gateway for the VPC subnet that the IP or IP range belongs to.

* `prefix` - The number of bits set in the `subnet_mask`.

* `region` - The region of the VPC.

* `subnet_mask` - The mask that separates host bits from network bits for the `address` or `address_range`.

* `linode_id` - The identifier for the Linode the VPC interface currently belongs to.

* `nat_1_1` - The public IP address used for NAT 1:1 with the VPC. This is `null` if the VPC interface uses an `address_range` or NAT 1:1 isn't used.

### Global

An object representing an IPv6 pool.

* `prefix` - The prefix length of the address, denoting how many addresses can be assigned from this pool calculated as 2^128-prefix.

* `range` - The IPv6 range of addresses in this pool.

* `region` - (Filterable) The region for this pool of IPv6 addresses.

* `route_target` - (Nullable) The last address in this block of IPv6 addresses.

### Link_Local

A link-local IPv6 address that exists in Linode’s system.

* `address` - The IPv6 link-local address.

* `gateway` - The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to.

* `prefix` - The network prefix.

* `public` - Whether this is a public or private IP address.

* `rdns` - The reverse DNS assigned to this address.

* `region` - (Filterable) The Region this address resides in.

* `subnet_mask` - The subnet mask.

* `type` - The type of address this is.

### SLAAC

A SLAAC IPv6 address object that exists in Linode’s system.

* `address` - The address.

* `gateway` - The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to.

* `prefix` - The network prefix.

* `public` - Whether this is a public or private IP address.

* `rdns` - The reverse DNS assigned to this address.

* `region` - (Filterable) The Region this address resides in.

* `subnet_mask` - The subnet mask.

* `type` - The type of address this is.
