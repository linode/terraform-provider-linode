---
page_title: "Linode: linode_vpc_ips"
description: |-
  Lists all ips under a Linode account or under a Linode VPC.
---

# Data Source: linode\_vpc\_ips

Provides information about a list of Linode VPC IPs that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpcs-ips).

Provides information about a list of Linode VPC IPs in a specific VPC that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-vpc-ips).

## Example Usage

The following example shows how one might use this data source to list VPC IPs.

```hcl
data "linode_vpc_ips" "filtered-ips" {
    filter {
        name = "address"
        values = ["10.0.0.0"]
    }
}

output "vpc_ips" {
  value = data.linode_vpc_ips.filtered-ips.vpc_ips
}
```

One might also use this data source to list all VPC IPs in a specific VPC. The following example shows how to do this.

```hcl
data "linode_vpc_ips" "specific-vpc-ips" {
    vpc_id = 123
}

output "vpc_ips" {
  value = data.linode_vpc_ips.specific-vpc-ips.vpc_ips
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional) The id of the parent VPC for the list of VPC IPs.

* [`filter`](#filter) - (Optional) A set of filters used to select Linode VPC IPs that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VPC IP will be stored in the `vpc_ips` attribute and will export the following attributes:

* `address` - An IPv4 address configured for this VPC interface. These follow the RFC 1918 private address format. Null if an address_range.

* `gateway` - The default gateway for the VPC subnet that the IP or IP range belongs to.

* `linode_id` - The identifier for the Linode the VPC interface currently belongs to.

* `prefix` - The number of bits set in the subnet mask.

* `region` - The region of the VPC.

* `subnet_mask` - The mask that separates host bits from network bits for the address or address_range.

* `nat_1_1` - The public IP address used for NAT 1:1 with the VPC. This is empty if NAT 1:1 isn't used.

* `subnet_id` - The id of the VPC Subnet for this interface.

* `config_id` - The globally general entity identifier for the Linode configuration profile where the VPC is included.

* `interface_id` - The globally general API entity identifier for the Linode interface.

* `address_range` - A range of IPv4 addresses configured for this VPC interface. Null if a single address.

* `vpc_id` - The unique globally general API entity identifier for the VPC.

* `active` - True if the VPC interface is in use, meaning that the Linode was powered on using the config_id to which the interface belongs. Otherwise false.

## Filterable Fields

* `active`

* `config_id`

* `linode_id`

* `region`

* `vpc_id`
