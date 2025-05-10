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

* `id` - The computed ID of the data source, which matches the `slug` attribute.
* `inbound` - A list of firewall rules specifying allowed inbound network traffic.
* `inbound_policy` - The default behavior for inbound traffic. This can be overridden by individual firewall rules.
* `outbound` - A list of firewall rules specifying allowed outbound network traffic.
* `outbound_policy` - The default behavior for outbound traffic. This can be overridden by individual firewall rules.
