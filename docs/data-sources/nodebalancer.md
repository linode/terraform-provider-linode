---
page_title: "Linode: linode_nodebalancer"
description: |-
  Provides details about a NodeBalancer.
---

# Data Source: linode\_nodebalancer

Provides details about a Linode NodeBalancer.

## Example Usage

```terraform
data "linode_nodebalancer" "my-nodebalancer" {
    id = 123
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The NodeBalancer's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the Linode NodeBalancer

* `client_conn_throttle` - Throttle connections per second (0-20).

* `created` – When this Linode NodeBalancer was created

* `linode_id` - The ID of a Linode Instance where the NodeBalancer should be attached.

* `tags` - A list of tags applied to this object. Tags are for organizational purposes only.

* `hostname` - This NodeBalancer's hostname, ending with .ip.linodeusercontent.com

* `ipv4` - The Public IPv4 Address of this NodeBalancer

* `ipv6` - The Public IPv6 Address of this NodeBalancer

* `region` - The Region where this Linode NodeBalancer is located. NodeBalancers only support backends in the same Region.

* [`transfer`](#transfer) - The network transfer stats for the current month

* `updated` – When this Linode NodeBalancer was last updated

* [`firewalls`](#firewalls) - A list of Firewalls assigned to this NodeBalancer.

### transfer

The following attributes are available on transfer:

* `in` - The total transfer, in MB, used by this NodeBalancer for the current month

* `out` - The total inbound transfer, in MB, used for this NodeBalancer for the current month

* `total` - The total outbound transfer, in MB, used for this NodeBalancer for the current month

### firewalls

The following attributes are available on firewalls:

* `id` - The Firewall's ID.

* `label` - The label for the firewall.

* `tags` - The tags applied to the firewall.

* [`inbound`](#inbound-and-outbound) - A firewall rule that specifies what inbound network traffic is allowed.

* `inbound_policy` - The default behavior for inbound traffic. (`ACCEPT`, `DROP`)

* [`outbound`](#inbound-and-outbound) - A firewall rule that specifies what outbound network traffic is allowed.

* `outbound_policy` - The default behavior for outbound traffic. (`ACCEPT`, `DROP`)

* `status` - The status of the firewall. (`enabled`, `disabled`, `deleted`)

* `created` - When this firewall was created.

* `updated` - When this firewall was last updated.

#### inboud and outbound

The following arguments are supported in the inbound and outbound rule blocks:

* `label` - Used to identify this rule. For display purposes only.

* `action` - Controls whether traffic is accepted or dropped by this rule. Overrides the Firewall’s inbound_policy if this is an inbound rule, or the outbound_policy if this is an outbound rule.

* `protocol` - The network protocol this rule controls. (`TCP`, `UDP`, `ICMP`)

* `ports` - A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - A list of IPv4 addresses or networks. Must be in IP/mask format.

* `ipv6` - A list of IPv6 addresses or networks. Must be in IP/mask format.
