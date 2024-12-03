---
page_title: "Linode: linode_image"
description: |-
  Manages a Linode Image.
---

# linode\_image

Provides a Linode Image resource.  This can be used to create, modify, and delete Linodes Images.  Linode Images are snapshots of a Linode Instance Disk which can then be used to provision more Linode Instances.  Images can be used across regions.

For more information, see [Linode's documentation on Images](https://www.linode.com/docs/platform/disk-images/linode-images/) and the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-image).

## Example Usage

Creating an image from an existing Linode Instance and deploying another instance with that image:

```hcl
resource "linode_instance" "foo" {
    type = "g6-nanode-1"
    region = "us-central"
    image = "linode/ubuntu22.04"
    root_pass = "insecure-p4ssw0rd!!"
}

resource "linode_image" "bar" {
    label = "foo-sda-image"
    description = "Image taken from foo"
    disk_id = linode_instance.foo.disk.0.id
    linode_id = linode_instance.foo.id
    tags = ["image-tag", "test"]
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
    tags = ["image-tag", "test"]
  
    file_path = "path/to/image.img.gz"
    file_hash = filemd5("path/to/image.img.gz")
}
```

Upload and replicate an image from a local file:

```hcl
resource "linode_image" "foobar" {
    label = "foobar-image"
    description = "An image uploaded from Terraform!"
    region = "us-southeast"
    tags = ["image-tag", "test"]
  
    file_path = "path/to/image.img.gz"
    file_hash = filemd5("path/to/image.img.gz")
    
    replica_regions = ["us-southeast", "us-east", "eu-west"]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) A short description of the Image. Labels cannot contain special characters.

* `description` - (Optional) A detailed description of this Image.

* `tags` - (Optional) A list of customized tags.

* `replica_regions` - (Optional) A list of regions that customer wants to replicate this image in. At least one valid region is required and only core regions allowed. Existing images in the regions not passed will be removed. See Replicate an Image [here](https://techdocs.akamai.com/linode-api/reference/post-replicate-image) for more details.

* `wait_for_replications` - (Optional) Whether to wait for all image replications become `available`. Default to false.

- - -

The following arguments apply to creating an image from an existing Linode Instance:

* `disk_id` - (Required) The ID of the Linode Disk that this Image will be created from.

* `linode_id` - (Required) The ID of the Linode that this Image will be created from.

- - -

~> **NOTICE:** Uploading images is currently in beta. Ensure `LINODE_API_VERSION` is set to `v4beta` in order to use this functionality.

The following arguments apply to uploading an image:

* `file_path` - (Required) The path of the image file to be uploaded.

* `file_hash` - (Optional) The MD5 hash of the file to be uploaded. This is used to trigger file updates.

* `region` - (Required) The region of the image. See all regions [here](https://techdocs.akamai.com/linode-api/reference/get-regions).

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the instance image (until the instance is available)

## Attributes Reference

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

* `total_size` - The total size of the image in all available regions.

* `replications` - A list of image replications region and corresponding status.
  * `region` - The region of an image replica.
  * `status` - The status of an image replica.

## Import

Linodes Images can be imported using the Linode Image `id`, e.g.

```sh
terraform import linode_image.myimage 1234567
```
