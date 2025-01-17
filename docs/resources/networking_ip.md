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

* `linode_id` - (Optional) The ID of the Linode to which the IP address will be assigned. Conflicts with `region`.

### Reserved-Specific Argument Reference

-> **Note:** IP reservation is not currently accessible to all users.

The following arguments are only available when reserving an IP address:

* `reserved` - (Optional) Whether this IP address should be a reserved IP.

* `region` - (Optional) The region where the reserved IP should be allocated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address` - The IP address.

* `id` - The IP address.

## Import

IP addresses can be imported using the IP address ID, e.g.

```sh
terraform import linode_networking_ip.example_ip 172.104.30.209
