---
layout: "linode"
page_title: "Linode: linode_lke_cluster"
sidebar_current: "docs-linode-resource-lke-cluster"
description: |-
  Manages a Linode instance.
---

# linode\_lke_cluster

Manages an LKE cluster.

## Example Usage

Creating a basic LKE cluster:

```terraform
resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.21"
    region      = "us-central"
    tags        = ["prod"]

    pool {
        type  = "g6-standard-2"
        count = 3
    }
}
```

Creating an LKE cluster with autoscaler:

```terraform
resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.21"
    region      = "us-central"
    tags        = ["prod"]

    pool {
        type  = "g6-standard-2"
        count = 3

        autoscaler {
          min = 3
          max = 10
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) This Kubernetes cluster's unique label.

* `k8s_version` - (Required) The desired Kubernetes version for this Kubernetes cluster in the format of `major.minor` (e.g. `1.21`), and the latest supported patch version will be deployed.

* `region` - (Required) This Kubernetes cluster's location.

* [`pool`](#pool) - (Required) The Node Pool specifications for the Kubernetes cluster. At least one Node Pool is required.

* [`control_plane`](#control_plane) (Optional) Defines settings for the Kubernetes Control Plane.

* `tags` - (Optional) An array of tags applied to the Kubernetes cluster. Tags are for organizational purposes only.

### pool

The following arguments are supported in the `pool` specification block:

* `type` - (Required) A Linode Type for all of the nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

* `count` - (Required) The number of nodes in the Node Pool.

* [`autoscaler`](#autoscaler) - (Optional) If defined, an autoscaler will be enabled with the given configuration.

### autoscaler

The following arguments are supported in the `autoscaler` specification block:

* `min` - (Required) The minimum number of nodes to autoscale to.

* `max` - (Required) The maximum number of nodes to autoscale to.

### control_plane

The following arguments are supported in the `control_plane` specification block:

* `high_availability` - (Optional) Defines whether High Availability is enabled for the cluster Control Plane. This is an **irreversible** change. **NOTICE:** High Availability Control Planes are currently available through early access. To learn more, see the [early access documentation](https://github.com/linode/terraform-provider-linode/tree/master/EARLY_ACCESS.md).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster.

* `status` - The status of the cluster.

* `api_endpoints` - The endpoints for the Kubernetes API server.

* `kubeconfig` - The base64 encoded kubeconfig for the Kubernetes cluster.

* `pool` - Additional nested attributes:

  * `id` - The ID of the Node Pool.

  * [`nodes`](#nodes) - The nodes in the Node Pool.

### nodes

The following attributes are available on nodes:

* `id` - The ID of the node.

* `instance_id` - The ID of the underlying Linode instance.

* `status` - The status of the node. (`ready`, `not_ready`)

## Import

LKE Clusters can be imported using the `id`, e.g.

```sh
terraform import linode_lke_cluster.my_cluster 12345
```
