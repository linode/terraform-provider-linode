---
page_title: "Linode: linode_networking_ip"
description: |-
  Manages the allocation and assignment of an IP addresses.
---

# linode\_networking\_ip

Manages allocation of reserved IPv4 address in a region and optionally assigning the reserved address to a Linode instance.

For more information, see the corresponding [API documentation](https://techdocs.akamai.com/linode-api/reference/post-allocate-ip).

## Example Usage

```hcl
resource "linode_networking_ip" "test_ip" {
  type     = "ipv4"
  linode_id = 12345
  public   = true
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of IP address. (ipv4, ipv6, etc.)

* `public` - (Optional) Whether the IP address is public. Defaults to true.

* `linode_id` - (Optional) The ID of the Linode to which the IP address will be assigned. Updating this field on an ephemeral IP will trigger a recreation. Conflicts with `region`.

### Reserved-Specific Argument Reference

-> **Note:** IP reservation is not currently accessible to all users.

The following arguments are only available when reserving an IP address:

* `reserved` - (Optional) Whether this IP address should be a reserved IP.

* `region` - (Optional) The region where the reserved IP should be allocated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of this IP address.

* `address` - The IP address.

* `gateway` - The default gateway for this address.

* `prefix` - The number of bits set in the subnet mask.

* `rdns` - The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.

  * `subnet_id` - The `id` of the VPC Subnet for this Interface.

  * `vpc_id` - The `id` of the VPC configured for this Interface.

## Import

IP addresses can be imported using the IP address ID, e.g.

```sh
terraform import linode_networking_ip.example_ip 172.104.30.209
```
