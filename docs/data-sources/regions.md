---
page_title: "Linode: linode_regions"
description: |-
  Lists the Regions available for Linode services. Not all services are guaranteed to be available in all Regions.
---

# linode\_regions

Provides information about Linode regions that match a set of filters.

```hcl
data "linode_regions" "filtered-regions" {
  filter {
    name = "status"
    values = ["ok"]
  }

  filter {
    name = "capabilities"
    values = ["NodeBalancers"]
  }
}

output "regions" {
  value = data.linode_regions.filtered-regions.regions
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode regions that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode region will be stored in the `regions` attribute and will export the following attributes:

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

## Filterable Fields

* `status`

* `country`

* `capabilities`

* `site_type`
