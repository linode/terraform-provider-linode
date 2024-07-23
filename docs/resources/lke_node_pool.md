---
page_title: "Linode: linode_lke_node_pool"
description: |-
  Manages an LKE Node Pool.
---

# linode\_lke\_node\_pool

Manages an LKE Node Pool.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-lke-cluster-pools).

~> **Notice** To prevent LKE node pools managed by this resource from being
recreated by the linode_lke_cluster resource, the cluster's external_pool_tags
 attribute must match the tags attribute of this resource. Please review the
[Externally Managed Node Pools](lke_cluster.md#externally-managed-node-pools)
section for more information.

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

Creating an LKE Node Pool for a Terraform-managed LKE cluster:

```terraform
locals {
  external_pool_tag = "external"
}

resource "linode_lke_node_pool" "my-pool" {
    cluster_id  = linode_lke_cluster.my-cluster.id
    type        = "g6-standard-2"
    node_count  = 3
  
    tags = [local.external_pool_tag]
}

resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.28"
    region      = "us-mia"
    
    # This tells the Linode provider to ignore 
    # node pools with the tag `external`, preventing
    # externally managed node pools from being deleted.
    external_pool_tags = [local.external_pool_tag]

  # Due to certain restrictions in Terraform and LKE, 
  # the cluster must be defined with at least one node pool.
    pool {
        type  = "g6-standard-1"
        count = 1
    }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - ID of the LKE Cluster where to create the current Node Pool.

* `type` - (Required) A Linode Type for all nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

* `node_count` - (Required; Optional with `autoscaler`) The number of nodes in the Node Pool. If undefined with an autoscaler the initial node count will equal the autoscaler minimum.

* `tags` - (Optional) An array of tags applied to the Node Pool. Tags can be used to flag node pools as externally managed, see [Externally Managed Node Pools](lke_cluster.md#externally-managed-node-pools) for more details.

* [`autoscaler`](#autoscaler) - (Optional) If defined, an autoscaler will be enabled with the given configuration.

### autoscaler

The following arguments are supported in the `autoscaler` specification block:

* `min` - (Required) The minimum number of nodes to autoscale to.

* `max` - (Required) The maximum number of nodes to autoscale to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Node Pool within LKE Cluster.

* `disk_encryption` - The disk encryption policy for nodes in this pool.

  * **NOTE: Disk encryption may not currently be available to all users.**

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
