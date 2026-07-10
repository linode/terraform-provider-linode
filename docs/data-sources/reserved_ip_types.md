---
page_title: "Linode: linode_reserved_ip_types"
description: |-
  Provides information about Linode Reserved IP types that match a set of filters.
---

# Data Source: linode\_reserved\_ip\_types

Provides information about Linode Reserved IP types that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-reserved-ip-types).

## Example Usage

Get information about all Reserved IP types:

```hcl
data "linode_reserved_ip_types" "all" {}

output "type_ids" {
  value = data.linode_reserved_ip_types.all.types.*.id
}
```

Get information about a specific Reserved IP type by ID:

```hcl
data "linode_reserved_ip_types" "specific" {
  filter {
    name   = "id"
    values = ["reserved-ip"]
  }
}

output "label" {
  value = data.linode_reserved_ip_types.specific.types[0].label
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Reserved IP types that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Reserved IP type will be stored in the `types` attribute and will export the following attributes:

* `id` - The unique ID of this Reserved IP type (e.g. `reserved-ip`).

* `label` - The human-readable label for this Reserved IP type.

* `price.0.hourly` - Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `region_prices.*.id` - The ID of the region this price entry applies to.

* `region_prices.*.hourly` - The hourly price for this Reserved IP type in the given region.

* `region_prices.*.monthly` - The monthly price for this Reserved IP type in the given region.

## Filterable Fields

* `label`
