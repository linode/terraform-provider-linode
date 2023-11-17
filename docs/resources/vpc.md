---
page_title: "Linode: linode_vpc"
description: |-
  Manages a Linode VPC.
---

# linode\_vpc

Manages a Linode VPC.

Please refer to [linode_vpc_subnet](vpc_subnet.html.markdown) to manage the subnets under a Linode VPC.

## Example Usage

Create a VPC:

```terraform
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-east"
    description = "My first VPC."
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the VPC. Only contains ascii letters, digits and dashes.

* `region` - (Required) The region of the VPC.

* `description` - (Optional) The user-defined description of this VPC.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was updated.
