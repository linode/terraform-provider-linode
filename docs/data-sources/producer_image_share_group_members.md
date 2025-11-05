---
page_title: "Linode: linode_producer_image_share_group_members"
description: |-
  Lists an Image Share Group's Members on your account.
---

# Data Source: linode\_producer\_image\_share\_group\_members

Provides information about a list of Members of an Image Share Group that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-sharegroup-members). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how one might use this data source to list Image Share Group Members.

```hcl
data "linode_producer_image_share_group_members" "all" {
    sharegroup_id = 12345
}

data "linode_producer_image_share_group_members" "filtered" {
    sharegroup_id = 12345
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-share-group-members" {
  value = data.linode_producer_image_share_group_members.all.members
}

output "filtered-share-group-members" {
  value = data.linode_producer_image_share_group_members.filtered.members
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Image Share Groups that meet certain requirements.

* `sharegroup_id` - (Required) The ID of the Image Share Group for which to list members.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `token_uuid` - The UUID of member's token.

* `label` - The label of the member.

* `status` - The status of the member.

* `created` - When the member was created.

* `updated` - When the member was last updated.

* `expiry` - When the member will expire.

## Filterable Fields

* `token_uuid`

* `label`

* `status`
