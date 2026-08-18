---
page_title: "Linode: linode_reserved_ip"
description: |-
  Manages a Linode reserved IP address.
---

# linode\_reserved\_ip

Manages a reserved IPv4 address in a Linode region. Reserved IPs persist independently of any Linode instance and can be assigned or unassigned at any time.

For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-reserved-ip).

~> **NOTICE:** IP reservation is not currently accessible to all users. Contact support to request access.

## Example Usage

### Basic Reserved IP

```terraform
resource "linode_reserved_ip" "example" {
  region = "us-east"
}
```

### Reserved IP with Tags

```terraform
resource "linode_reserved_ip" "example" {
  region = "us-east"
  tags   = ["production", "web"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region in which to reserve the IP address. Changing this forces a new resource to be created.

* `tags` - (Optional) A set of tags to assign to this reserved IP address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of this reserved IP address (same as `address`).

* `address` - The reserved IPv4 address.

* `gateway` - The default gateway for this address.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `prefix` - The number of bits set in the subnet mask.

* `type` - The type of IP address.

* `public` - Whether this is a public IP address.

* `rdns` - The reverse DNS assigned to this address.

* `linode_id` - The ID of the Linode this address is currently assigned to, if any.

* `reserved` - Whether this IP address is reserved. Always `true` for this resource.

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

* `assigned_entity` - The entity this reserved IP address is currently assigned to, if any. Contains:

  * `id` - The ID of the assigned entity.

  * `label` - The label of the assigned entity.

  * `type` - The type of the assigned entity (e.g. `linode`).

  * `url` - The API URL of the assigned entity.

## Import

Reserved IP addresses can be imported using the IP `address`:

```sh
terraform import linode_reserved_ip.example 203.0.113.1
```
