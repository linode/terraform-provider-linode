---
page_title: "Linode: linode_firewall_rulesets"
description: |-
  Provides information about Firewall Rule Sets that match a set of filters.
---

# Data Source: linode\_firewall\_rulesets

Provides information about Linode Firewall Rule Sets that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-firewall-rule-sets).

## Example Usage

Get information about all inbound rule sets:

```terraform
data "linode_firewall_rulesets" "inbound" {
  filter {
    name   = "type"
    values = ["inbound"]
  }
}

output "ruleset_labels" {
  value = data.linode_firewall_rulesets.inbound.rulesets.*.label
}
```

Get all rule sets:

```terraform
data "linode_firewall_rulesets" "all" {}

output "ruleset_ids" {
  value = data.linode_firewall_rulesets.all.rulesets.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Firewall Rule Sets that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Firewall Rule Set will be stored in the `rulesets` attribute and will export the following attributes:

* `label` - The label for the Rule Set.

* `description` - The description of the Rule Set.

* `type` - The type of rule set (`inbound` or `outbound`).

* [`rules`](#rules) - The firewall rules defined in this set.

* `is_service_defined` - Whether this Rule Set is service-defined (managed by Linode).

* `version` - The version number of the Rule Set.

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

## Filterable Fields

* `label`

* `type`
