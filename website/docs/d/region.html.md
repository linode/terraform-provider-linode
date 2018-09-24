---
layout: "linode"
page_title: "Linode: linode_region"
sidebar_current: "docs-linode-datasource-region"
description: |-
    Provides details about a specific service region
---

# Data Source: linode_region

`linode_region` provides details about a specific Linode region.

As well as validating a given region name this resource can be used to
discover the name of the region configured within the provider. The latter
can be useful in a child module which is inheriting an AWS provider
configuration from its parent module.

## Example Usage

The following example shows how the resource might be used to obtain
the name of the AWS region configured on the provider.

```hcl
data "linode_region" "region" {
  id = "us-east"
}
```

## Argument Reference

* `id` - (Required) The code name of the region to select.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `country` - The country the region resides in.
