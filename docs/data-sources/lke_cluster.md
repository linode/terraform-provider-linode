---
page_title: "Linode: linode_lke_cluster"
description: |-
  Provides details about an LKE Cluster.
---

# Data Source: linode\_lke_cluster

Provides details about an LKE Cluster.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-lke-cluster).

## Example Usage

```terraform
data "linode_lke_cluster" "my-cluster" {
    id = 123
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The LKE Cluster's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `k8s_version` - The Kubernetes version for this Kubernetes cluster in the format of `major.minor` (e.g. `1.17`).

* `region` - This Kubernetes cluster's location.

* `tags` - The tags applied to the cluster. Tags are case-insensitive and are for organizational purposes only.

* `status` - The status of the cluster.

* `label` - The unique label for the cluster.

* `created` - When this Kubernetes cluster was created.

* `updated` - When this Kubernetes cluster was updated.

* `api_endpoints` - The endpoints for the Kubernetes API server.

* `kubeconfig` - The base64 encoded kubeconfig for the Kubernetes cluster.

* `dashboard_url` - The Kubernetes Dashboard access URL for this cluster.

* `pools` - Node pools associated with this cluster.

  * `id` - The ID of the Node Pool.

  * `type` - The linode type for all of the nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

  * `count` - The number of nodes in the Node Pool.

  * `disk_encryption` - The disk encryption policy for nodes in this pool.

    * **NOTE: Disk encryption may not currently be available to all users.**

  * `tags` - An array of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

  * `nodes` - The nodes in the Node Pool.

    * `id` - The ID of the node.

    * `instance_id` - The ID of the underlying Linode instance.

    * `status` - The status of the node. (`ready`, `not_ready`)

  * `autoscaler` - The configuration options for the autoscaler. This field only contains an autoscaler configuration if autoscaling is enabled on this cluster.

    * `enabled` - Whether autoscaling is enabled for this Node Pool. Defaults to false.

    * `min` - The minimum number of nodes to autoscale to.

    * `max` - The maximum number of nodes to autoscale to.

  * `disks` - This Node Pool’s custom disk layout.

    * `size` - The size of this custom disk partition in MB.

    * `type` - This custom disk partition’s filesystem type.

* `control_plane` - The settings for the Kubernetes Control Plane.

  * `high_availability` - Whether High Availability is enabled for the cluster Control Plane.
  
  * `acl` - The ACL configuration for an LKE cluster's control plane. **NOTE: Control Plane ACLs may not currently be available to all users.**

    * `enabled` - The default policy. A value of true means a default policy of DENY. A value of false means a default policy of ALLOW.

    * `addresses` - A list of ip addresses to allow.

      * `ipv4` - A set of individual ipv4 addresses or CIDRs to ALLOW.

      * `ipv6` - A set of individual ipv6 addresses or CIDRs to ALLOW.
