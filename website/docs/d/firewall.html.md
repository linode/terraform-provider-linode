---
layout: "linode"
page_title: "Linode: linode_firewall"
sidebar_current: "docs-linode-datasource-firewall"
description: |-
Provides details about a Firewall.
---

# Data Source: linode\_firewall

Provides details about a Linode Firewall.

## Example Usage

```terraform
data "linode_firewall" "my-firewall" {
    id = 123
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The Firewall's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label for the firewall.

* `tags` - The tags applied to the firewall.

* `disabled` - If true, the firewall is inactive.

* [`inbound`](#inbound-and-outbound) - A firewall rule that specifies what inbound network traffic is allowed.

* `inbound_policy` - The default behavior for inbound traffic. (`ACCEPT`, `DROP`)

* [`outbound`](#inbound-and-outbound) - A firewall rule that specifies what outbound network traffic is allowed.

* `outbound_policy` - The default behavior for outbound traffic. (`ACCEPT`, `DROP`)

* `linodes` - The IDs of Linodes to apply this firewall to.

* `status` - The status of the firewall. (`enabled`, `disabled`, `deleted`)

* [`devices`](#devices) - The devices governed by the Firewall.

### inbound and outbound

The following arguments are supported in the inbound and outbound rule blocks:

* `label` - Used to identify this rule. For display purposes only.

* `action` - Controls whether traffic is accepted or dropped by this rule. Overrides the Firewall’s inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.

* `protocol` - The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`)

* `ports` - A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - A list of IPv4 addresses or networks. Must be in IP/mask format.

* `ipv6` - A list of IPv6 addresses or networks. Must be in IP/mask format.

### devices

The following attributes are available on devices:

* `id` - The ID of the Firewall Device.

* `entity_id` - The ID of the underlying entity this device references (i.e. the Linode's ID).

* `type` - The type of Firewall Device.

* `label` - The label of the underlying entity this device references.

* `url` The URL of the underlying entity this device references.
