---
page_title: "Linode: linode_volume_types"
description: |-
  Provides information about Linode Volume types that match a set of filters.
---

# Data Source: linode_volume_types

Provides information about Linode Volume types that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-volume-types).

## Example Usage

Get information about all Linode Volume types with a certain label:

```hcl
data "linode_volume_types" "specific-label" {
  filter {
    name = "label"
    values = "Storage Volume"
  }
}

output "type_id" {
  value = data.linode_volume_types.specific-label.id
}
```

Get information about all Linode Volume types:

```hcl
data "linode_volume_types" "all-types" {}

output "type_id" {
  value = data.linode_volume_types.all-types.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Volume types that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Volume type will export the following attributes:

* `id` - The ID representing the Volume type.

* `label` - The Volume type label is for display purposes only.

* `price.0.hourly` -  Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `region_prices.*.id` - The Region ID for these prices.

* `region_prices.*.hourly` - Cost per hour for this region, in US dollars.

* `region_prices.*.monthly` - Cost per month for this region, in US dollars.

* `transfer` - The monthly outbound transfer amount, in MB.

## Filterable Fields

* `label`

* `transfer`
