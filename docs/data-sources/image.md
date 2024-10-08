---
page_title: "Linode: linode_image"
description: |-
  Provides details about a Linode image.
---

# Data Source: linode\_image

Provides information about a Linode image
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-image).

## Example Usage

The following example shows how one might use this data source to access information about a Linode image.

```hcl
data "linode_image" "k8_master" {
    id = "linode/debian12"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this Image.  The ID of private images begin with `private/` followed by the numeric identifier of the private image, for example `private/12345`.

## Attributes Reference

The Linode Image resource exports the following attributes:

* `label` - A short description of the Image.

* `created` - When this Image was created.

* `created_by` - The name of the User who created this Image, or "linode" for official Images.

* `deprecated` - Whether or not this Image is deprecated. Will only be true for deprecated public Images.

* `description` - A detailed description of this Image.

* `is_public` - True if the Image is public.

* `size` - The minimum size this Image needs to deploy. Size is in MB. example: 2500

* `status` - The current status of this image. (`creating`, `pending_upload`, `available`)

* `type` - How the Image was created. Manual Images can be created at any time. "Automatic" Images are created automatically from a deleted Linode. (`manual`, `automatic`)

* `vendor` - The upstream distribution vendor. `None` for private Images.

* `tags` - A list of customized tags.

* `total_size` - The total size of the image in all available regions.

* `replications` - A list of image replication regions and corresponding status.
  * `region` - The region of an image replica.
  * `status` - The status of an image replica.
