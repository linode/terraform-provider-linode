---
page_title: "Linode: linode_account_availabilities"
description: |-
  Provides information about services which are unavailable for the current Linode account.
---

# linode\_account\_availabilities

Provides information about services which are unavailable for the current Linode account.

## Example Usage

The following example shows how one might use this data source to discover regions without specific service availability.

```hcl
data "linode_account_availabilities" "filtered-availabilities" {
  filter {
    name = "unavailable"
    values = ["Linodes"]
  }
}

output "regions-without-linodes" {
  value = data.linode_account_availabilities.filtered-availabilities.availabilities.*.region
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode account availabilities that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode account availability will be stored in the `availabilities` attribute and will export the following attributes:

* `region` - The region this availability entry refers to.

* `unavailable` - A set of services that are unavailable for the given region.

## Filterable Fields

* `region`

* `unavailable`
