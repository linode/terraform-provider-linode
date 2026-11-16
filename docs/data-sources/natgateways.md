---
page_title: "Linode: linode_natgateways"
description: |-
  Provides information about Linode NAT Gateways that match a set of filters.
---

# Data Source: linode\_natgateways

Provides information about Linode NAT Gateways that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateways). May not be currently available to all users even under v4beta.

## Example Usage

Get information about all Linode NAT Gateways with a certain label:

```hcl
data "linode_natgateways" "specific" {
  filter {
    name = "label"
    values = ["my-natgateway"]
  }
}

output "natgateway_id" {
  value = data.linode_natgateways.specific.natgateways.0.id
}
```

Get information about all Linode NAT Gateways in a certain region:

```hcl
data "linode_natgateways" "all" {
  filter {
    name = "region"
    values = ["us-ord"]
  }
}

output "natgateway_ids" {
  value = data.linode_natgateways.all.natgateways.*.id
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode NAT Gateways that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `desc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode NAT Gateway will be stored in the `natgateways` attribute and will export the following attributes:

* `natgateways` - (Nested Attribute List) The list of Linode NAT Gateways.

* `id` - The unique ID of the NAT Gateway.

* `region` - The region where the NAT Gateway is located.

* `label` - The label of the NAT Gateway.

* `default_ports_per_interface` - The default number of ports per interface for the NAT Gateway.

* [`addresses`](#addresses) - (Nested Attribute Set) The public IPv4 addresses assigned to this NAT Gateway.

* `address_autoscale_max` - The maximum number of auto-allocated IP addresses allowed for this NAT Gateway instance.

* `portset_assignments` - The number of portsets currently assigned to interfaces behind this NAT Gateway, excluding system-reserved ports.

* `portset_capacity` - The total number of portsets available to allocate across all current and potential autoscaled public IP addresses assigned to this NAT Gateway, excluding system-reserved ports.

* [`vpc_subnet`](#vpc_subnet) - Information about the VPC subnet attached to this NAT Gateway, if there is one attached.

* `created` - The date and time when the NAT Gateway was created.

* `updated` - The date and time when the NAT Gateway was last updated.

## Addresses

The following attributes are exported under each entry of the `addresses` field:

* `address` - The public IPv4 address assigned to this NAT Gateway.

## vpc_subnet

The following attributes are exported under the `vpc_subnet` field:

* `id` - The ID of the attached VPC subnet.

* `type` - The type of the attached resource. Always `vpc_subnet`.

* `url` - The API URL of the attached VPC subnet.

* `label` - The label of the VPC subnet.

* `vpc_id` - The ID of the parent VPC.

* `vpc_label` - The label of the parent VPC.

## Filterable Fields

* `id`

* `label`

* `region`
