---
page_title: "Linode: linode_producer_image_share_groups"
description: |-
  Lists Image Share Groups on your account.
---

# Data Source: linode\_producer\_image\_share\_groups

Provides information about a list of Image Share Groups that match a set of filters.
For more information, see the [Linode APIv4 docs](TODO).

## Example Usage

The following example shows how one might use this data source to list Image Share Groups.

```hcl
data "linode_producer_image_share_groups" "all" {}

data "linode_producer_image_share_groups" "filtered" {
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-share-groups" {
  value = data.linode_producer_image_share_groups.all.image_share_groups
}

output "filtered-share-groups" {
  value = data.linode_producer_image_share_groups.filtered.image_share_groups
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Image Share Groups that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Image Share Group.

* `uuid` - The UUID of the Image Share Group.

* `label` - The label of the Image Share Group.

* `description` - The description of the Image Share Group.

* `is_suspended` - Whether the Image Share Group is suspended.

* `images_count` - The number of images in the Image Share Group.

* `members_count` - The number of members in the Image Share Group.

* `created` - The date and time the Image Share Group was created.

* `updated` - The date and time the Image Share Group was last updated.

* `expiry` - The date and time the Image Share Group will expire.

## Filterable Fields

* `id`

* `label`

* `is_suspended`
