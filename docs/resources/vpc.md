---
page_title: "Linode: linode_vpc"
description: |-
  Manages a Linode VPC.
---

# linode\_vpc

Manages a Linode VPC.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-vpc).

Please refer to [linode_vpc_subnet](vpc_subnet.html.markdown) to manage the subnets under a Linode VPC.

## Example Usage

Create a VPC:

```terraform
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    description = "My first VPC."
}
```

Create a VPC with a `/52` IPv6 range prefix:

```terraform
# NOTE: IPv6 VPCs may not currently be available to all users.
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    
    ipv6 = [
      {
        range = "/52"
      }
    ]
}
```

Create a VPC with a custom IPv4 range:

```terraform
# NOTE: Custom VPC IPv4 Ranges may not currently be available to all users.
resource "linode_vpc" "test" {
    label = "test-vpc"
    region = "us-iad"
    
    ipv4 = [
      {
        range = "10.0.0.0/8"
      }
    ]
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the VPC. This field can only contain ASCII letters, digits and dashes.

* `region` - (Required) The region of the VPC.

* `description` - (Optional) The user-defined description of this VPC.

* [`ipv6`](#ipv6) - (Optional) A list of IPv6 allocations under this VPC.

* [`ipv4`](#ipv4) - (Optional) A list of IPv4 ranges under this VPC.

## IPv6

-> **Limited Availability** IPv6 VPCs may not currently be available to all users.

Configures a single IPv6 range under this VPC.

* `range` - (Optional) An existing IPv6 prefix owned by the current account or a forward slash (/) followed by a valid prefix length. If unspecified, a range with the default prefix will be allocated for this VPC.

* `allocation_class` - (Optional) Indicates the labeled IPv6 Inventory that the VPC Prefix should be allocated from.

* `allocated_range` - (Read-Only) The value of range computed by the API. This is necessary when needing to access the range for an implicit allocation.

## IPv4

-> **Limited Availability** Custom VPC IPv4 Ranges may not currently be available to all users.

Configures a single IPv4 range under this VPC. Unlike IPv6, IPv4 ranges can be updated in-place without requiring resource replacement.

* `range` - (Required) The IPv4 range in CIDR format to assign to this VPC (e.g. `10.0.0.0/8`).

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `id` - The ID of the VPC.

* `created` - The date and time when the VPC was created.

* `updated` - The date and time when the VPC was last updated.
