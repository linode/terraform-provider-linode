---
page_title: "Linode: linode_producer_image_share_group_member"
description: |-
  Manages a member of an Image Share Group.
---

# linode\_producer\_image\_share\_group\_member

Manages a member of an Image Share Group.
For more information, see the [Linode APIv4 docs](TODO). May not be currently available to all users even under v4beta.

## Example Usage

Accept a member into an Image Share Group:

```terraform
resource "linode_producer_image_share_group_member" "example" {
  sharegroup_id   = 12345
  token  = "abcdefghijklmnopqrstuvwxyz0123456789"
  label = "example-member"
}
```

## Argument Reference

The following arguments are supported:

* `sharegroup_id` - (Required) The ID of the Image Share Group to which the member will be added.

* `token` - (Required) The token of the prospective member.

* `label` - (Required) A label for the member.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `token_uuid` - The UUID of member's token.

* `status` - The status of the member.

* `created` - When the member was created.

* `updated` - When the member was last updated.

* `expiry` - When the member will expire.
