---
page_title: "Linode: linode_natgateway_interfaces"
description: |-
  Provides information about the interfaces attached to a Linode NAT Gateway.
---

# Data Source: linode\_natgateway\_interfaces

Provides information about the interfaces attached to a Linode NAT Gateway. When `address` is set, only the interfaces using that public IPv4 address are returned.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway-interfaces). May not be currently available to all users even under v4beta.

## Example Usage

Get all interfaces attached to a NAT Gateway:

```hcl
data "linode_natgateway_interfaces" "all" {
  natgateway_id = linode_natgateway.example.id
}

output "interface_ids" {
  value = data.linode_natgateway_interfaces.all.interfaces.*.id
}
```

Get the interfaces using a specific public address of a NAT Gateway:

```hcl
data "linode_natgateway_interfaces" "by_address" {
  natgateway_id = linode_natgateway.example.id
  address       = "192.0.2.10"
}

output "linode_labels" {
  value = data.linode_natgateway_interfaces.by_address.interfaces.*.linode.label
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `natgateway_id` - (Required) The ID of the NAT Gateway to list interfaces for.

* `address` - (Optional) A public IPv4 address of the NAT Gateway. When set, only interfaces using this address are returned.

## Attributes Reference

Each NAT Gateway interface will be stored in the `interfaces` attribute and will export the following attributes:

* `id` - The data source's unique ID. This is the NAT Gateway ID, or the NAT Gateway ID and address when `address` is set.

* [`interfaces`](#interfaces) - (Nested Attribute List) The list of NAT Gateway interfaces.

## interfaces

The following attributes are exported under each entry of the `interfaces` field:

* `id` - The unique ID of the NAT Gateway interface.

* [`linode`](#linode) - (Nested Attribute) The Linode the interface belongs to.

* `addresses` - (List of Strings) The public IPv4 addresses assigned to the interface.

* [`portsets`](#portsets) - (Nested Attribute List) The portsets assigned to the interface. May not be available for all accounts.

## linode

The following attributes are exported under the `linode` field:

* `id` - The ID of the Linode.

* `label` - The label of the Linode.

* `type` - The type of the attached resource. Always `linode`.

* `url` - The API URL of the Linode.

## portsets

The following attributes are exported under each entry of the `portsets` field:

* `address` - The public IPv4 address the portset belongs to.

* [`ports`](#ports) - (Nested Attribute List) The port ranges in this portset.

## ports

The following attributes are exported under each entry of the `ports` field:

* `start` - The first port in the range.

* `end` - The last port in the range.

