---
page_title: "Linode: linode_networking_ip"
description: |-
  Provides details about a Linode Networking IP Address.
---

# Data Source: linode\_network\_ip

Provides information about a Linode Networking IP Address
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-ip).

## Example Usage

The following example shows how one might use this data source to access information about a Linode Networking IP Address.

```hcl
data "linode_networking_ip" "ns1_linode_com" {
    address = "162.159.27.72"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The IP Address to access.  The address must be associated with the account and a resource that the user has access to view.

## Attributes Reference

The Linode Network IP Address resource exports the following attributes:

* `address` - The IP address.

* `gateway` - The default gateway for this address.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `prefix` - The number of bits set in the subnet mask.

* `type` - The type of address this is (ipv4, ipv6, ipv6/pool, ipv6/range).

* `public` - Whether this is a public or private IP address.

* `rdns` - The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `linode_id` - The ID of the Linode this address currently belongs to.

* `region` - The Region this IP address resides in. See all regions [here](https://api.linode.com/v4/regions).

* `reserved` - Whether this IP address is a reserved IP.

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.

  * `subnet_id` - The `id` of the VPC Subnet for this Interface.

  * `vpc_id` - The `id` of the VPC configured for this Interface.
