---
layout: "linode"
page_title: "Linode: linode_vpc_subnet"
sidebar_current: "docs-linode-resource-vpc-subnet"
description: |-
  Manages a Linode VPC subnet.
---

# linode\_vpc\_subnet

Manages a Linode VPC subnet.

## Example Usage

Create a VPC subnet:

```terraform
resource "linode_vpc_subnet" "test" {
    vpc_id = 123
    label = "test-subnet"
    ipv4 = "10.0.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the parent VPC for this VPC Subnet.

* `label` - (Required) The label of the VPC. Only contains ascii letters, digits and dashes.

* `ipv4` - (Required) The IPv4 range of this subnet in CIDR format.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC Subnet.

* `linodes` - A list of Linode IDs that added to this subnet.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was updated.
