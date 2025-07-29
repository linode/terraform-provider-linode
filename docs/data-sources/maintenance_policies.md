---
page_title: "Linode: linode_maintenance_policies"
description: |-
  Provides details about the Maintenance Policies available to apply to Accounts and Instances.
---

# linode\_maintenance\_policies

Provides details about the Maintenance Policies available to apply to Accounts and Instances.
For more information, see the [Linode APIv4 docs](TODO).

## Example Usage

The following example shows how one might use this data source to access information about Maintenance Policies:

```hcl
data "linode_maintenance_policies" "example" {}

output "example_output" {
  value = data.linode_maintenance_policies.example
}
```

The following example shows how one might use this data source to access information about Maintenance Policies with a certain slug:

```hcl
data "linode_maintenance_policies" "example-specific-slug" {
  filter {
    name = "slug"
    values = "linode/migrate"
  }
}

output "example_slug_output" {
  value = data.linode_maintenance_policies.example-specific-slug
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode LKE types that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Maintenance Policy will be stored in the `maintenance_policies` attribute and will export the following attributes:

* `slug` - Unique identifier for this policy

* `label` - The label for this policy.

* `description` - Description of this policy

* `type` - The type of action taken during maintenance.

* `notification_period_sec` - The notification lead time in seconds.

* `is_default` - Whether this is the default policy for the account.

## Filterable Fields

* `slug`
