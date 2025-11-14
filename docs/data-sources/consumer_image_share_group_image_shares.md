---
page_title: "Linode: linode_consumer_image_share_group_image_shares"
description: |-
  Lists Images shared in the Image Share Group the provided Token has been accepted into.
---

# Data Source: linode\_consumer\_image\_share\_group\_image\_shares

Provides information about a list of Images that match a set of filters that have been
shared in the Image Share Group that the provided Token has been accepted into.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-sharegroup-images-by-token). May not be currently available to all users even under v4beta.

## Example Usage

The following example shows how one might use this data source to list Images shared in an Image Share Group.

```hcl
data "linode_consumer_image_share_group_image_shares" "all" {}

data "linode_consumer_image_share_group_image_shares" "filtered" {
    token_uuid = "54e1adf3-e499-4685-82be-10d29d4e8fae"
    filter {
        name = "label"
        values = ["my-label"]
    }
}

output "all-shared-images" {
  value = data.linode_consumer_image_share_group_image_shares.all.image_shares
}

output "filtered-shared-images" {
  value = data.linode_consumer_image_share_group_image_shares.filtered.image_shares
}
```

## Argument Reference

The following arguments are supported:

* `token_uuid` - (Required) The UUID of the Token that has been accepted into the Image Share Group to list shared Images from.

* [`filter`](#filter) - (Optional) A set of filters used to select Image Share Groups that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Image Share will be stored in the `images_shares` attribute and will export the following attributes:

* `id` - The unique ID assigned to this Image Share.

* `label` - The label of the Image Share.

* `capabilities` - The capabilities of the Image represented by the Image Share.

* `created` - When this Image Share was created.

* `deprecated` - Whether this Image is deprecated.

* `description` - A description of the Image Share.

* `is_public` - True if the Image is public.

* `image_sharing` - Details about image sharing, including who the image is shared with and by.
  * `shared_with` - Details about who the image is shared with.
    * `sharegroup_count` - The number of sharegroups the private image is present in.
    * `sharegroup_list_url` - The GET api url to view the sharegroups in which the image is shared.
  * `shared_by` - Details about who the image is shared by.
    * `sharegroup_id` - The sharegroup_id from the im_ImageShare row.
    * `sharegroup_uuid` - The sharegroup_uuid from the im_ImageShare row.
    * `sharegroup_label` - The label from the associated im_ImageShareGroup row.
    * `source_image_id` - The image id of the base image (will only be shown to producers, will be null for consumers).

* `size` - The minimum size this Image needs to deploy. Size is in MB. example: 2500

* `status` - The current status of this image. (`creating`, `pending_upload`, `available`)

* `type` - How the Image was created. Manual Images can be created at any time. "Automatic" Images are created automatically from a deleted Linode. (`manual`, `automatic`)

* `tags` - A list of customized tags.

* `total_size` - The total size of the image in all available regions.

## Filterable Fields

* `id`

* `label`
