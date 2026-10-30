il---
page_title: "Linode: linode_natgateway_addresses"
description: |-
  Provides information about the addresses assigned to a Linode NAT Gateway that match a set of filters.
---

# Data Source: linode\_natgateway\_addresses

Provides information about the addresses assigned to a Linode NAT Gateway that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway-addresses). May not be currently available to all users even under v4beta.

## Example Usage

Get information about the addresses of a NAT Gateway that are in use:

```hcl
data "linode_natgateway_addresses" "in_use" {
  natgateway_id = 123

  filter {
    name = "in_use"
    values = ["true"]
  }
}

output "address_values" {
  value = data.linode_natgateway_addresses.in_use.addresses.*.address
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `natgateway_id` - (Required) The ID of the NAT Gateway whose addresses are listed.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select addresses that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `desc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each address will be stored in the `addresses` attribute and will export the following attributes:

* `addresses` - (Nested Attribute List) The list of addresses assigned to the NAT Gateway.

* `address` - The public IPv4 address assigned to the NAT Gateway.

* `in_use` - Whether the address is currently in use by one or more interfaces.

* `interface_count` - The number of interfaces using this address.

* `interface_url` - The API URL of the interfaces using this address.

* `portset_assignments` - The number of portsets currently assigned to interfaces using this address, excluding system-reserved ports.

* `portset_capacity` - The total number of portsets available to allocate for this address, excluding system-reserved ports.

## Filterable Fields

* `address`

* `in_use`
