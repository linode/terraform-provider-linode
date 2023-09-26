---
layout: "linode"
page_title: "Linode: linode_vpc_subnets"
sidebar_current: "docs-linode-datasource-vpc-subnets"
description: |-
  Lists all subnets under a Linode VPC.
---

# Data Source: linode\_vpc\_subnets

Provides information about a list of Linode VPC subnets that match a set of filters.

## Example Usage

The following example shows how one might use this data source to list VPC subnets.

```hcl
data "linode_vpc_subnets" "filtered-subnets" {
    vpc_id = 123
    filter {
        name = "label"
        values = ["test"]
    }
}

output "vpc_subnets" {
  value = data.linode_vpc_subnets.filtered-subnets.vpc_subnets
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the parent VPC for the list of VPCs.

* [`filter`](#filter) - (Optional) A set of filters used to select Linode users that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VPC subnet will be stored in the `vpc_subnets` attribute and will export the following attributes:

* `id` - The unique id of the VPC subnet.

* `label` - The label of the VPC subnet.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* `linodes` - A list of Linode IDs that added to this subnet.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was updated.

## Filterable Fields

* `id`

* `label`

* `ipv4`