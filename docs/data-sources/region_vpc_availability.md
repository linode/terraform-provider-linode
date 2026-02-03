---
page_title: "Linode: linode_region_vpc_availability"
description: |-
  Provides details about a specific region VPC availability.
---

# Data Source: linode\_region\_vpc\_availability

`linode_region_vpc_availability` provides details about a specific Linode region VPC availability. See all regions VPC availability [here](https://api.linode.com/v4/regions/vpc-availability).
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-region-vpc-availability).

## Example Usage

The following example shows how the resource might be used to obtain additional information about a Linode region VPC availability.

```hcl
data "linode_region_vpc_availability" "region_vpc_availability" {
  id = "us-east"
}
```

## Argument Reference

* `id` - (Required) The code name of the region to select.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available` - Whether VPCs can be created in the region.

* `available_ipv6_prefix_lengths` - The IPv6 prefix lengths allowed when creating a VPC in the region.