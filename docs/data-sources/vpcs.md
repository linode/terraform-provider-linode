---
page_title: "Linode: linode_vpcs"
description: |-
  Lists Linode VPCs on your account.
---

# Data Source: linode\_vpcs

Provides information about a list of Linode VPCs that match a set of filters.

## Example Usage

The following example shows how one might use this data source to list VPCs.

```hcl
data "linode_vpcs" "filtered-vpcs" {
    filter {
        name = "label"
        values = ["test"]
    }
}

output "vpcs" {
  value = data.linode_vpcs.filtered-vpcs.vpcs
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode VPCs that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VPC will be stored in the `vpcs` attribute and will export the following attributes:

* `id` - The unique id of this VPC.

* `label` - The label of the VPC.

* `description` - The user-defined description of this VPC.

* `region` - The region where the VPC is deployed.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was updated.

## Filterable Fields

* `id`

* `label`

* `description`

* `region`
