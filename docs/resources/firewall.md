---
page_title: "Linode: linode_firewall"
description: |-
  Manages a Linode Firewall.
---

# linode\_firewall

Manages a Linode Firewall.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-firewalls).

## Example Usage

Accept only inbound HTTP(s) requests and drop outbound HTTP(s) requests:

```terraform
resource "linode_firewall" "my_firewall" {
  label = "my_firewall"

  inbound {
    label    = "allow-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
  
  inbound {
    label    = "allow-https"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
  
  inbound_policy = "DROP"

  outbound {
    label    = "reject-http"
    action   = "DROP"
    protocol = "TCP"
    ports    = "80"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
  
  outbound {
    label    = "reject-https"
    action   = "DROP"
    protocol = "TCP"
    ports    = "443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
  
  outbound_policy = "ACCEPT"

  linodes = [linode_instance.my_instance.id]
}

resource "linode_instance" "my_instance" {
  label      = "my_instance"
  image      = "linode/ubuntu22.04"
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
  
* `inbound_policy` - (Required) The default behavior for inbound traffic. This setting can be overridden by updating the inbound.action property of the Firewall Rule. (`ACCEPT`, `DROP`)

* [`outbound`](#outbound) - (Optional) A firewall rule that specifies what outbound network traffic is allowed.
  
* `outbound_policy` - (Required) The default behavior for outbound traffic. This setting can be overridden by updating the outbound.action property for an individual Firewall Rule. (`ACCEPT`, `DROP`)

* `linodes` - (Optional) A list of IDs of Linodes this Firewall should govern network traffic for.

* `nodebalancers` - (Optional) A list of IDs of NodeBalancers this Firewall should govern network traffic for.

* `tags` - (Optional) A list of tags applied to the Kubernetes cluster. Tags are case-insensitive and are for organizational purposes only.

### inbound and outbound

**NOTE:** Firewall rules can be dynamically generated using [dynamic blocks](https://www.terraform.io/language/expressions/dynamic-blocks).

The following arguments are supported in the inbound and outbound rule blocks:

* `label` - (required) Used to identify this rule. For display purposes only.
  
* `action` - (required) Controls whether traffic is accepted or dropped by this rule (`ACCEPT`, `DROP`). Overrides the Firewall’s inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.

* `protocol` - (Required) The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`)

* `ports` - (Optional) A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").
  
* `ipv4` - (Optional) A list of IPv4 addresses or networks. Must be in IP/mask (CIDR) format.

* `ipv6` - (Optional) A list of IPv6 addresses or networks. Must be in IP/mask (CIDR) format.

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
