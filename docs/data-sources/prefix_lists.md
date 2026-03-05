---
page_title: "Linode: linode_prefix_lists"
description: |-
  Provides information about Prefix Lists that match a set of filters.
---

# Data Source: linode\_prefix\_lists

Provides information about Linode Prefix Lists that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-prefix-lists).

## Example Usage

Get information about all account-visible prefix lists:

```terraform
data "linode_prefix_lists" "account" {
  filter {
    name   = "visibility"
    values = ["account"]
  }
}

output "prefix_list_names" {
  value = data.linode_prefix_lists.account.prefix_lists.*.name
}
```

Get all prefix lists:

```terraform
data "linode_prefix_lists" "all" {}

output "prefix_list_ids" {
  value = data.linode_prefix_lists.all.prefix_lists.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Prefix Lists that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Prefix List will be stored in the `prefix_lists` attribute and will export the following attributes:

* `name` - The name of the Prefix List.

* `description` - A description of the Prefix List.

* `visibility` - The visibility of the Prefix List. (`account`, `restricted`)

* `source_prefixlist_id` - The ID of the source prefix list, if this is a derived list.

* `ipv4` - A list of IPv4 addresses or networks in CIDR format contained in this prefix list.

* `ipv6` - A list of IPv6 addresses or networks in CIDR format contained in this prefix list.

* `version` - The version number of this Prefix List.

* `created` - When this Prefix List was created.

* `updated` - When this Prefix List was last updated.

## Filterable Fields

* `name`

* `visibility`
