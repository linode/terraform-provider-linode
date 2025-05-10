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

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Cloud Firewalls that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

The following attributes are exported:

* `templates` - A list of firewall templates, where each template includes:
  * `slug` - The slug of the firewall template.
  * `inbound` - A list of firewall rules specifying allowed inbound network traffic.
  * `inbound_policy` - The default behavior for inbound traffic.
  * `outbound` - A list of firewall rules specifying allowed outbound network traffic.
  * `outbound_policy` - The default behavior for outbound traffic.

## Filterable Fields

* `slug`
