---
layout: "linode"
page_title: "Linode: linode_vpc_subnet"
sidebar_current: "docs-linode-datasource-vpc-subnet"
description: |-
  Provides details about a subnet under a Linode VPC.
---

# Data Source: linode\_vpc\_subnet

Provides information about a Linode VPC subnet.

## Example Usage

The following example shows how one might use this data source to access information about a Linode VPC subnet.

```hcl
data "linode_vpc_subnet" "foo" {
    vpc_id = 123
    id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the parent VPC for this VPC Subnet.

* `id` - (Required) The unique id of this VPC subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the VPC subnet.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* `linodes` - A list of Linode IDs that added to this subnet.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was updated.