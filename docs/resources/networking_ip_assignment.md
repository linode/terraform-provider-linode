---
page_title: "Linode: linode_networking_ip_assignment"
description: |-
  Managed the assignment multiple IPv4 addresses and/or IPv6 ranges to multiple Linodes in a Region.
---

# linode_networking_ip_assignment

Manages the assignment of multiple IPv4 addresses and/or IPv6 ranges to multiple Linodes in a specified region.


For more information, see the corresponding [API documentation](https://techdocs.akamai.com/linode-api/reference/post-assign-ips).

## Example Usage

```hcl
resource "linode_networking_ip_assignment" "foobar" {
  region = "us-mia"
  
  assignments = [
    {
      address   = linode_networking_ip.reserved_ip1.address
      linode_id = linode_instance.terraform-web1.id
    },
    {
      address   = linode_networking_ip.reserved_ip2.address
      linode_id = linode_instance.terraform-web2.id
    },
  ]
}
```

## Argument Reference

* `region` - (Required) The region where the IP addresses will be assigned.

* `assignments` - (Required) A list of IP/Linode assignments to apply.

## assignments

The following attributes can be defined under each entry in the `assignments` field:

* `address` - (Required) The IPv4 address or IPv6 range to assign.

* `linode_id` - (Required) The ID of the Linode to which the IP address will be assigned.

## Attribute Reference

* `id` - The unique ID of this resource.

## Import

Network IP assignments cannot be imported.
