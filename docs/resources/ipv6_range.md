---
page_title: "Linode: linode_ipv6_range"
description: |-
  Manages a Linode IPv6 range.
---

# linode\_ipv6\_range

Manages a Linode IPv6 range.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-ipv6-range).

~> **NOTICE:** We highly recommend that users do not remove an IPv6 range created by this Terraform resource outside of Terraform. This is because if a user manually removes an IPv6 range created by Terraform, and then assigns some IPv6 ranges to other linodes outside of Terraform, there is a chance that the same IPv6 range can be assigned to another linode, even though the new range is randomly selected. This will result in the newly assigned IPv6 range being managed by this Terraform resource. In this case, the user should manually taint this resource.

## Example Usage

```terraform
resource "linode_instance" "foobar" {
  label = "my-linode"
  image = "linode/alpine3.19"
  type = "g6-nanode-1"
  region = "us-southeast"
}

resource "linode_ipv6_range" "foobar" {
  linode_id = linode_instance.foobar.id

  prefix_length = 64
}
```

## Argument Reference

The following arguments are supported:

**Note:** Only one of`linode_id` and `route_target` should be specified.

* `prefix_length` - (Required) The prefix length of the IPv6 range.

* `linode_id` - (Required) The ID of the Linode to assign this range to. This field may be updated to reassign the IPv6 range.

* `route_target` - (Required) The IPv6 SLAAC address to assign this range to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `is_bgp` - Whether this IPv6 range is shared.

* `linodes` - A list of Linodes targeted by this IPv6 range. Includes Linodes with IP sharing.

* `range` - The IPv6 range of addresses in this pool.

* `region` - The region for this range of IPv6 addresses.
