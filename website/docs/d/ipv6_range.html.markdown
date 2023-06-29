---
layout: "linode"
page_title: "Linode: linode_ipv6_range"
sidebar_current: "docs-linode-datasource-ipv6-range"
description: |-
  Provides details about a Linode IPv6 Range.
---

# Data Source: linode\_ipv6\_range

Provides information about a Linode IPv6 Range.

## Example Usage

Get information about an IPv6 range assigned to a Linode:

```hcl
data "linode_ipv6_range" "range-info" {
  range = "2001:0db8::"
}
```

## Argument Reference

The following arguments are supported:

* `range` - (Required) The IPv6 range to retrieve information about.

## Attributes Reference

The `linode_ipv6_range` data source exports the following attributes:

* `ip_bgp` - Whether this IPv6 range is shared.

* `linodes` - A set of Linodes targeted by this IPv6 range. Includes Linodes with IP sharing.

* `prefix` - The prefix length of the address, denoting how many addresses can be assigned from this range.

* `region` - The region for this range of IPv6 addresses.
