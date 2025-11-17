---
page_title: "Linode: linode_producer_image_share_group"
description: |-
  Provides details about an Image Share Group created by a producer.
---

# Data Source: linode\_producer\_image\_share\_group

`linode_producer_image_share_group` provides details about an Image Share Group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-sharegroup). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how the datasource might be used to obtain additional information about an Image Share Group.

```hcl
data "linode_producer_image_share_group" "sg" {
  id = 12345
}
```

## Argument Reference

* `id` - (Required) The ID of the Image Share Group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - The UUID of the Image Share Group.

* `label` - The label of the Image Share Group.

* `description` - The description of the Image Share Group.

* `is_suspended` - Whether the Image Share Group is suspended.

* `images_count` - The number of images in the Image Share Group.

* `members_count` - The number of members in the Image Share Group.

* `created` - The date and time the Image Share Group was created.

* `updated` - The date and time the Image Share Group was last updated.

* `expiry` - The date and time the Image Share Group will expire.
