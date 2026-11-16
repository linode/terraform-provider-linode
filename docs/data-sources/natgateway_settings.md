---
page_title: "Linode: linode_natgateway_settings"
description: |-
  Provides information about Linode NAT Gateway settings for the current account.
---

# Data Source: linode\_natgateway\_settings

Provides information about Linode NAT Gateway settings and limits for the current account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-natgateway-settings). May not be currently available to all users even under v4beta.

## Example Usage

Get the NAT Gateway settings:

```hcl
data "linode_natgateway_settings" "natgw" {}

output "allowed_ports_per_interface" {
  value = data.linode_natgateway_settings.natgw.allowed_ports_per_interface
}

output "maximum_reserved_addresses_per_natgateway" {
  value = data.linode_natgateway_settings.natgw.maximum_reserved_addresses_per_natgateway
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

The following attributes are exported:

* `allowed_ports_per_interface` - (List of Integers) The port counts allowed per NAT Gateway interface.

* `maximum_autoscaling_addresses_per_natgateway` - (Integer) The maximum number of autoscaling addresses allowed per NAT Gateway.

* `maximum_reserved_addresses_per_natgateway` - (Integer) The maximum number of reserved addresses allowed per NAT Gateway.
