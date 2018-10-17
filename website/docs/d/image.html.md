---
layout: "linode"
page_title: "Linode: linode_image"
sidebar_current: "docs-linode-datasource-image"
description: |-
  Provides details about a Linode image.
---

# Data Source: image

Provides information about a Linode image

## Example Usage

The following example shows how one might use this data source to define a Linode instance.

```hcl
data "linode_image" "k8_master" {
    id = "linode/debian8"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this Image. 

## Attributes

This Linode Instance resource exports the following attributes:

* `label` - A short description of the Image.

* `created` - When this Image was created.

* `created_by` - The name of the User who created this Image, or "linode" for official Images.

* `deprecated` - Whether or not this Image is deprecated. Will only be true for deprecated public Images.

* `description` - A detailed description of this Image.

* `is_public` - True if the Image is public.

* `size` - The minimum size this Image needs to deploy. Size is in MB. example: 2500

* `type` - How the Image was created. Manual Images can be created at any time. image"Automatic" Images are created automatically from a deleted Linode.

* `vendor` - The upstream distribution vendor. `None` for private Images. 
