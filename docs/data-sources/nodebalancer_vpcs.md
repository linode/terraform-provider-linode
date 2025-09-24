---
page_title: "Linode: linode_nodebalancer_vpcs"
description: |-
  Provides information about NodeBalancers VPC configurations that match a set of filters.
  For more information, see corresponding [Linode APIv4 documentation](https://techdocs.akamai.com/linode-api/reference/get-node-balancer-vpcs).
---

# linode_nodebalancer_vpcs

-> **Limited Availability** VPC-attached NodeBalancers may not currently be available to all users.

Provides information about Linode NodeBalancers VPC configurations that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-node-balancers).

## Example Usage

Retrieve all VPC configurations under a NodeBalancer:

```hcl
data "linode_nodebalancer_vpcs" "vpc-configs" {
  nodebalancer_id = 12345
}
```

Retrieve all VPC configurations under a NodeBalancer with an IPv4 range of "10.0.0.4/30":

```hcl
data "linode_nodebalancer_vpcs" "vpc-configs" {
  nodebalancer_id = 12345
  
  filter {
    name = "ipv4_range"
    values = ["10.0.0.4/30"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `nodebalancer_id` - (Required) The ID of the NodeBalancer to list VPC configurations for.

* [`filter`](#filter) - (Optional) A set of filters used to select VPC configurations that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each VPC configuration will be stored in the `vpc_configs` attribute and will export the following attributes:

* `nodebalancer_id` - The ID of the parent NodeBalancer for this VPC configuration.

* `id` - The ID of the VPC configuration.

* `ipv4_range` - A CIDR range for the VPC's IPv4 addresses. The NodeBalancer sources IP addresses from this range when routing traffic to the backend VPC nodes.

* `subnet_id` - The ID of this configuration's VPC subnet.

* `vpc_id` - The ID of this configuration's VPC.

## Filterable Fields

* `id`

* `ipv4_range`

* `nodebalancer_id`

* `subnet_id`

* `vpc_id`
