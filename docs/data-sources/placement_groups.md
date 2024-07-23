---
page_title: "Linode: linode_placement_groups"
description: |-
  Lists Linode Placement Groups on your account.
---

# Data Source: linode\_placement\_groups

**NOTE: Placement Groups may not currently be available to all users.**

Provides information about a list of Linode Placement Groups that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-placement-groups).

## Example Usage

The following example shows how one might use this data source to list Placement Groups.

```hcl
data "linode_placement_groups" "all" {}

data "linode_placement_groups" "filtered" {
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-pgs" {
  value = data.linode_placement_groups.all.placement_groups
}

output "filtered-pgs" {
  value = data.linode_placement_groups.filtered.placement_groups
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Placement Groups that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Placement Group will be stored in the `placement_groups` attribute and will export the following attributes:

* `label` - The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

* `region` - The region of the Placement Group.

* `affinity_type` - The affinity policy to use when placing Linodes in this group.

* `is_strict` - Whether Linodes must be able to become compliant during assignment. (Default `true`)

* `is_compliant` - Whether all Linodes in this group are currently compliant with the group's affinity policy.

* `members` - A set of Linodes currently assigned to this Placement Group.

  * `linode_id` - The ID of the Linode.

  * `is_compliant` - Whether this Linode is currently compliant with the group's affinity policy.

## Filterable Fields

* `id`

* `label`

* `region`

* `affinity_type`

* `is_strict`

* `is_compliant`
