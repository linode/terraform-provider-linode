---
page_title: "Linode: linode_region"
description: |-
  Provides details about a specific service region
---

# Data Source: linode\_region

`linode_region` provides details about a specific Linode region. See all regions [here](https://api.linode.com/v4/regions).
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-region).

## Example Usage

The following example shows how the resource might be used to obtain additional information about a Linode region.

```hcl
data "linode_region" "region" {
  id = "us-east"
}
```

## Argument Reference

* `id` - (Required) The code name of the region to select.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `country` - The country the region resides in.

* `label` - Detailed location information for this Region, including city, state or region, and country.

* `capabilities` - A list of capabilities of this region.

* `status` - This region’s current operational status (ok or outage).

* `site_type` - The type of this region.

* [`resolvers`] (#resolvers) - An object representing the IP addresses for this region's DNS resolvers.

* [`placement_group_limits`] (#placement-group-limits) - An object representing the limits relating to placement groups in this region.

### Resolvers

* `ipv4` - The IPv4 addresses for this region’s DNS resolvers, separated by commas.

* `ipv6` - The IPv6 addresses for this region’s DNS resolvers, separated by commas.

### Placement Group Limits

* `maximum_pgs_per_customer` - The maximum number of placement groups allowed for the current user in this region.

* `maximum_linodes_per_pg` - The maximum number of Linodes allowed to be assigned to a placement group in this region.
