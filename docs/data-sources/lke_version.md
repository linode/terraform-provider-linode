---
page_title: "Linode: linode_lke_version"
description: |-
  Provides details about a Kubernetes version available for deployment to a Kubernetes cluster.
---

# linode\_lke\_version

Provides details about a specific Kubernetes versions available for deployment to a Kubernetes cluster.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-lke-version).

## Example Usage

The following example shows how one might use this data source to access information about a Linode LKE Version.

```hcl
data "linode_lke_version" "example" {id = "1.31"}
```

The following example shows how one might use this data source to access information about a Linode LKE Version
with additional information about the Linode LKE Version's tier (`enterprise` or `standard`).

> **_NOTE:_**  This functionality may not be currently available to all users and can only be used with v4beta.

```hcl
data "linode_lke_version" "example" {
    id = "1.31"
    tier = "standard"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this Linode LKE Version.

* `tier` - (Optional) The tier (`standard` or `enterprise`) of Linode LKE Version to fetch.

## Attributes Reference

The Linode LKE Version datasource exports the following attributes:

* `id` - The Kubernetes version numbers available for deployment to a Kubernetes cluster in the format of [major].[minor], and the latest supported patch version.

* `tier` - The Kubernetes version tier. Only exported if `tier` was provided when using the datasource.
