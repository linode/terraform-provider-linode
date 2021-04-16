---
layout: "linode"
page_title: "Linode: linode_firewall"
sidebar_current: "docs-linode-firewall"
description: |-
  Manages a Linode Firewall.
---

# linode\_firewall

~> **NOTICE:** The Firewall feature is currently available through early access. To learn more, see the [early access documentation](https://github.com/linode/terraform-provider-linode/tree/main/EARLY_ACCESS.md).

Manages a Linode Firewall.

## Example Usage

```terraform
resource "linode_firewall" "my_firewall" {
  label = "my_firewall"
  tags  = ["test"]

  inbound {
    label    = "allow-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["ff00::/8"]
  }
  
  inbound {
    label    = "allow-https"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["ff00::/8"]
  }
  
  inbound_policy = "DROP"

  outbound {
    label    = "reject-http"
    action   = "DROP"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["ff00::/8"]
  }
  
  outbound {
    label    = "reject-https"
    action   = "DROP"
    protocol = "TCP"
    ports    = "443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["ff00::/8"]
  }
  
  outbound_policy = "ACCEPT"

  linodes = [linode_instance.my_instance.id]
}

resource "linode_instance" "my_instance" {
  label      = "my_instance"
  image      = "linode/ubuntu18.04"
  region     = "us-southeast"
  type       = "g6-standard-1"
  root_pass  = "bogusPassword$"
  swap_size  = 256
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) This Firewall's unique label.

* `disabled` - (Optional) If `true`, the Firewall's rules are not enforced (defaults to `false`).

* [`inbound`](#inbound) - (Optional) A firewall rule that specifies what inbound network traffic is allowed.
  
* `inbound_policy` - (Required) The default behavior for inbound traffic. This setting can be overridden by updating the inbound.action property of the Firewall Rule.

* [`outbound`](#outbound) - (Optional) A firewall rule that specifies what outbound network traffic is allowed.
  
* `outbound_policy` - (Required) The default behavior for outbound traffic. This setting can be overridden by updating the action property for an individual Firewall Rule.

* `linodes` - (Optional) A list of IDs of Linodes this Firewall should govern it's network traffic for.

* `tags` - (Optional) A list of tags applied to the Kubernetes cluster. Tags are for organizational purposes only.

### inbound and outbound

The following arguments are supported in the inbound and outbound rule blocks:

* `label` - (required) Used to identify this rule. For display purposes only.
  
* `action` - (required) Controls whether traffic is accepted or dropped by this rule. Overrides the Firewall’s inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.

* `protocol` - (Required) The network protocol this rule controls.

* `ports` - (Optional) A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").
  
* `ipv4` - (Optional) A list of IPv4 addresses or networks. Must be in IP/mask format.

* `ipv6` - (Optional) A list of IPv6 addresses or networks. Must be in IP/mask format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Firewall.

* `status` - The status of the Firewall.

* [`devices`](#devices) - The devices governed by the Firewall.

### devices

The following attributes are available on devices:

* `id` - The ID of the Firewall Device.

* `entity_id` - The ID of the underlying entity this device references (i.e. the Linode's ID).

* `type` - The type of Firewall Device.

* `label` - The label of the underlying entity this device references.

* `url` The URL of the underlying entity this device references.

## Import

Firewalls can be imported using the `id`, e.g.

```sh
terraform import linode_firewall.my_firewall 12345
```
