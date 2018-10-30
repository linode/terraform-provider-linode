---
layout: "linode"
page_title: "Linode: linode_region"
sidebar_current: "docs-linode-datasource-region"
description: |-
  Provides details about a specific service region
---

# Data Source: linode\_region

`linode_region` provides details about a specific Linode region.

## Example Usage

The following example shows how the resource might be used to obtain additional information about a Linode region.

```hcl
data "linode_region" "region" {
  id = "us-east"
}
```

## Argument Reference

- `id` - (Required) The code name of the region to select.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `country` - The country the region resides in.
