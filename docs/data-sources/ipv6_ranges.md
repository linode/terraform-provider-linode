---
page_title: "Linode: linode_ipv6_ranges"
description: |-
  Lists IPv6 ranges on your Account.
---

# linode\_ipv6\_ranges

Provides information about Linode IPv6 ranges that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-ipv6-ranges).

~> Some fields may not be accessible directly the results of this data source.
For additional information about a specific IPv6 range consider using the [linode_ipv6_range](ipv6_range.md)
data source.

```hcl
data "linode_ipv6_ranges" "filtered-ranges" {
  filter {
    name = "region"
    values = ["us-mia"]
  }
}

output "ranges" {
  value = data.linode_ipv6_ranges.filtered-ranges
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode IPv6 ranges that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode IPv6 range will be stored in the `ranges` attribute and will export the following attributes:

* `range` - The IPv6 address of this range.

* `route_target` - The IPv6 SLAAC address.

* `prefix` - The prefix length of the address, denoting how many addresses can be assigned from this range.

* `region` - The region for this range of IPv6 addresses.

## Filterable Fields

* `range`

* `route_target`

* `prefix`

* `region`
