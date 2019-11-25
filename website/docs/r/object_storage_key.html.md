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

## Attributes

This resource exports the following attributes:

* `access_key` - This keypair's access key. This is not secret.

* `secret_key` - This keypair's secret key.
