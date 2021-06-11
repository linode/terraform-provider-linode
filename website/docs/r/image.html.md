---
layout: "linode"
page_title: "Linode: linode_image"
sidebar_current: "docs-linode-resource-image"
description: |-
  Manages a Linode Image.
---

# linode\_image

Provides a Linode Image resource.  This can be used to create, modify, and delete Linodes Images.  Linode Images are snapshots of a Linode Instance Disk which can then be used to provision more Linode Instances.  Images can be used across regions.

For more information, see [Linode's documentation on Images](https://www.linode.com/docs/platform/disk-images/linode-images/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createImage).

## Example Usage

Creating an image from an existing Linode Instance and deploying another instance with that image:

```hcl
resource "linode_instance" "foo" {
    type = "g6-nanode-1"
    region = "us-central"
}

resource "linode_image" "bar" {
    label = "foo-sda-image"
    description = "Image taken from foo"
    disk_id = linode_instance.foo.disk.0.id
    linode_id = linode_instance.foo.id
}

resource "linode_instance" "bar_based" {
    type = linode_instance.foo.type
    region = "eu-west"
    image = linode_image.bar.id
}
```

Creating and uploading an image from a local file:

```hcl
resource "linode_image" "foobar" {
    label = "foobar-image"
    description = "An image uploaded from Terraform!"
    region = "us-southeast"
  
    file_path = "path/to/image.img.gz"
    file_hash = filemd5("path/to/image.img.gz")
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) A short description of the Image. Labels cannot contain special characters.

* `description` - (Optional) A detailed description of this Image.

- - -

The following arguments apply to creating an image from an existing Linode Instance:

* `disk_id` - (Required) The ID of the Linode Disk that this Image will be created from.

* `linode_id` - (Required) The ID of the Linode that this Image will be created from.

- - -

~> **NOTICE:** Uploading images is currently in beta. Ensure `LINODE_API_VERSION` is set to `v4beta` in order to use this functionality.

The following arguments apply to uploading an image:

* `file_path` - (Required) The path of the image file to be uploaded.

* `file_hash` - (Optional) The MD5 hash of the file to be uploaded. This is used to trigger file updates.

* `region` - (Required) The region of the image. See all regions [here](https://api.linode.com/v4/regions).

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the instance image (until the instance is available)

## Attributes

This resource exports the following attributes:

* `id` - The unique ID of this Image.  The ID of private images begin with `private/` followed by the numeric identifier of the private image, for example `private/12345`.

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
