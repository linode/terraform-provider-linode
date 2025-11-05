---
page_title: "Linode: linode_consumer_image_share_group_token"
description: |-
  Provides details about a Token for an Image Share Group.
---

# Data Source: linode\_consumer\_image\_share\_group\_token

`linode_consumer_image_share_group_token` provides details about a Token for an Image Share Group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-sharegroup-token). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how the datasource might be used to obtain additional information about a Token for an Image Share Group.

```hcl
data "linode_consumer_image_share_group_token" "token" {
  token_uuid = "db58ab2e-3021-4b08-9426-8e456f6dd268"
}
```

## Argument Reference

The following arguments are supported:

* `token_uuid` - (Required) The UUID of the token.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `label` - A label for the token.

* `status` - The status of the token.

* `created` - When the token was created.

* `updated` - When the token was last updated.

* `expiry` - When the token will expire.

* `valid_for_sharegroup_uuid` - The UUID of the Image Share Group for which to create a token.

* `sharegroup_uuid` - The UUID of the Image Share Group that the token is for.

* `sharegroup_label` - The label of the Image Share Group that the token is for.
