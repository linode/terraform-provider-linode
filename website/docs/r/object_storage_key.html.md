---
layout: "linode"
page_title: "Linode: linode_object_storage_key"
sidebar_current: "docs-linode-resource-object-storage-key"
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

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label given to this key. For display purposes only.

- - -

* `bucket_access` - (Optional) Defines this key as a Limited Access Key. Limited Access Keys restrict this Object Storage key’s access to only the bucket(s) declared in this array and define their bucket-level permissions. Not providing this block will not limit this Object Storage Key.

### bucket_access

The following arguments are supported in the bucket_access block:

* `bucket_name` - The unique label of the bucket to which the key will grant limited access.

* `cluster` - The Object Storage cluster where a bucket to which the key is granting access is hosted.

* `permissions` - This Limited Access Key’s permissions for the selected bucket. *Changing `permissions` forces the creation of a new Object Storage Key.* (`read_write`, `read_only`)

## Attributes

This resource exports the following attributes:

* `access_key` - This keypair's access key. This is not secret.

* `secret_key` - This keypair's secret key.

* `limited` - Whether or not this key is a limited access key.
