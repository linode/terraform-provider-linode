---
layout: "linode"
page_title: "Linode: linode_lke_version"
sidebar_current: "docs-linode-datasource-lke-version"
description: |-
  Provides details about the Kubernetes versions available for deployment to a Kubernetes cluster.
---

# linode\_lke\_version

Provides details about the Kubernetes versions available for deployment to a Kubernetes cluster.

## Example Usage

The following example shows how one might use this data source to access information about a Linode LKE Version.

```hcl
data "linode_lke_version" "example" {
    id = "1.25"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) A Kubernetes version number available for deployment to a Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version.

## Attributes Reference

The Linode LKE Version resource exports the following attributes:

* `id` - The Kubernetes version number available for deployment to a Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version.
