---
page_title: "Linode: linode_firewall_templates"
description: |-
  Lists Linode Firewall Templates available on your account.
---

# Data Source: linode\_firewall\_templates

Provides information about all Linode Firewall Templates.

## Example Usage

The following example shows how one might use this data source to list all available Firewall Templates:

```hcl
data "linode_firewall_templates" "all" {}

output "firewall_template_slugs" {
  value = data.linode_firewall_templates.all.firewall_templates
}
```

Or with some filters to get a subset of the results.

```hcl
data "linode_firewall_templates" "filtered" {
  filter {
    name     = "slug"
    values   = ["public"]
    match_by = "exact"
  }
}

output "firewall_template_slugs" {
  value = data.linode_firewall_templates.filtered.firewall_templates
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode Cloud Firewalls that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

The following attributes are exported:

* `firewall_templates` - (Nested Attribute List) The returned list of firewall templates. Referenced by index (e.g. `firewall_templates[0].slug`).

* `templates` - A list of firewall templates, where each template includes:
  * `slug` - The slug of the firewall template.
  * `inbound` - (Read-Only Object List) A list of firewall rules specifying allowed inbound network traffic. Referenced with an index (e.g. `inbound.0.action`).
  * `inbound_policy` - The default behavior for inbound traffic.
  * `outbound` - (Read-Only Object List) A list of firewall rules specifying allowed outbound network traffic. Referenced with an index (e.g. `outbound.0.action`).
  * `outbound_policy` - The default behavior for outbound traffic.

## Filterable Fields

* `slug`
