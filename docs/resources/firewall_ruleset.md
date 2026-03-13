---
page_title: "Linode: linode_firewall_ruleset"
description: |-
  Manages a Linode Firewall Rule Set.
---

# linode\_firewall\_ruleset

Manages a Linode Firewall Rule Set.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-firewall-rule-set).

## Example Usage

Create a reusable inbound rule set that allows SSH and HTTP traffic:

```terraform
resource "linode_firewall_ruleset" "allow_web_ssh" {
  label       = "allow-web-ssh"
  description = "Allow inbound SSH and HTTP traffic"
  type        = "inbound"

  rules {
    label    = "allow-ssh"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "22"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  rules {
    label    = "allow-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "80, 443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }
}
```

Reference a rule set in a firewall:

```terraform
resource "linode_firewall" "my_firewall" {
  label = "my-firewall"

  inbound_ruleset = [linode_firewall_ruleset.allow_web_ssh.id]

  inbound_policy  = "DROP"
  outbound_policy = "ACCEPT"

  linodes = [linode_instance.my_instance.id]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label for this Rule Set. Must be between 3 and 32 characters.

* `type` - (Required) The type of rule set. Must be `inbound` or `outbound`. Changing this forces a new resource to be created.

* `description` - (Optional) A description for this Rule Set.

* [`rules`](#rules) - (Optional) One or more rule blocks defining the firewall rules in this set.

### rules

The following arguments are supported in the `rules` block:

* `label` - (Required) The label for this rule. Must be between 3 and 32 characters. For display purposes only.

* `action` - (Required) Controls whether traffic is accepted or dropped by this rule. (`ACCEPT`, `DROP`)

* `protocol` - (Required) The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`, `IPENCAP`)

* `description` - (Optional) A description for this rule.

* `ports` - (Optional) A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - (Optional) A list of IPv4 addresses or networks in CIDR format, or prefix list tokens (e.g. `pl::subnets:123`).

* `ipv6` - (Optional) A list of IPv6 addresses or networks in CIDR format, or prefix list tokens (e.g. `pl::subnets:123`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Firewall Rule Set.

* `is_service_defined` - Whether this Rule Set is service-defined (managed by Linode).

* `version` - The version number of this Rule Set. This is incremented each time the rules are updated.

* `created` - When this Rule Set was created.

* `updated` - When this Rule Set was last updated.

## Import

Firewall Rule Sets can be imported using the `id`, e.g.

```sh
terraform import linode_firewall_ruleset.allow_web_ssh 12345
```
