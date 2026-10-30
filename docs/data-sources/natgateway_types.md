---
page_title: "Linode: linode_natgateway_types"
description: |-
  Provides information about Linode NAT Gateway types that match a set of filters.
---

# Data Source: linode\_natgateway\_types

Provides information about Linode NAT Gateway types that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway-types). May not be currently available to all users even under v4beta.

## Example Usage

Get information about the Linode NAT Gateway type with a certain ID:

```hcl
data "linode_natgateway_types" "specific" {
  filter {
    name = "id"
    values = ["g1-natgateway"]
  }
}

output "type_label" {
  value = data.linode_natgateway_types.specific.types.0.label
}
```

Get information about all Linode NAT Gateway types with a certain label:

```hcl
data "linode_natgateway_types" "labeled" {
  filter {
    name = "label"
    values = ["NAT Gateway"]
  }
}

output "type_ids" {
  value = data.linode_natgateway_types.labeled.types.*.id
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode NAT Gateway types that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `desc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode NAT Gateway type will export the following attributes:

* `types` - (Nested Attribute List) The returned list of NAT Gateway types. Referenced by index (e.g. `types[0].id`).

* `id` - The ID representing the NAT Gateway type.

* `label` - The NAT Gateway type label is for display purposes only.

* `price` - (Read-Only Object List) Pricing information for this NAT Gateway type. Referenced with an index (e.g. `price.0.hourly`).

* `price.0.hourly` - Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

## Filterable Fields

* `id`

* `label`
