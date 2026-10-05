---
page_title: "Linode: linode_network_transfer_prices"
description: |-
  Provides information about Linode Network Transfer Prices that match a set of filters.
---

# Data Source: linode_network_transfer_prices

Provides information about Linode Network Transfer Prices that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-network-transfer-prices).

## Example Usage

Get information about all Linode Network Transfer Prices with a certain label:

```hcl
data "linode_network_transfer_prices" "specific-label" {
  filter {
    name = "label"
    values = "Network Transfer"
  }
}

output "price_id" {
  value = data.linode_network_transfer_prices.specific-label.id
}
```

Get information about all Linode Network Transfer Prices:

```hcl
data "linode_network_transfer_prices" "all-prices" {}

output "price_id" {
  value = data.linode_network_transfer_prices.all-prices.*.id
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode Network Transfer Prices that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Network Transfer Price will export the following attributes:

* `types` - (Nested Attribute List) The Network Transfer Prices returned by this data source.

* `id` - The ID representing the Network Transfer Price.

* `label` - The Network Transfer Price label is for display purposes only.

* `price` - (Read-Only Object List) The base cost for this Network Transfer Price. Referenced with an index (e.g. `price.0.hourly`).

* `price.0.hourly` -  Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `region_prices` - (Read-Only Object List) Region-specific pricing for this Network Transfer Price. Referenced with an index (e.g. `region_prices.0.id`).

* `region_prices.*.id` - The Region ID for these prices.

* `region_prices.*.hourly` - Cost per hour for this region, in US dollars.

* `region_prices.*.monthly` - Cost per month for this region, in US dollars.

* `transfer` - The monthly outbound transfer amount, in MB.

## Filterable Fields

* `label`

* `transfer`
