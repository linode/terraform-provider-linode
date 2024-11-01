---
page_title: "Linode: linode_reserved_ip_assign"
description: |-
  Manages assigning a reserved IP address to a linode instance.
---

# linode\_reserved\_ip\_assign

Manages assigning a reserved IP address to a linode instance.
For more information, see the [API SPEC file](https://docs.google.com/document/d/1cEnLaPLdo2_I8xXywlIscsvFuP3IVve6AgzCRShbXWs/edit?tab=t.0#heading=h.2sxf35d4jytx).

~> **NOTICE:** You may need to contact support to increase your instance IP limit before you can allocate additional IPs.

~> **NOTICE:** This resource will reboot the specified instance following IP allocation.

## Example Usage

```hcl
resource "linode_instance" "web" {
  label           = "simple_instance"
  image           = "linode/ubuntu22.04"
  region          = "us-central"
  type            = "g6-standard-1"
  authorized_keys = ["ssh-rsa AAAA...Gw== user@example.local"]
  root_pass       = "this-is-not-a-safe-password"
}


resource "linode_instance_ip" "foo" {
    linode_id = linode_instance.foo.id
    public = true
    address = "desired_reservd_ip_address"
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to allocate an IPv4 address for.

* `public` - (Optional) Whether the IPv4 address is public or private. Defaults to true.

* `rdns` - (Optional) The reverse DNS assigned to this address.

* `apply_immediately` - (Optional) If true, the instance will be rebooted to update network interfaces.

* `address` - (Required) The reserved IP address that will be assigned additionally to the linode instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway` - The default gateway for this address

* `prefix` - The number of bits set in the subnet mask.

* `address` - The resulting IPv4 address.

* `region` - The region this IP resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of IP address. (`ipv4`, `ipv6`, `ipv6/pool`, `ipv6/range`)

* `id` - The resulting IPv4 address.

* `linode_id` - The ID of the linode the additional reserved IP address is assigned to.

* `public` - Whether the IPv4 address is public or private.

* `rdns` - The reverse DNS assigned to this address

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

* `apply_immediately` - If true, the instance will be rebooted to update network interfaces.

* `reserved` - The reservation status of the IP address.



