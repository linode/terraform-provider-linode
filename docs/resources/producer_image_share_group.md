---
page_title: "Linode: linode_producer_image_share_group"
description: |-
  Manages an Image Share Group.
---

# linode\_producer\_image\_share\_group

Manages an Image Share Group.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-sharegroups). May not be currently available to all users even under v4beta.

## Example Usage

Create an Image Share Group without any Images:

```terraform
resource "linode_producer_image_share_group" "test-empty" {
    label = "my-image-share-group"
    description = "My description."
}
```

Create an Image Share Group with one Image:

```terraform
resource "linode_producer_image_share_group" "test-images" {
    label = "my-image-share-group"
    description = "My description."
    images = [
        {
            id = "private/12345"
            label = "my-image"
            description = "My image description."
        },
    ]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Image Share Group.

* `description` - (Optional) The description of the Image Share Group

* [`images`](#images) - (Optional) A list of Images to include in the Image Share Group.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the Image Share Group.

* `uuid` - The UUID of the Image Share Group.

* `is_suspended` - Whether the Image Share Group is suspended.

* `images_count` - The number of images in the Image Share Group.

* `members_count` - The number of members in the Image Share Group.

* `created` - The date and time the Image Share Group was created.

* `updated` - The date and time the Image Share Group was last updated.

* `expiry` - The date and time the Image Share Group will expire.

### Images

Represents a single Image shared in an Image Share Group.

* `id` - (Required) The ID of the Image to share. This must be in the format `private/<image_id>`.

* `label` - (Optional) The label of the Image Share.

* `description` - (Optional) The description of the Image Share.
