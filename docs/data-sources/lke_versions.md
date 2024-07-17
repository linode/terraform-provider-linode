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

## Attributes Reference

Each Linode LKE Version will be stored in the `versions` attribute and will export the following attributes:

* `id` - The Kubernetes version numbers available for deployment to a Kubernetes cluster in the format of [major].[minor], and the latest supported patch version.
