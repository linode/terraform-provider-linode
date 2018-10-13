---
layout: "linode"
page_title: "Linode: linode_image"
sidebar_current: "docs-linode-resource-image"
description: |-
  Manages a Linode Image.
---

# linode\_image

Provides a Linode Image resource.  This can be used to create,
modify, and delete Linodes Volumes. For more information, see [Linode's documentation on Images](https://www.linode.com/docs/platform/disk-images/linode-images/)
and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createImage).

## Example Usage

The following example shows how one might use this resource to configure a Volume attached to a Linode instance.

```hcl
resource "linode_instance" "foobaz" {
    root_pass = "3X4mp13"
    type = "g6-nanode-1"
    region = "us-west"
}

resource "linode_image" "foobar" {
    label = "foo-volume"
    description = "My new disk image"
    disk_id = "${linode_instance.disk[0].id}"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Linode Image

* `disk_id` - (Required) The ID of the Disk to use for the Image.

- - -

* `description` - (Optional) Description of the image.

## Attributes

This resource exports the following attributes:

## Import

Linodes Images can be imported using the Linode Image `id`, e.g.

```sh
terraform import linode_image.myimage 1234567
```
