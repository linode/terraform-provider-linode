---
page_title: "Linode: linode_firewall_template"
description: |-
  Provides details about a Linode Firewall Template.
---

# Data Source: linode\_firewall\_template

Provides information about a Linode Firewall Template.

## Example Usage

The following example shows how one might use this data source to access information about a specific Firewall Template:

```hcl
data "linode_firewall_template" "public-template" {
  slug = "public"
}

output "firewall_template_id" {
  value = data.linode_firewall_template.public-template.id
}
```

## Argument Reference

The following arguments are supported:

* `slug` - (Required) The slug of the firewall template.

## Attributes Reference

The following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* `id` - The computed ID of the data source, which matches the `slug` attribute.
* `inbound` - (Read-Only Object List) A list of firewall rules specifying allowed inbound network traffic. Referenced with an index (e.g. `inbound.0.action`).
* `inbound_policy` - The default behavior for inbound traffic. This can be overridden by individual firewall rules.
* `outbound` - (Read-Only Object List) A list of firewall rules specifying allowed outbound network traffic. Referenced with an index (e.g. `outbound.0.action`).
* `outbound_policy` - The default behavior for outbound traffic. This can be overridden by individual firewall rules.
