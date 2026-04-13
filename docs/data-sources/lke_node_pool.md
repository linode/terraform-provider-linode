---
page_title: "Linode: linode_lke_node_pool"
description: |-
  Provides details about a specific LKE Cluster Node Pool.
---

# Data Source: linode\_lke\_node\_pool

Provides details about a specific LKE Cluster Node Pool.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-lke-node-pool).

## Example Usage

```terraform
data "linode_lke_node_pool" "my-node-pool" {
    id         = 123
    cluster_id = 321
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The LKE Cluster's Node Pool ID.

* `cluster_id` - (Required) The LKE Cluster's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `autoscaler` - When enabled, the number of nodes autoscales within the defined minimum and maximum values.

  * `enabled` - Whether autoscaling is enabled for this node pool.

  * `max` - The maximum number of nodes to autoscale to.

  * `min` - The minimum number of nodes to autoscale to.

* `disk_encryption` - Indicates the local disk encryption setting for this LKE node pool.

* `isolation` - Network isolation settings for this node pool.

  * `public_ipv4` - Whether nodes have public IPv4 addresses.

  * `public_ipv6` - Whether nodes have public IPv6 addresses.

* `disks` - This node pool's custom disk layout.

  * `size` - The size of this custom disk partition in MB.

  * `type` - This custom disk partition's filesystem type.

* `firewall_id` - The ID of the Cloud Firewall assigned to this node pool. This field is available as part of the beta API and can only be used by accounts with access to LKE Enterprise.

* `k8s_version` - The Kubernetes version used for the worker nodes within this node pool. This field is available as part of the beta API and can only be used by accounts with access to LKE Enterprise.

* `label` - The optional label defined for this node pool.

* `labels` - Key-value pairs added as labels to nodes in the node pool.

* `node_count` - The number of nodes in the node pool.

* `nodes` - Status information for the nodes that are members of this node pool.

  * `id` - The Node's ID.

  * `instance_id` - The Linode's ID. When no Linode is currently provisioned for this node, this is null.

  * `status` - The creation status of this node.

* `tags` - An array of tags applied to this object.

* `taints` - Kubernetes taints to add to node pool nodes.

  * `effect` - The Kubernetes taint effect.

  * `key` - The Kubernetes taint key.

  * `value` - The Kubernetes taint value.

* `type` - The Linode type for all of the nodes in the node pool.

* `update_strategy` - Determines when the worker nodes within this node pool upgrade to the latest selected Kubernetes version. This field is available as part of the beta API.
