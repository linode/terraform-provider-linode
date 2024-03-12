---
page_title: "Linode: linode_lke_node_pool"
description: |-
  Manages a Linode Node Pool.
---

# linode\_nodepool

Manages an LKE Node Pool.

## Example Usage

Creating a basic LKE Node Pool:

```terraform
resource "linode_lke_node_pool" "my-pool" {
  
    cluster_id  = 150003
    type  = "g6-standard-2"
    node_count = 3
  
}
```

Creating an LKE Node Pool with autoscaler:

```terraform
resource "linode_lke_node_pool" "my-pool" {

    cluster_id  = 150003
    type  = "g6-standard-2"
  
    autoscaler {
      min = 3
      max = 10
    }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - ID of the LKE Cluster where to create the current Node Pool.

* `type` - (Required) A Linode Type for all of the nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

* `node_count` - (Required; Optional with `autoscaler`) The number of nodes in the Node Pool. If undefined with an autoscaler the initial node count will equal the autoscaler minimum.

* `tags` - (Optional) An array of tags applied to the Node Pool. Tags are for organizational purposes only.

* [`autoscaler`](#autoscaler) - (Optional) If defined, an autoscaler will be enabled with the given configuration.

### autoscaler

The following arguments are supported in the `autoscaler` specification block:

* `min` - (Required) The minimum number of nodes to autoscale to.

* `max` - (Required) The maximum number of nodes to autoscale to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Node Pool within LKE Cluster.

* [`nodes`](#nodes) - The nodes in the Node Pool.

### nodes

The following attributes are available on nodes:

* `id` - The ID of the node.

* `instance_id` - The ID of the underlying Linode instance.

* `status` - The status of the node. (`ready`, `not_ready`)

## Import

LKE Node Pools can be imported using the `cluster_id,id`, e.g.

```sh
terraform import linode_lke_node_pool.my_pool 150003,12345
```
