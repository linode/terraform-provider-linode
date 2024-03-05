---
page_title: "Linode: linode_lke_cluster"
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
    k8s_version = "1.28"
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
    k8s_version = "1.28"
    region      = "us-central"
    tags        = ["prod"]

    pool {
        # NOTE: If count is undefined, the initial node count will
        # equal the minimum autoscaler node count.
        type  = "g6-standard-2"

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

* `tags` - (Optional) An array of tags applied to the Kubernetes cluster. Tags are case-insensitive and are for organizational purposes only.

* `external_pool_tags` - (Optional) An array of tags indicating that node pools having those tags are defined with a separate `linode_lke_node_pool` resource, rather than inside the current cluster resource.

### pool

~> **Notice** Due to limitations in Terraform, the order of pools in the `linode_lke_cluster` resource is treated as significant.
For example, the removal of the first listed pool in a cluster may result in all other node pools
being updated accordingly. See the [Nested Node Pool Caveats](#nested-node-pool-caveats) section for more details.

The following arguments are supported in the `pool` specification block:

* `type` - (Required) A Linode Type for all of the nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

* `count` - (Required; Optional with `autoscaler`) The number of nodes in the Node Pool. If undefined with an autoscaler the initial node count will equal the autoscaler minimum.

* [`autoscaler`](#autoscaler) - (Optional) If defined, an autoscaler will be enabled with the given configuration.

### autoscaler

The following arguments are supported in the `autoscaler` specification block:

* `min` - (Required) The minimum number of nodes to autoscale to.

* `max` - (Required) The maximum number of nodes to autoscale to.

### control_plane

The following arguments are supported in the `control_plane` specification block:

* `high_availability` - (Optional) Defines whether High Availability is enabled for the cluster Control Plane. This is an **irreversible** change.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster.

* `status` - The status of the cluster.

* `api_endpoints` - The endpoints for the Kubernetes API server.

* `kubeconfig` - The base64 encoded kubeconfig for the Kubernetes cluster.

* `dashboard_url` - The Kubernetes Dashboard access URL for this cluster.

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

## Nested Node Pool Caveats

Due to limitations in Terraform, there are some minor caveats that may cause unexpected behavior when updating
nested `pool` blocks in this resource.
Primarily, the order of `pool` blocks is significant because the ID of each pool is resolved from
the Terraform state.

For example, updating the following configuration:

```terraform
resource "linode_lke_cluster" "my-cluster" {
  # ...
  
  pool {
    type  = "g6-standard-1"
    count = 2
  }

  pool {
    type  = "g6-standard-2"
    count = 3
  }
}
```

to this:

```terraform
resource "linode_lke_cluster" "my-cluster" {
  # ...

  pool {
    type  = "g6-standard-2"
    count = 3
  }
}
```

will produce the following plan:

```terraform
~ resource "linode_lke_cluster" "my-cluster" {
      ~ pool {
            id = ... -> null
          ~ count = 2 -> 3
          ~ type  = "g6-standard-1" -> "g6-standard-2"
        }
      - pool {
          - count = 3 -> null
          - id    = ... -> null
          - nodes = [
              ...
            ] -> null
        }
  }
```

In this case, the first node pool from the original configuration will be updated to match
the second node pool's configuration.

Although not ideal, this functionality guarantees that updates to nested node pools will be reliable and predictable.
