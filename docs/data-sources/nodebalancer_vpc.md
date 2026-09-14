---
page_title: "Linode: linode_nodebalancer_vpc"
description: |-
  Provides information about NodeBalancer VPC configuration.
---

# linode_nodebalancer_vpc

-> **Limited Availability** VPC-attached NodeBalancers may not currently be available to all users and may require the `api_version` provider argument must be set to `v4beta`.

Provides information about a NodeBalancer VPC configuration.
For more information, see the corresponding [Linode APIv4 documentation](https://techdocs.akamai.com/linode-api/reference/get-node-balancer-vpc-config).

## Example Usage

Retrieve information about a NodeBalancer VPC configuration:

```hcl
data "linode_nodebalancer_vpc" "vpc-config" {
  nodebalancer_id = 123
  id = 456
}
```

## Arguments Reference

This data source accepts the following arguments:

* `nodebalancer_id` - (Required) The ID of the parent NodeBalancer of the VPC configuration.

* `id` - (Required) The ID of the VPC configuration.

## Attributes Reference

This data source exports the following attributes:

* `ipv4_range` - A CIDR range for the VPC's IPv4 addresses, used for backend node routing or frontend IPs depending on `purpose`.

* `ipv6_range` - A CIDR range for the VPC's IPv6 addresses, used for backend node routing or frontend IPs depending on `purpose`.

* `subnet_id` - The ID of this configuration's VPC subnet.

* `vpc_id` - The ID of this configuration's VPC.

* `purpose` - Indicates whether the VPC configuration applies to backend nodes that serve requests or to the NodeBalancer frontend which manages incoming traffic.
