---
page_title: "Linode: linode_lke_versions"
description: |-
  Provides details about the Kubernetes versions available for deployment to a Kubernetes cluster.
---

# linode\_lke\_versions

Provides details about the Kubernetes versions available for deployment to a Kubernetes cluster.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-lke-versions).

## Example Usage

The following example shows how one might use this data source to access information about a Linode LKE Version.

```hcl
data "linode_lke_versions" "example" {}
```

The following example shows how one might use this data source to access information about a Linode LKE Version
with additional information about the Linode LKE Version's tier (`enterprise` or `standard`).

> **_NOTE:_**  This functionality may not be currently available to all users and can only be used with v4beta.

```hcl
data "linode_lke_versions" "example" {tier = "enterprise"}
```

## Argument Reference

The following arguments are supported:

* `tier` - (Optional) The tier (`standard` or `enterprise`) of Linode LKE Versions to fetch.

## Attributes Reference

Each Linode LKE Version will be stored in the `versions` attribute and will export the following attributes:

* `id` - The Kubernetes version numbers available for deployment to a Kubernetes cluster in the format of [major].[minor], and the latest supported patch version.

* `tier` - The Kubernetes version tier. Only exported if `tier` was provided when using the datasource.
