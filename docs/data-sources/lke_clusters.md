---
page_title: "Linode: linode_lke_clusters"
description: |-
  Provides information about LKE Clusters that match a set of filters.
---

# Data Source: linode\_lke\_clusters

Provides information about a list of current Linode Kubernetes (LKE) clusters on your account that match a set of filters.

## Example Usage

Get information about all LKE clusters with a specific tag:

```hcl
data "linode_lke_clusters" "specific" {
  filter {
    name = "tags"
    values = ["test-tag"]
  }
}

output "lke_cluster" {
  value = data.linode_lke_clusters.specific.lke_clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select LKE Clusters that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each LKE Cluster will be stored in the `lke_clusters` attribute and will export the following attributes:

* `id` - The LKE Cluster's ID.

* `created` - When this Kubernetes cluster was created.

* `updated` - When this Kubernetes cluster was updated.

* `label` - The unique label for the cluster.

* `k8s_version` - The Kubernetes version for this Kubernetes cluster in the format of `major.minor` (e.g. `1.17`).

* `tags` - An array of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

* `status` - The status of the cluster.

* `region` - This Kubernetes cluster's location.

* `control_plane.high_availability` - Whether High Availability is enabled for the cluster Control Plane.

To get more information about a cluster, i.e. node pools, please refer to the [linode_lke_cluster](lke_cluster.html.markdown) data source.

## Filterable Fields

* `k8s_version`

* `label`

* `region`

* `tags`

* `status`

* `created`

* `updated`
