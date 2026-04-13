---
page_title: "Linode: linode_firewall_rules_expansion"
description: |-
  Provides the expanded (resolved) firewall rules for a Firewall.
---

# Data Source: linode\_firewall\_rules\_expansion

Provides the expanded (resolved) firewall rules for a Linode Firewall. This data source resolves all prefix list tokens and rule set references into their concrete IP addresses and individual rules, giving you the effective rule set that the firewall is currently enforcing.

For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-firewall-rules-expansion).

## Example Usage

```terraform
resource "linode_firewall" "my_firewall" {
  label = "my-firewall"

  inbound_ruleset = [linode_firewall_ruleset.allow_web.id]

  inbound_policy  = "DROP"
  outbound_policy = "ACCEPT"

  linodes = [linode_instance.my_instance.id]
}

data "linode_firewall_rules_expansion" "expanded" {
  firewall_id = linode_firewall.my_firewall.id
}

output "effective_inbound_rules" {
  value = data.linode_firewall_rules_expansion.expanded.inbound
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Required) The ID of the Firewall to get the expanded rules for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* [`inbound`](#rules) - The expanded inbound firewall rules with all prefix list tokens and rule set references resolved.

* `inbound_policy` - The default behavior for inbound traffic. (`ACCEPT`, `DROP`)

* [`outbound`](#rules) - The expanded outbound firewall rules with all prefix list tokens and rule set references resolved.

* `outbound_policy` - The default behavior for outbound traffic. (`ACCEPT`, `DROP`)

* `version` - The version number of the Firewall's rule configuration.

### rules

Each expanded rule exports the following attributes:

* `label` - The label for this rule.

* `action` - Controls whether traffic is accepted or dropped by this rule. (`ACCEPT`, `DROP`)

* `protocol` - The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`, `IPENCAP`)

* `description` - The description for this rule.

* `ports` - A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - A list of resolved IPv4 addresses or networks in CIDR format.

* `ipv6` - A list of resolved IPv6 addresses or networks in CIDR format.
