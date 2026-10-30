---
page_title: "Linode: linode_natgateway_address"
description: |-
  Provides details about a Linode NAT Gateway address.
---

# Data Source: linode\_natgateway\_address

Provides information about an address assigned to a Linode NAT Gateway.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway-address). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how one might use this data source to access information about an address assigned to a Linode NAT Gateway.

```hcl
data "linode_natgateway_address" "foo" {
    natgateway_id = 123
    address = "1.2.3.4"
}

output "natgateway_address" {
  value = data.linode_natgateway_address.foo
}
```

## Argument Reference

The following arguments are supported:

* `natgateway_id` - (Required) The ID of the NAT Gateway the address is assigned to.

* `address` - (Required) The public IPv4 address assigned to the NAT Gateway. This is the identifier of the NAT Gateway address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `in_use` - Whether the address is currently in use by one or more interfaces.

* `interface_count` - The number of interfaces using this address.

* `interface_url` - The API URL of the interfaces using this address.

* `portset_assignments` - The number of portsets currently assigned to interfaces using this address, excluding system-reserved ports.

* `portset_capacity` - The total number of portsets available to allocate for this address, excluding system-reserved ports.

