---
title: "linode_firewall_settings"
description: |-
  Manages Linode account-level firewall settings.
---

# linode_firewall_settings

Manages Linode account-level firewall settings. Resetting default firewall IDs
to null is not available to all customers and unsupported in this resource.

## Example Usage

```hcl
resource "linode_firewall_settings" "example" {
  default_firewall_ids = {
    linode           = 12345
    nodebalancer     = 12345
    public_interface = 12345
    vpc_interface    = 12345
  }
}
```

## Argument Reference

* `default_firewall_ids` - (Optional) A map of default firewall IDs for various interfaces.
  * `linode` - (Optional) The Linode's default firewall.
  * `nodebalancer` - (Optional) The NodeBalancer's default firewall.
  * `public_interface` - (Optional) The public interface's default firewall.
  * `vpc_interface` - (Optional) The VPC interface's default firewall.

## API Reference

See the [Linode API documentation](https://techdocs.akamai.com/linode-api/reference/put-firewall-settings) for more details.
