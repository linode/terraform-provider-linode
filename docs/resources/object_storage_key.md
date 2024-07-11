---
page_title: "Linode: linode_object_storage_key"
description: |-
  Manages a Linode Object Storage Key.
---

# linode\_object\_storage\_key

Provides a Linode Object Storage Key resource. This can be used to create, modify, and delete Linodes Object Storage Keys.

## Example Usage

The following example shows how one might use this resource to create an Object Storage Key.

```hcl
resource "linode_object_storage_key" "foo" {
    label = "image-access"
}

```

The following example shows a key with limited access.

```hcl
resource "linode_object_storage_key" "foobar" {
  label   = "my-key"

  bucket_access {
    bucket_name = "my-bucket-name"
    region      = "us-mia"
    permissions = "read_write"
  }
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label given to this key. For display purposes only.

* `regions` - A set of regions where the key will grant access to create buckets.

- - -

* `bucket_access` - (Optional) Defines this key as a Limited Access Key. Limited Access Keys restrict this Object Storage key’s access to only the bucket(s) declared in this array and define their bucket-level permissions. Not providing this block will not limit this Object Storage Key.

### bucket_access

The following arguments are supported in the bucket_access block:

* `bucket_name` - The unique label of the bucket to which the key will grant limited access.

* `cluster` - (Deprecated) The Object Storage cluster where the bucket resides. Deprecated in favor of `region`.

* `region` - The region where the bucket resides.

* `permissions` - This Limited Access Key’s permissions for the selected bucket. *Changing `permissions` forces the creation of a new Object Storage Key.* (`read_write`, `read_only`)

## Attributes Reference

This resource exports the following attributes:

* `access_key` - This keypair's access key. This is not secret.

* `secret_key` - This keypair's secret key.

* `limited` - Whether or not this key is a limited access key.

* `regions_details` - A set of objects containing the detailed info of the regions where this key can access.

  * `id` - The ID of the region.

  * `s3_endpoint` - The S3-compatible hostname you can use to access the Object Storage buckets in this region.
