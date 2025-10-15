---
page_title: "Linode: linode_vpc_subnets"
description: |-
  Lists all subnets under a Linode VPC.
---

# Data Source: linode\_vpc\_subnets

Provides information about a list of Linode VPC subnets that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpc-subnets).

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

* [`filter`](#filter) - (Optional) A set of filters used to select Linode VPC subnets that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VPC subnet will be stored in the `vpc_subnets` attribute and will export the following attributes:

* `id` - The unique id of the VPC subnet.

* `label` - The label of the VPC subnet.

* `ipv4` - The IPv4 range of this subnet in CIDR format.

* `linodes` - A list of Linodes added to this subnet.
  * `id` - ID of the Linode
  * `interfaces` - A list of networking interfaces objects.
    * `id` - ID of the interface.
    * `config_id` - ID of Linode Config that the interface is associated with. `null` for a Linode Interface.
    * `active` - Whether the Interface is actively in use.

* [`ipv6`](#ipv6) - A list of IPv6 ranges under this subnet.

* `created` - The date and time when the VPC Subnet was created.

* `updated` - The date and time when the VPC Subnet was last updated.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

The following attributes are exported under each entry of the `ipv6` field:

* `range` - An IPv6 range allocated to this subnet in CIDR format.

## Filterable Fields

* `id`

* `label`

* `ipv4`
