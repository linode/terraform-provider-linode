---
page_title: "Linode: linode_object_storage_bucket"
description: |-
  Provides details about a Linode Object Storage Bucket.
---

# Data Source: linode_object_storage_bucket

Provides information about a Linode Object Storage Bucket
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-bucket).

## Example Usage

The following example shows how one might use this data source to access information about a Linode Object Storage Bucket.

```hcl
data "linode_object_storage_bucket" "my-bucket" {
    label  = "my-bucket"
    region = "us-mia"
}
```

## Argument Reference

* `label` - (Required) The name of this bucket.

* `region` - The ID of the region this bucket is in. Required if `cluster` is not configured.

* `cluster` - (Deprecated) The ID of the Object Storage Cluster this bucket is in. Required if `region` is not configured.

### Attributes Reference

* `created` - When this bucket was created.

* `hostname` - The hostname where this bucket can be accessed.This hostname can be accessed through a browser if the bucket is made public.

* `id` - The id of this bucket.

* `objects` - The number of objects stored in this bucket.

* `size` - The size of the bucket in bytes.
