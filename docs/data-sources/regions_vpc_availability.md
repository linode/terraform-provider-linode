---
page_title: "Linode: linode_regions_vpc_availability"
description: |-
  Provides details about all regions VPC availability.
---

# linode\_regions\_vpc\_availability

Provides information about Linode regions VPC availability.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-regions-vpc-availability).

```hcl
data "linode_regions_vpc_availability" "all-regions" {}

output "regions_vpc_availability" {
  value = data.linode_regions_vpc_availability.all-regions.regions_vpc_availability
}
```

## Attributes Reference

Each Linode region VPC availability will be stored in the `regions_vpc_availability` attribute and will export the following attributes:

* `available` - Whether VPCs can be created in the region.

* `available_ipv6_prefix_lengths` - The IPv6 prefix lengths allowed when creating a VPC in the region.

* `id` - The Region ID.
