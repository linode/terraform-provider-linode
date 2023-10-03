---
page_title: "Linode: linode_nodebalancers"
description: |-
  Provides information about Linode NodeBalancers that match a set of filters.
---

# linode_nodebalancers

Provides information about Linode NodeBalancers that match a set of filters.

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

* [`filter`](#filter) - (Optional) A set of filters used to select Linode NodeBalancers that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode NodeBalancer will be stored in the `nodebalancers` attribute and will export the following attributes:

* `label` - The label of the Linode NodeBalancer

* `client_conn_throttle` - Throttle connections per second (0-20)

* `created` – When this Linode NodeBalancer was created

* `linode_id` - The ID of a Linode Instance where the NodeBalancer should be attached

* `tags` - A list of tags applied to this object. Tags are for organizational purposes only.

* `hostname` - This NodeBalancer's hostname, ending with .ip.linodeusercontent.com

* `id` - The Linode NodeBalancer's unique ID

* `ipv4` - The Public IPv4 Address of this NodeBalancer

* `ipv6` - The Public IPv6 Address of this NodeBalancer

* `region` - The Region where this Linode NodeBalancer is located. NodeBalancers only support backends in the same Region.

* `updated` – When this Linode NodeBalancer was last updated

* [`transfer`](#transfer) - The network transfer stats for the current month

### transfer

The following attributes are available on transfer:

* `in` - The total transfer, in MB, used by this NodeBalancer for the current month

* `out` - The total inbound transfer, in MB, used for this NodeBalancer for the current month

* `total` - The total outbound transfer, in MB, used for this NodeBalancer for the current month

## Filterable Fields

* `label`

* `tags`

* `ipv4`

* `ipv6`

* `hostname`

* `region`

* `client_conn_throttle`
