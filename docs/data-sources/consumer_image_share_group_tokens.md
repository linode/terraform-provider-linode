---
page_title: "Linode: linode_consumer_image_share_group_tokens"
description: |-
  Lists Image Share Group Tokens on your account.
---

# Data Source: linode\_consumer\_image\_share\_group\_tokens

Provides information about a list of Image Share Group Tokens that match a set of filters.
For more information, see the [Linode APIv4 docs](TODO). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how one might use this data source to list Image Share Groups.

```hcl
data "linode_consumer_image_share_group_tokens" "all" {}

data "linode_consumer_image_share_group_tokens" "filtered" {
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-share-group-tokens" {
  value = data.linode_consumer_image_share_group_tokens.all.tokens
}

output "filtered-share-group-tokens" {
  value = data.linode_consumer_image_share_group_tokens.filtered.tokens
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

* `token_uuid` - The UUID of the token.

* `label` - A label for the token.

* `status` - The status of the token.

* `created` - When the token was created.

* `updated` - When the token was last updated.

* `expiry` - When the token will expire.

* `valid_for_sharegroup_uuid` - The UUID of the Image Share Group for which to create a token.

* `sharegroup_uuid` - The UUID of the Image Share Group that the token is for.

* `sharegroup_label` - The label of the Image Share Group that the token is for.

## Filterable Fields

* `token_uuid`

* `label`

* `status`

* `valid_for_sharegroup_uuid`

* `sharegroup_uuid`

* `sharegroup_label`
