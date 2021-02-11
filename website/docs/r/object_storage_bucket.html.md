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

* [`cert`](#cert) - (Optional) The bucket's TLS/SSL certificate.

### cert

The following arguments are supported in the cert specification block:

* `certificate` - (Required) The Base64 encoded and PEM formatted SSL certificate.

* `private_key` - (Required) The private key associated with the TLS/SSL certificate.

## Import

Linodes Object Storage Buckets can be imported using the resource `id` which is made of `cluster:label`, e.g.

```sh
terraform import linode_object_storage_bucket.mybucket us-east-1:foobar
```
