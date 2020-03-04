---
layout: "linode"
page_title: "Linode: linode_object_storage_cluster"
sidebar_current: "docs-linode-datasource-object-storage-cluster"
description: |-
  Provides details about a Linode Object Storage Cluster.
---

# Data Source: linode\_object\_storage\_cluster

Provides information about a Linode Object Storage Cluster

## Example Usage

The following example shows how one might use this data source to access information about a Linode Object Storage Cluster.

```hcl
data "linode_object_storage_cluster" "primary" {
    id = "us-east-1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this cluster.

## Attributes

The Linode Object Storage Cluster resource exports the following attributes:

* `domain` - The base URL for this cluster.

* `status` - This cluster's status.

* `region` - The region this cluster is located in.

* `static_site_domain` - The base URL for this cluster used when hosting static sites.
