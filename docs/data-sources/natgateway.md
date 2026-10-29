---
page_title: "Linode: linode_natgateway"
description: |-
  Provides details about a Linode NAT Gateway.
---

# Data Source: linode\_natgateway

Provides information about a Linode NAT Gateway.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how one might use this data source to access information about a Linode NAT Gateway.

```hcl
data "linode_natgateway" "foo" {
    id = 123
}

output "natgateway" {
  value = data.linode_natgateway.foo
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this NAT Gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

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
