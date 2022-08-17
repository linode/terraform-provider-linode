---
layout: "linode"
page_title: "Linode: linode_instance_types"
sidebar_current: "docs-linode-datasource-instance-types"
description: |-
Provides information about Linode Instance types that match a set of filters.
---

# Data Source: linode\_instance_types

Provides information about Linode Instance types that match a set of filters.

## Example Usage

Get information about all Linode Instance types with a certain number of VCPUs:

```hcl
data "linode_instance_types" "specific-types" {
  filter {
    name = "vcpus"
    values = [2]
  }
}
```

Get information about all Linode Instance types:

```hcl
data "linode_instance_types" "all-types" {}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Instance types that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Instance type will be stored in the `types` attribute and will export the following attributes:

* `id` - The ID representing the Linode Type.

* `label` - The Linode Type's label is for display purposes only.

* `class` - The class of the Linode Type. See all classes [here](https://www.linode.com/docs/api/linode-types/#type-view__responses).

* `disk` - The Disk size, in MB, of the Linode Type.

* `price.0.hourly` -  Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `addons.0.backups.0.price.0.hourly` - The cost (in US dollars) per hour to add Backups service.

* `addons.0.backups.0.price.0.monthly` - The cost (in US dollars) per month to add Backups service.

* `network_out` - The Mbits outbound bandwidth allocation.

* `memory` - The amount of RAM included in this Linode Type.

* `transfer` - The monthly outbound transfer amount, in MB.

* `vcpus` - The number of VCPU cores this Linode Type offers.

## Filterable Fields

* `class`

* `disk`

* `gpus`

* `label`

* `memory`

* `network_out`

* `transfer`

* `vcpus`
