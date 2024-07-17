---
page_title: "Linode: linode_firewalls"
description: |-
  Provides information about Linode Cloud Firewalls that match a set of filters.
---

# Data Source: linode\_firewalls

Provides information about Linode Cloud Firewalls that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-firewalls).

## Example Usage

Get information about all Linode Cloud Firewalls with a certain label and visibility:

```hcl
data "linode_firewalls" "specific" {
  filter {
    name = "label"
    values = ["my-firewalls"]
  }

  filter {
    name = "tags"
    values = ["my-tag"]
  }
}

output "firewall_id" {
  value = data.linode_firewalls.specific.firewalls.0.id
}
```

Get information about all Linode images associated with the current token:

```hcl
data "linode_firewalls" "all" {}

output "firewall_ids" {
  value = data.linode_firewalls.all.firewalls.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Cloud Firewalls that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode image will be stored in the `firewalls` attribute and will export the following attributes:

* `id` - The unique ID assigned to this Firewall.

* `label` - The label for the Firewall. For display purposes only. If no label is provided, a default will be assigned.

* `tags` - An array of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

* `disabled` - If true, the Firewall is inactive.

* [`devices`](#firewall-device) - The devices associated with this firewall.

* [`inbound`](#firewall-rule) - A set of firewall rules that specify what inbound network traffic is allowed.

* `inbound_policy` - The default behavior for inbound traffic.

* [`outbound`](#firewall-rule) - A set of firewall rules that specify what outbound network traffic is allowed.

* `outbound_policy` - The default behavior for outbound traffic.

* `linodes` - The IDs of Linodes this firewall is applied to.

* `status` - The status of the firewall.

* `created` - When this firewall was created.

* `updated` - When this firewall was last updated.

## Firewall Rule

* `label` - The label of this rule for display purposes only.

* `action` - Controls whether traffic is accepted or dropped by this rule (ACCEPT, DROP).

* `protocol` - The network protocol this rule controls. (TCP, UDP, ICMP)

* `ports` - A string representation of ports and/or port ranges (i.e. "443" or "80-90, 91").

* `ipv4` - A list of IPv4 addresses or networks in IP/mask format.

* `ipv6` - A list of IPv6 addresses or networks in IP/mask format.

## Firewall Device

* `id` - The unique ID of this Firewall Device assignment.

* `entity_id` - The ID of the underlying entity this device references.

* `type` - The type of the assigned entity.

* `label` - The label of the assigned entity.

* `url` - The URL of the assigned entity.

## Filterable Fields

* `id`

* `label`

* `status`

* `tags`
