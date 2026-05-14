---
page_title: "Linode: linode_reserved_ip_assignment"
description: |-
  Manages assignment of a reserved IP address to a Linode instance.
---

# linode\_reserved\_ip\_assignment

Manages the assignment of a reserved IPv4 address to a Linode instance.

For more information, see the corresponding [API documentation](https://techdocs.akamai.com/linode-api/reference/post-add-linode-ip).

## Example Usage

```hcl
resource "linode_reserved_ip_assignment" "example" {
  linode_id = linode_instance.example.id
  address   = linode_networking_ip.reserved.address
  public    = true
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to assign the reserved IP to. Changing this forces creation of a new resource.

* `address` - (Required) The reserved IPv4 address to assign to the Linode.

* `public` - (Optional) Whether the IP address is public. Defaults to `true`. This must match the reserved IP's existing public/private status. Changing this forces creation of a new resource.

* `rdns` - (Optional) The reverse DNS assigned to this address. Configured via a separate API call after the IP is assigned.

* `apply_immediately` - (Optional) If true, the instance will be rebooted to update network interfaces. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IPv4 address (the address itself).

* `gateway` - The default gateway for this address.

* `prefix` - The number of bits set in the subnet mask.

* `region` - The region this IP resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of IP address.

* `reserved` - The reservation status of the IP address.

* `tags` - A set of tags associated with this IP address.

* `assigned_entity` - The entity this IP address has been assigned to. This is null if the address is not assigned to an entity.

  * `id` - The ID of the entity.

  * `label` - The label of the entity.

  * `type` - The type of the entity.

  * `url` - The URL of the entity.

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.

  * `subnet_id` - The `id` of the VPC Subnet for this Interface.

  * `vpc_id` - The `id` of the VPC configured for this Interface.
