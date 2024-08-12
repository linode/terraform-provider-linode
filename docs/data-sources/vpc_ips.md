---
page_title: "Linode: linode_vpc_ips"
description: |-
  Lists all ips under a Linode account or under a Linode VPC.
---

# Data Source: linode\_vpc\_ips

Provides information about a list of Linode VPC IPs that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-ips).

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

* `address` - The IP address in CIDR format.

* `gateway` - The default gateway for this address.

* `linode_id` - The ID of the Linode this address currently belongs to. For IPv4 addresses, this defaults to the Linode that this address was assigned to on creation.

* `prefix` - The number of bits set in the subnet mask.

* `region` - The Region this IP address resides in.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `nat_1_1` - IPv4 address configured as a 1:1 NAT for this Interface. If no address is configured as a 1:1 NAT, null is returned. Only allowed for vpc type interfaces.

* `subnet_id` - The ID of the subnet this IP address currently belongs to.

* `config_id` - The ID of the config this IP address is associated with.

* `interface_id` - The ID of the interface this IP address is associated with.

* `address_range` - The IP address range that this IP address is associated with.

* `vpc_id` - The ID of the VPC this IP address is associated with.

* `active` - Indicates whether this IP address is active or not.

## Filterable Fields

* `address`

* `prefix`

* `region`
