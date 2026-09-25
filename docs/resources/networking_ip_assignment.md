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

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `region` - (Required) The region where the IP addresses will be assigned.

* `assignments` - (Required, Nested Attribute List) A list of IP/Linode assignments to apply.

## assignments

The following attributes can be defined under each entry in the `assignments` field:

* `address` - (Required) The IPv4 address or IPv6 range to assign.

* `linode_id` - (Required) The ID of the Linode to which the IP address will be assigned.

## Attribute Reference

* `id` - The unique ID of this resource.

* `assignments` - (Nested Attribute List) The list of IP/Linode assignments. In addition to the configurable arguments above, each entry exposes the following computed attributes after assignment:

  * `reserved` - Whether this IP address is a reserved IP.

  * `tags` - A set of tags associated with this IP address.

  * `assigned_entity` - (Read-Only Object) The entity this IP address has been assigned to. Referenced directly (e.g. `assigned_entity.id`).

    * `id` - The ID of the entity.

    * `label` - The label of the entity.

    * `type` - The type of the entity.

    * `url` - The URL of the entity.

## Import

Network IP assignments cannot be imported.
