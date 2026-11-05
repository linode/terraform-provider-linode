---
page_title: "Linode: linode_natgateway"
description: |-
  Manages a Linode NAT Gateway.
---

# linode\_natgateway

Manages a Linode NAT Gateway.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-natgateways). May not be currently available to all users even under v4beta.

## Example Usage

Create a NAT Gateway that uses autoscaling:

```terraform
resource "linode_natgateway" "autoscaling" {
    label = "test-natgateway"
    region = "us-iad"
    use_autoscaling = true
}
```

Create a NAT Gateway with manually managed addresses:

```terraform
resource "linode_networking_ip" "test" {
    type = "ipv4"
    region = "us-iad"
    public = true
    reserved = true
}

resource "linode_natgateway" "manual" {
    label = "test-natgateway"
    region = "us-iad"
    use_autoscaling = false

    addresses = [
        { address = linode_networking_ip.test.address },
    ]
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `label` - (Required) The label for the NAT Gateway.

* `region` - (Required) The region where the NAT Gateway will be created. Changing this value forces the creation of a new NAT Gateway.

* `default_ports_per_interface` - (Optional) The default number of ports per interface for the NAT Gateway. If omitted, the API assigns a value. Changing this value forces the creation of a new NAT Gateway.

* `use_autoscaling` - (Optional) Whether the NAT Gateway automatically scales its public IPv4 addresses. If omitted, defaults to `true` when `addresses` is empty and `false` otherwise. This value is write-only and is not returned by the API. Changing this value forces the creation of a new NAT Gateway.

* [`addresses`](#addresses) - (Optional, Nested Attribute Set) The public IPv4 addresses assigned to this NAT Gateway. Must be empty when `use_autoscaling` is `true`. If omitted, the API assigns addresses.

## Addresses

Configures a single public IPv4 address assigned to this NAT Gateway.

* `address` - (Required) The public IPv4 address to assign to this NAT Gateway.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the NAT Gateway.

* `address_autoscale_max` - The maximum number of auto-allocated IP addresses allowed for this NAT Gateway instance.

* `portset_assignments` - The number of portsets currently assigned to interfaces behind this NAT Gateway, excluding system-reserved ports.

* `portset_capacity` - The total number of portsets available to allocate across all current and potential autoscaled public IP addresses assigned to this NAT Gateway, excluding system-reserved ports.

* [`vpc_subnet`](#vpc_subnet) - Information about the VPC subnet attached to this NAT Gateway, if there is one attached.

* `created` - The date and time when the NAT Gateway was created.

* `updated` - The date and time when the NAT Gateway was last updated.

## vpc_subnet

The following attributes are exported under the `vpc_subnet` field:

* `id` - The ID of the attached VPC subnet.

* `type` - The type of the attached resource. Always `vpc_subnet`.

* `url` - The API URL of the attached VPC subnet.

* `label` - The label of the VPC subnet.

* `vpc_id` - The ID of the parent VPC.

* `vpc_label` - The label of the parent VPC.
