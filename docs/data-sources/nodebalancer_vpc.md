---
page_title: "Linode: linode_nodebalancer_vpc"
description: |-
  Provides information about NodeBalancer VPC configuration.
---

# linode_nodebalancer_vpc

-> **Limited Availability** VPC-attached NodeBalancers may not currently be available to all users.

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

* `ipv4_range` - A CIDR range for the VPC's IPv4 addresses. The NodeBalancer sources IP addresses from this range when routing traffic to the backend VPC nodes.

* `subnet_id` - The ID of this configuration's VPC subnet.

* `vpc_id` - The ID of this configuration's VPC.
