---
layout: "linode"
page_title: "Linode: linode_ipv6_range"
sidebar_current: "docs-linode-resource-ipv6-range"
description: |-
  Manages a Linode IPv6 range.
---

# linode\_ipv6\_range

Manages a Linode IPv6 range.

When an IPv6 range initially assigned to a linode by Terraform, then user removes it outside of the Terraform and reassigns it to another linode, there is a chance that same range can be assigned to another linode, although the range is randomly selected. In this case, user should be manually tainting it.

## Example Usage

```terraform
resource "linode_instance" "foobar" {
  label = "my-linode"
  image = "linode/alpine3.14"
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
