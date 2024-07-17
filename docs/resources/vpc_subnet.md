---
page_title: "Linode: linode_vpc_subnet"
description: |-
  Manages a Linode VPC subnet.
---

# linode\_vpc\_subnet

Manages a Linode VPC subnet.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-vpc-subnet).

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

* `label` - (Required) The label of the VPC. Only contains ASCII letters, digits and dashes.

* `ipv4` - (Required) The IPv4 range of this subnet in CIDR format.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC Subnet.

* `linodes` - A list of Linode IDs that added to this subnet.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was last updated.

## Import

Linode Virtual Private Cloud (VPC) Subnet can be imported using the `vpc_id` followed by the subnet `id` separated by a comma, e.g.

```sh
terraform import linode_vpc_subnet.my_subnet_duplicated 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
