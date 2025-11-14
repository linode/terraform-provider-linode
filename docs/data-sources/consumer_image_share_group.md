---
page_title: "Linode: linode_consumer_image_share_group"
description: |-
  Provides details about an Image Share Group a consumer's token has been accepted into.
---

# Data Source: linode\_consumer\_image\_share\_group

`linode_consumer_image_share_group` provides details about an Image Share Group that the user's token has been accepted into.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-sharegroup-by-token). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how the datasource might be used to obtain additional information about an Image Share Group.

```hcl
data "linode_consumer_image_share_group" "sg" {
  token_uuid = "7548d17e-8db4-4a91-b47c-a8e1203063d9"
}
```

## Argument Reference

* `token_uuid` - (Required) The UUID of the token that has been accepted into the Image Share Group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Image Share Group.

* `uuid` - The UUID of the Image Share Group.

* `label` - The label of the Image Share Group.

* `description` - The description of the Image Share Group.

* `is_suspended` - Whether the Image Share Group is suspended.

* `created` - The date and time the Image Share Group was created.

* `updated` - The date and time the Image Share Group was last updated.
