---
page_title: "Linode: linode_object_storage_bucket"
description: |-
  Manages a Linode Object Storage Bucket.
---

# linode\_object\_storage\_bucket

Provides a Linode Object Storage Bucket resource. This can be used to create, modify, and delete Linodes Object Storage Buckets.

## Example Usage

The following example shows how one might use this resource to create an Object Storage Bucket:

```hcl
data "linode_object_storage_cluster" "primary" {
  id = "us-east-1"
}

resource "linode_object_storage_bucket" "foobar" {
  cluster = data.linode_object_storage_cluster.primary.id
  label = "mybucket"
}

```

Creating an Object Storage Bucket with Lifecycle rules:

```hcl

resource "linode_object_storage_key" "mykey" {
  label = "image-access"
}

resource "linode_object_storage_bucket" "mybucket" {
  access_key = linode_object_storage_key.mykey.access_key
  secret_key = linode_object_storage_key.mykey.secret_key

  cluster = "us-east-1"
  label   = "mybucket"

  lifecycle_rule {
    id      = "my-rule"
    enabled = true

    abort_incomplete_multipart_upload_days = 5

    expiration {
      date = "2021-06-21"
    }
  }
}
```

Creating an Object Storage Bucket with Lifecycle rules using provider-level object credentials

```hcl
provider "linode" {
    obj_access_key = ${your-access-key}
    obj_secret_key = ${your-secret-key}
}

resource "linode_object_storage_bucket" "mybucket" {
  # no need to specify the keys with the resource
  cluster = "us-east-1"
  label   = "mybucket"

  lifecycle_rule {
    # ... details of the lifecycle
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster` - (Required) The cluster of the Linode Object Storage Bucket.

* `label` - (Required) The label of the Linode Object Storage Bucket.

* `acl` - (Optional) The Access Control Level of the bucket using a canned ACL string. See all ACL strings [in the Linode API v4 documentation](https://linode.com/docs/api/object-storage/#object-storage-bucket-access-update__request-body-schema).

* `access_key` - (Optional) The access key to authenticate with. If not specified with the resource, the value of [`obj_access_key`](../index.md#configuration-reference) from provider-level will be used.

* `secret_key` - (Optional) The secret key to authenticate with. If not specified with the resource, the value of [`obj_secret_key`](../index.md#configuration-reference) from provider-level will be used.

* `cors_enabled` - (Optional) If true, the bucket will have CORS enabled for all origins.

* `versioning` - (Optional) Whether to enable versioning. Once you version-enable a bucket, it can never return to an unversioned state. You can, however, suspend versioning on that bucket. (Requires `access_key` and `secret_key`)

* [`lifecycle_rule`](#lifecycle_rule) - (Optional) Lifecycle rules to be applied to the bucket. (Requires `access_key` and `secret_key`)

* [`cert`](#cert) - (Optional) The bucket's TLS/SSL certificate.

### cert

The following arguments are supported in the cert specification block:

* `certificate` - (Required) The Base64 encoded and PEM formatted SSL certificate.

* `private_key` - (Required) The private key associated with the TLS/SSL certificate.

### lifecycle_rule

The following arguments are supported in the lifecycle_rule specification block:

* `id` - (Optional) The unique identifier for the rule.

* `prefix` - (Optional) The object key prefix identifying one or more objects to which the rule applies.

* `enabled` - (Optional) Specifies whether the lifecycle rule is active.

* `abort_incomplete_multipart_upload_days` - (Optional) Specifies the number of days after initiating a multipart upload when the multipart upload must be completed.

* [`expiration`](#expiration) - (Optional) Specifies a period in the object's expire.

* [`noncurrent_version_expiration`](#noncurrent_version_expiration) - (Optional) Specifies when non-current object versions expire.

### expiration

The following arguments are supported in the expiration specification block:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect.

* `days` - (Optional) Specifies the number of days after object creation when the specific rule action takes effect.

* `expired_object_delete_marker` - (Optional) On a versioned bucket (versioning-enabled or versioning-suspended bucket), you can add this element in the lifecycle configuration to direct Linode Object Storage to delete expired object delete markers. This cannot be specified with Days or Date in a Lifecycle Expiration Policy.

### noncurrent_version_expiration

The following arguments are supported in the noncurrent_version_expiration specification block:

* `days` - (Required) Specifies the number of days non-current object versions expire.

## Import

Linodes Object Storage Buckets can be imported using the resource `id` which is made of `cluster:label`, e.g.

```sh
terraform import linode_object_storage_bucket.mybucket us-east-1:foobar
```
