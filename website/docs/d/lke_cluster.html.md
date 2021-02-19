---
layout: "linode"
page_title: "Linode: linode_lke_cluster"
sidebar_current: "docs-linode-datasource-resource-lke-cluster"
description: |-
  Provides details about an LKE Cluster.
---

# Data Source: linode\_lke_cluster

Provides details about an LKE Cluster.

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

* `tags` - The tags applied to the cluster.

* `status` - The status of the cluster.

* `api_endpoints` - The endpoints for the Kubernetes API server.

* `kubeconfig` - The base64 encoded kubeconfig for the Kubernetes cluster.

* `pools` - Node pools associated with this cluster.

  * `id` - The ID of the Node Pool.

  * `type` - The linode type for all of the nodes in the Node Pool.

  * `count` - The number of nodes in the Node Pool.

  * `nodes` - The nodes in the Node Pool.

    * `id` - The ID of the node.

    * `instance_id` - The ID of the underlying Linode instance.

    * `status` - The status of the node.
