---
layout: "linode"
page_title: "Linode: linode_vlan"
sidebar_current: "docs-linode-vlan"
description: |-
  Manages a Linode VLAN.
---

# linode\_vlan

~> **NOTICE:** The VLAN feature is currently available through early access. To learn more, see the [early access documentation](https://github.com/linode/terraform-provider-linode/tree/master/EARLY_ACCESS.md).

Manages a Linode VLAN.

## Example Usage

```terraform
resource "linode_vlan" "my_vlan" {
  description = "my VLAN"
  region      = "ca-central"
  linodes     = [linode_instance.my_instance.id]
  cidr        = "0.0.0.0/0"

}

resource "linode_instance" "my_instance" {
  label      = "my_instance"
  image      = "linode/ubuntu18.04"
  region     = "ca-central"
  type       = "g6-standard-1"
  root_pass  = "bogusPassword$"
  swap_size  = 256
}
```

## Argument Reference

The following arguments are supported:

* `region` - (required) The region of where the VLAN is deployed.

* `description` - (Optional) Description of the vlan for display purposes only.

* `linodes` - (Optional) A list of IDs of Linodes to attach to this VLAN.

* `cidr_block` - (Optional) The CIDR block for this VLAN.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* [`attached_linodes`](#attached_linodes) - The devices governed by the Firewall.

### attached_linodes

The following attributes are available on attached linodes:

* `id` - The ID of the Linode.

* `mac_address` - The mac address of the Linode.

* `ipv4_address` - The IPv4 address of the Linode.
