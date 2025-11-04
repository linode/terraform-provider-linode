---
page_title: "Linode: linode_consumer_image_share_group_token"
description: |-
  Manages a token for an Image Share Group.
---

# linode\_consumer\_image\_share\_group\_token

Manages a token for an Image Share Group.
For more information, see the [Linode APIv4 docs](TODO). May not be currently available to all users even under v4beta.

## Example Usage

Create a token for an Image Share Group:

```terraform
resource "linode_consumer_image_share_group_token" "example" {
  valid_for_sharegroup_uuid = "03fbb93e-c27d-4c4a-9180-67f6e0cd74ca"
  label                     = "example-token"
}
```

## Argument Reference

The following arguments are supported:

* `valid_for_sharegroup_uuid` - (Required) The UUID of the Image Share Group for which to create a token.

* `label` - (Optional) A label for the token.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `token` - The one-time-use token to be provided to the Image Share Group Producer.

* `token_uuid` - The UUID of the token.

* `status` - The status of the token.

* `created` - When the token was created.

* `updated` - When the token was last updated.

* `expiry` - When the token will expire.

* `sharegroup_uuid` - The UUID of the Image Share Group that the token is for.

* `sharegroup_label` - The label of the Image Share Group that the token is for.
