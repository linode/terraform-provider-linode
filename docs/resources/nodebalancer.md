---
page_title: "Linode: linode_nodebalancer"
description: |-
  Manages a Linode NodeBalancer.
---

# linode\_nodebalancer

Provides a Linode NodeBalancer resource.  This can be used to create, modify, and delete Linodes NodeBalancers in Linode's managed load balancer service.
For more information, see [Getting Started with NodeBalancers](https://www.linode.com/docs/platform/nodebalancer/getting-started-with-nodebalancers/) and the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-node-balancer).

The Linode Guide, [Create a NodeBalancer with Terraform](https://www.linode.com/docs/applications/configuration-management/create-a-nodebalancer-with-terraform/), provides step-by-step guidance and additional examples.

## Example Usage

The following example shows how one might use this resource to configure a NodeBalancer.

```hcl
resource "linode_nodebalancer" "foobar" {
    label = "mynodebalancer"
    region = "us-east"
    client_conn_throttle = 20
    tags = ["foobar"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region where this NodeBalancer will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc. See all regions [here](https://api.linode.com/v4/regions).  *Changing `region` forces the creation of a new Linode NodeBalancer.*.

- - -

* `label` - (Optional) The label of the Linode NodeBalancer

* `client_conn_throttle` - (Optional) Throttle connections per second (0-20). Set to 0 (default) to disable throttling.

* `tags` - (Optional) A list of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

## Attributes Reference

This resource exports the following attributes:

* `hostname` - This NodeBalancer's hostname, ending with .nodebalancer.linode.com

* `ipv4` - The Public IPv4 Address of this NodeBalancer

* `ipv6` - The Public IPv6 Address of this NodeBalancer

* `created` - When this NodeBalancer was created

* `updated` - When this NodeBalancer was last updated.

* [`transfer`](#transfer) - The network transfer stats for the current month

* [`firewalls`](#firewalls) - A list of Firewalls assigned to this NodeBalancer.

### transfer

The following attributes are available on transfer:

* `in` - The total transfer, in MB, used by this NodeBalancer for the current month

* `out` - The total inbound transfer, in MB, used for this NodeBalancer for the current month

* `total` - The total outbound transfer, in MB, used for this NodeBalancer for the current month

### firewalls

The following attributes are available on firewalls:

* `id` - (Required) The Firewall's ID.

* `label` - The label for the firewall.

* `tags` - The tags applied to the firewall. Tags are case-insensitive and are for organizational purposes only.

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

## Import

Linodes NodeBalancers can be imported using the Linode NodeBalancer `id`, e.g.

```sh
terraform import linode_nodebalancer.mynodebalancer 1234567
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for NodeBalancers and other Linode resource types.
