---
page_title: "Linode: linode_firewall_ruleset"
description: |-
  Provides details about a Firewall Rule Set.
---

# Data Source: linode\_firewall\_ruleset

Provides details about a Linode Firewall Rule Set.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-firewall-rule-set).

## Example Usage

```terraform
data "linode_firewall_ruleset" "example" {
  id = "12345"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Firewall Rule Set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label for the Rule Set.

* `description` - The description of the Rule Set.

* `type` - The type of rule set (`inbound` or `outbound`).

* [`rules`](#rules) - The firewall rules defined in this set.

* `is_service_defined` - Whether this Rule Set is service-defined (managed by Linode).

* `version` - The version number of this Rule Set.

* `created` - When this Rule Set was created.

* `updated` - When this Rule Set was last updated.

### rules

Each rule exports the following attributes:

* `label` - The label for this rule.

* `action` - Controls whether traffic is accepted or dropped by this rule. (`ACCEPT`, `DROP`)

* `protocol` - The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`, `IPENCAP`)

* `description` - The description for this rule.

* `ports` - A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - A list of IPv4 addresses or networks in CIDR format, or prefix list tokens.

* `ipv6` - A list of IPv6 addresses or networks in CIDR format, or prefix list tokens.
