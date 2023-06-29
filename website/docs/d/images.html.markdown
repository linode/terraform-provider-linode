---
layout: "linode"
page_title: "Linode: linode_images"
sidebar_current: "docs-linode-datasource-images"
description: |-
  Provides information about Linode images that match a set of filters.
---

# Data Source: linode\_images

Provides information about Linode images that match a set of filters.

## Example Usage

Get information about all Linode images with a certain label and visibility:

```hcl
data "linode_images" "specific-images" {
  filter {
    name = "label"
    values = ["Debian 11"]
  }

  filter {
    name = "is_public"
    values = ["true"]
  }
}

output "image_id" {
  value = data.linode_images.specific-images.images.0.id
}
```

Get information about all Linode images associated with the current token:

```hcl
data "linode_images" "all-images" {}

output "image_ids" {
  value = data.linode_images.all-images.images.*.id
}
```

## Argument Reference

The following arguments are supported:

* `latest` - (Optional) If true, only the latest image will be returned. Images without a valid `created` field are not included in the result.

* [`filter`](#filter) - (Optional) A set of filters used to select Linode images that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode image will be stored in the `images` attribute and will export the following attributes:

* `id` - The unique ID of this Image.  The ID of private images begin with `private/` followed by the numeric identifier of the private image, for example `private/12345`.

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

## Filterable Fields

* `created_by`

* `deprecated`

* `description`

* `id`

* `is_public`

* `label`

* `size`

* `status`

* `vendor`
