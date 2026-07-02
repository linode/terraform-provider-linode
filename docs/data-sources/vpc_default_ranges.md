---
page_title: "Linode: linode_vpc_default_ranges"
description: |-
  Provides details about the default and forbidden IPv4 address ranges for VPCs.
---

# Data Source: linode\_vpc\_default\_ranges

Provides information about the default and forbidden IPv4 address ranges for VPCs.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpcs-default-ranges).

## Example Usage

The following example shows how one might use this data source to access information about the default and forbidden IPv4 address ranges for VPCs.

```hcl
data "linode_vpc_default_ranges" "foo" {}

output "vpc_default_ranges" {
  value = data.linode_vpc_default_ranges.foo
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_ipv4_ranges` - The default IPv4 address ranges (CIDR blocks) used by the system when creating a VPC if no custom ranges are provided.

* `forbidden_ipv4_ranges` - IPv4 address ranges (CIDR blocks) that are forbidden and can't be used as part of a VPC's address space.
