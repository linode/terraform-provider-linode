---
layout: "linode"
page_title: "Linode: linode_instance_networking"
sidebar_current: "docs-linode-datasource-instance-networking"
description: |-
  Provides details about the networking configuration of an Instance.
---

# Data Source: linode\_instance\_networking

Provides details about the networking configuration of an Instance.

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

* `linode_id` - The ID of the Linode this address currently belongs to. For IPv4 addresses, this is by default the Linode that this address was assigned to on creation, and these addresses my be moved using the [/networking/ipv4/assign](https://www.linode.com/docs/api/networking/#ips-to-linodes-assign) endpoint. For SLAAC and link-local addresses, this value may not be changed.

* `prefix` - The number of bits set in the subnet mask.

* `public` - Whether this is a public or private IP address.

* `rdns` - (Nullable) The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `region` - (Filterable) The Region this IP address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of address this is.

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
