---
layout: "linode"
page_title: "Linode: linode_kernel"
sidebar_current: "docs-linode-datasource-kernel"
description: |-
  Provides details about Linode Kernels that match a set of filters.
---

# Data Source: linode\_kernels

Provides information about Linode Kernels that match a set of filters.

## Example Usage

The following example shows how one might use this data source to access information about a Linode Kernel.

```hcl
data "linode_kernels" "filtered_kernels" {
    filter {
        name = "label"
        values = ["my-kernel"]
    }

    filter {
        name = "architecture"
        values = ["x86_64"]
    }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Kernels that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Kernel will be stored in the `kernel` attribute and will export the following attributes:

* `id` - The unique ID of this Kernel.

* `architecture` - The architecture of this Kernel.

* `deprecated` - Whether or not this Kernel is deprecated.

* `kvm` - If this Kernel is suitable for KVM Linodes.

* `label` - The friendly name of this Kernel.

* `pvops` - If this Kernel is suitable for paravirtualized operations.

* `version` - Linux Kernel version

* `xen` - If this Kernel is suitable for Xen Linodes.

## Filterable Fields

* `id`

* `architecture`

* `deprecated`

* `kvm`

* `label`

* `pvops`

* `version`

* `xen`