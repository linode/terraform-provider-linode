---
layout: "linode"
page_title: "Linode: linode_instance_ip"
sidebar_current: "docs-linode-instance-ip"
description: |-
  Manages a Linode instance IP.
---

# linode\_instance\_ip

~> **NOTICE:** You may need to contact support to increase your instance IP limit before you can allocate additional IPs.

~> **NOTICE:** This resource will reboot the specified instance following IP allocation.

Manages a Linode instance IP.

## Example Usage

```terraform
resource "linode_instance" "foo" {
    image = "linode/alpine3.16"
    label = "foobar-test"
    type = "g6-nanode-1"
    region = "us-east"
}

resource "linode_instance_ip" "foo" {
    linode_id = linode_instance.foo.id
    public = true
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to allocate an IPv4 address for.

* `public` - (Optional) Whether the IPv4 address is public or private. Defaults to true.

* `rdns` - (Optional) The reverse DNS assigned to this address.

* `apply_immediately` - (Optional) If true, the instance will be rebooted to update network interfaces.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway` - The default gateway for this address

* `prefix` - The number of bits set in the subnet mask.

* `address` - The resulting IPv4 address.

* `region` - The region this IP resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `type` - The type of IP address. (`ipv4`, `ipv6`, `ipv6/pool`, `ipv6/range`)
