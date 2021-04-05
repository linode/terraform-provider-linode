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
    values = ["Debian 8"]
  }

  filter {
    name = "is_public"
    values = ["true"]
  }
}
```

Get information about all Linode images associated with the current token:

```hcl
data "linode_images" "all-images" {}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode images that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

## Attributes

Each Linode image will be stored in the `images` attribute and will export the following attributes:

* `id` - The unique ID of this Image.  The ID of private images begin with `private/` followed by the numeric identifier of the private image, for example `private/12345`.

* `label` - A short description of the Image.

* `created` - When this Image was created.

* `created_by` - The name of the User who created this Image, or "linode" for official Images.

* `deprecated` - Whether or not this Image is deprecated. Will only be true for deprecated public Images.

* `description` - A detailed description of this Image.

* `is_public` - True if the Image is public.

* `size` - The minimum size this Image needs to deploy. Size is in MB. example: 2500

* `type` - How the Image was created. Manual Images can be created at any time. "Automatic" Images are created automatically from a deleted Linode.

* `vendor` - The upstream distribution vendor. `None` for private Images.

## Filterable Fields

* `deprecated`

* `is_public`

* `label`

* `size`

* `vendor`
