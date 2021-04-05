---
layout: "linode"
page_title: "Linode: linode_object_storage_bucket"
sidebar_current: "docs-linode-resource-object-storage-bucket"
description: |-
  Manages a Linode Object Storage Bucket.
---

# linode\_object\_storage\_bucket

Provides a Linode Object Storage Bucket resource. This can be used to create, modify, and delete Linodes Object Storage Buckets.

## Example Usage

The following example shows how one might use this resource to create an Object Storage Bucket.

```hcl
data "linode_object_storage_cluster" "primary" {
  id = "us-east-1"
}

resource "linode_object_storage_bucket" "foobar" {
  cluster = data.linode_object_storage_cluster.primary.id
  label = "%s"
}

```

## Argument Reference

The following arguments are supported:

* `cluster` - (Required) The cluster of the Linode Object Storage Bucket.

* `label` - (Required) The label of the Linode Object Storage Bucket.

* `acl` - (Optional) The Access Control Level of the bucket using a canned ACL string. See all ACL strings [in the Linode API v4 documentation](linode.com/docs/api/object-storage/#object-storage-bucket-access-update__request-body-schema).

* `cors_enabled` - (Optional) If true, the bucket will have CORS enabled for all origins.

* `versioning` - (Optional) Whether to enable versioning. Once you version-enable a bucket, it can never return to an unversioned state. You can, however, suspend versioning on that bucket.

* [`lifecycle_rule`](#lifecycle_rule) - (Optional) Lifecycle rules to be applied to the bucket.

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
