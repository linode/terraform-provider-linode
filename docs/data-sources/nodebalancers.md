---
page_title: "Linode: linode_nodebalancers"
description: |-
  Provides information about Linode NodeBalancers that match a set of filters.
---

# linode_nodebalancers

Provides information about Linode NodeBalancers that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-node-balancers).

## Example Usage

The following example shows how one might use this data source to access information about a Linode NodeBalancer.

```hcl
data "linode_nodebalancers" "specific-nodebalancers" {
  filter {
    name = "label"
    values = ["my-nodebalancer"]
  }

  filter {
    name = "region"
    values = ["us-iad"]
  }
}

output "nodebalancer_id" {
  value = data.linode_nodebalancers.specific-nodebalancers.nodebalancers.0.id
}
```

## Argument Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select Linode NodeBalancers that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned (`asc`, `desc`; default `asc`).

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by (`exact`, `regex`, `substring`; default `exact`).

## Attributes Reference

Each Linode NodeBalancer will be stored in the `nodebalancers` attribute and will export the following attributes:

* `nodebalancers` - (Nested Attribute List) The returned list of NodeBalancers. Referenced by index (e.g. `nodebalancers[0].id`).

* `label` - The label of the Linode NodeBalancer.

* `client_conn_throttle` - Throttle connections per second (0-20).

* `client_udp_sess_throttle` - Throttle UDP sessions per second (0-20).

  * **NOTE: This attribute may not be generally available.**

* `created` – When this Linode NodeBalancer was created.

* `tags` - A list of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

* `hostname` - This NodeBalancer's hostname, ending with .ip.linodeusercontent.com.

* `id` - The Linode NodeBalancer's unique ID.

* `ipv4` - The Public IPv4 Address of this NodeBalancer.

* `ipv6` - The Public IPv6 Address of this NodeBalancer.

* `region` - The Region where this Linode NodeBalancer is located. NodeBalancers only support backends in the same Region.

* `updated` – When this Linode NodeBalancer was last updated.

* [`transfer`](#transfer) - (Read-Only Object List) The network transfer stats for the current month. Referenced with an index (e.g. `transfer.0.in`).

* [`lke_cluster`](#lke_cluster) - (Nested Attribute List) The LKE cluster that manages this NodeBalancer, if any. The list will be empty if this NodeBalancer isn't related to an LKE cluster.

* `type` - The type of this NodeBalancer.

* `frontend_address_type` - Indicates whether incoming requests are routed to NodeBalancers using VPC frontend IPs or public frontend IPs.

* `frontend_vpc_subnet_id` - The VPC subnet assigned to this NodeBalancer.

### transfer

The following attributes are available on transfer:

* `in` - The total transfer, in MB, used by this NodeBalancer for the current month.

* `out` - The total inbound transfer, in MB, used for this NodeBalancer for the current month.

* `total` - The total outbound transfer, in MB, used for this NodeBalancer for the current month.

### lke_cluster

The following attributes are available on `lke_cluster`:

* `id` - The ID of the related LKE cluster.

* `label` - The label of the related LKE cluster.

* `type` - The type of the related LKE cluster.

* `url` - The URL where you can access the related LKE cluster.

## Filterable Fields

* `label`

* `tags`

* `ipv4`

* `ipv6`

* `hostname`

* `region`

* `client_conn_throttle`
