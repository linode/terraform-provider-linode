---
layout: "linode"
page_title: "Linode: linode_vpc"
sidebar_current: "docs-linode-datasource-vpc"
description: |-
  Provides details about a Linode VPC.
---

# Data Source: linode\_vpc

Provides information about a Linode VPC.

## Example Usage

The following example shows how one might use this data source to access information about a Linode VPC.

```hcl
data "linode_vpc" "foo" {
    id = 123
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique id of this VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the VPC.

* `description` - The user-defined description of this VPC.

* `region` - The region of the VPC.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was updated.

### Subnets Reference

To list all subnets under a VPC, please refer to the [linode_vpc_subnets](vpc_subnets.html.markdown) data source. 