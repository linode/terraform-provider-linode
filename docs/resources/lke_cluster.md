---
page_title: "Linode: linode_lke_cluster"
description: |-
  Manages a Linode instance.
---

# linode\_lke\_cluster

Manages an LKE cluster.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-lke-cluster).

## Example Usage

Creating a basic LKE cluster:

```terraform
resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.32"
    region      = "us-central"
    tags        = ["prod"]

    pool {
        type  = "g6-standard-2"
        count = 3
    }
}
```

Creating an enterprise LKE cluster:

```terraform
resource "linode_lke_cluster" "test" {
    label       = "lke-e-cluster"
    region      = "us-lax"
    k8s_version = "v1.31.8+lke5"
    tags        = ["test"]
    tier = "enterprise"

    pool {
      type  = "g7-premium-2"
      count = 3
      tags  = ["test"]
    }
}
```

Creating an LKE cluster with autoscaler:

```terraform
resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.32"
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

Creating an LKE cluster with control plane:

```terraform
resource "linode_lke_cluster" "test" {
    label       = "my-cluster"     
    k8s_version = "1.32"           
    region      = "us-central"     
    tags        = ["prod"]         

    control_plane {
        high_availability = true
      
        acl {
            enabled = true
            addresses {
                ipv4 = ["0.0.0.0/0"]
                ipv6 = ["2001:db8::/32"]
            }
        }
    }

    pool {
        type  = "g6-standard-2"
        count = 1
    }
}
```

Creating an LKE cluster with labeled node pools:

```terraform
resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.32"
    region      = "us-central"
    tags        = ["prod"]

    pool {
        type  = "g6-standard-2"
        count = 2
        labels = {
            "role" = "database"
            "environment" = "production"
        }
    }

    pool {
        type  = "g6-standard-1"
        count = 3
        labels = {
            "role" = "application"
            "environment" = "production"
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

* `external_pool_tags` - (Optional) A set of node pool tags to ignore when planning and applying this cluster. This prevents externally managed node pools from being deleted or unintentionally updated on subsequent applies. See [Externally Managed Node Pools](#externally-managed-node-pools) for more details.

* `tier` - (Optional) The desired Kubernetes tier. (**Note: v4beta only and may not currently be available to all users.**)

* `subnet_id` - (Optional) The ID of the VPC subnet to use for the Kubernetes cluster. This subnet must be dual stack (IPv4 and IPv6 should both be enabled). (**Note: v4beta only and may not currently be available to all users.**)

* `vpc_id` - (Optional) The ID of the VPC to use for the Kubernetes cluster.

* `stack_type` - (Optional) The networking stack type of the Kubernetes cluster.

### pool

~> **Notice** Due to limitations in Terraform, the order of pools in the `linode_lke_cluster` resource is treated as significant.
For example, the removal of the first listed pool in a cluster may result in all other node pools
being updated accordingly. See the [Nested Node Pool Caveats](#nested-node-pool-caveats) section for more details.

The following arguments are supported in the `pool` specification block:

* `type` - (Required) A Linode Type for all of the nodes in the Node Pool. See all node types [here](https://api.linode.com/v4/linode/types).

* `count` - (Required; Optional with `autoscaler`) The number of nodes in the Node Pool. If undefined with an autoscaler the initial node count will equal the autoscaler minimum.

* `labels` - (Optional) A map of key/value pairs to apply to all nodes in the pool. Labels are used to identify and organize Kubernetes resources within your cluster.

* `tags` - (Optional) A set of tags applied to this node pool. Tags can be used to flag node pools as externally managed. See [Externally Managed Node Pools](#externally-managed-node-pools) for more details.

* `taint` - (Optional) Kubernetes taints to add to node pool nodes. Taints help control how pods are scheduled onto nodes, specifically allowing them to repel certain pods. See [Add Labels and Taints to your LKE Node Pools](https://www.linode.com/docs/products/compute/kubernetes/guides/deploy-and-manage-cluster-with-the-linode-api/#add-labels-and-taints-to-your-lke-node-pools).

  * `effect` - (Required) The Kubernetes taint effect. Accepted values are `NoSchedule`, `PreferNoSchedule`, and `NoExecute`. For the descriptions of these values, see [Kubernetes Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/).

  * `key` - (Required) The Kubernetes taint key.

  * `value` - (Required) The Kubernetes taint value.

* [`autoscaler`](#autoscaler) - (Optional) If defined, an autoscaler will be enabled with the given configuration.

* `k8s_version` - (Optional) The k8s version of the nodes in this Node Pool. For LKE enterprise only and may not currently available to all users even under v4beta.

* `update_strategy` - (Optional) The strategy for updating the Node Pool k8s version. For LKE enterprise only and may not currently available to all users even under v4beta.

### autoscaler

The following arguments are supported in the `autoscaler` specification block:

* `min` - (Required) The minimum number of nodes to autoscale to.

* `max` - (Required) The maximum number of nodes to autoscale to.

### control_plane

The following arguments are supported in the `control_plane` specification block:

* `high_availability` - (Optional) Defines whether High Availability is enabled for the cluster Control Plane. This is an **irreversible** change.

* `audit_logs_enabled` - (Optional) Enables audit logs on the cluster's control plane.

* [`acl`](#acl) - (Optional) Defines the ACL configuration for an LKE cluster's control plane.

### acl

The following arguments are supported in the `acl` specification block:

* `enabled` - (Optional) Defines default policy. A value of true results in a default policy of DENY. A value of false results in default policy of ALLOW, and has the same effect as delete the ACL configuration.

* [`addresses`](#addresses) - (Optional) A list of ip addresses to allow.

### addresses

The following arguments are supported in the `addresses` specification block:

* `ipv4` - (Optional) A set of individual ipv4 addresses or CIDRs to ALLOW.

* `ipv6` - (Optional) A set of individual ipv6 addresses or CIDRs to ALLOW.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster.

* `status` - The status of the cluster.

* `api_endpoints` - The endpoints for the Kubernetes API server.

* `kubeconfig` - The base64 encoded kubeconfig for the Kubernetes cluster.

* `dashboard_url` - The Kubernetes Dashboard access URL for this cluster. LKE Enterprise does not have a dashboard URL.

* `apl_enabled` - Enables the App Platform Layer

* `pool` - Additional nested attributes:

  * `id` - The ID of the Node Pool.

  * `disk_encryption` - The disk encryption policy for nodes in this pool.

    * **NOTE: Disk encryption may not currently be available to all users.**

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

## Externally Managed Node Pools

By default, the `linode_lke_cluster` resource will account for all node pools under the corresponding cluster, meaning
any node pools created externally or managed by other resources will be removed on subsequent applies.

To signal the provider to ignore externally managed node pools, the `external_pool_tags` attribute can be defined with
tags matching a tag on an externally managed node pool.

For example:

```terraform
locals {
  external_pool_tag = "external"
}

resource "linode_lke_cluster" "my-cluster" {
    label       = "my-cluster"
    k8s_version = "1.32"
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

resource "linode_lke_node_pool" "my-pool" {
  cluster_id  = linode_lke_cluster.my-cluster.id
  type        = "g6-standard-2"
  node_count  = 3

  tags = [local.external_pool_tag]
}
```
