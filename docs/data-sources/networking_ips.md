---
page_title: "Linode: linode_networking_ips"
description: |-
  Retrieves a list of IP addresses on your account, including unassigned reserved IP addresses.
---

# linode\_networking\_ips

Provides information about all IP addresses associated with the current Linode account, including both assigned and unassigned reserved IP addresses.

## Example Usage

Retrieve all IPs under the current account:

```hcl
data "linode_networking_ips" "all" {}
```

Retrieve all IPs under the current account in a specific region:

```hcl
data "linode_networking_ips" "filtered" {
  filter {
    name = "region"
    values = ["us-mia"]
  }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select IP addresses that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each IP address will be stored in the `ip_addresses` attribute and will export the following attributes:

* `address` - The IP address.

* `gateway` - The default gateway for this address.

* `subnet_mask` - The mask that separates host bits from network bits for this address.

* `prefix` - The number of bits set in the subnet mask.

* `type` - The type of address this is (ipv4, ipv6, ipv6/pool, ipv6/range).

* `public` - Whether this is a public or private IP address.

* `rdns` - The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.

* `linode_id` - The ID of the Linode this address currently belongs to.

* `interface_id` - The ID of the interface this address is assigned to.

* `region` - The Region this IP address resides in. See all regions [here](https://api.linode.com/v4/regions).

* `reserved` - Whether this IP address is a reserved IP.

* `vpc_nat_1_1` - Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.

  * `address` - The IPv4 address that is configured as a 1:1 NAT for this VPC interface.

  * `subnet_id` - The `id` of the VPC Subnet for this Interface.

  * `vpc_id` - The `id` of the VPC configured for this Interface.

## Filterable Fields

* `address`

* `gateway`

* `subnet_mask`

* `prefix`

* `type`

* `public`

* `rdns`

* `linode_id`

* `region`

* `reserved`
