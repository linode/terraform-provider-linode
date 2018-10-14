---
layout: "linode"
page_title: "Linode: linode_image"
sidebar_current: "docs-linode-resource-image"
description: |-
  Manages a Linode Image.
---

# linode\_image

Provides a Linode Image resource.  This can be used to create,
modify, and delete Linodes Images. For more information, see [Linode's documentation on Images](https://www.linode.com/docs/platform/disk-images/linode-images/)
and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createImage).

## Example Usage

The following example shows how one might use this resource to configure an Image from a Linode Instance Disk.

```hcl
resource "linode_instance" "foobaz" {
    root_pass = "3X4mp13"
    type = "g6-nanode-1"
    region = "us-west"
}

resource "linode_image" "foobar" {
    label = "foo-volume"
    description = "My new disk image"
    disk_id = "${linode_instance.foobaz.disk.0.id}"
    linode_id = "${linode_instance.foobaz.id}"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) A short description of the Image. Labels cannot contain special characters.

* `disk_id` - (Required) The ID of the Linode Disk that this Image will be created from.

* `linode_id` - (Required) The ID of the Linode that this Image will be created from.

- - -

* `description` - (Optional) A detailed description of this Image.

## Attributes

This resource exports the following attributes:

* `created` - When this Image was created.

* `created_by` - The name of the User who created this Image.

* `deprecated` - Whether or not this Image is deprecated. Will only be True for deprecated public Images.

* `is_public` - True if the Image is public.

* `size` - The minimum size this Image needs to deploy. Size is in MB.

* `type` - How the Image was created. 'Manual' Images can be created at any time. 'Automatic' images are created automatically from a deleted Linode.

* `expiry` - Only Images created automatically (from a deleted Linode; type=automatic) will expire.

* `vendor` - The upstream distribution vendor. Nil for private Images.

## Import

Linodes Images can be imported using the Linode Image `id`, e.g.

```sh
terraform import linode_image.myimage 1234567
```
