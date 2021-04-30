---
layout: "linode"
page_title: "Linode: linode_vlans"
sidebar_current: "docs-linode-datasource-vlans"
description: |-
Provides details about Linode VLANs.
---

# Data Source: linode\_vlans

Provides details about Linode VLANs.

## Example Usage

```terraform
resource "linode_instance" "my_instance" {
  label      = "my_instance"
  image      = "linode/ubuntu18.04"
  region     = "us-southeast"
  type       = "g6-standard-1"
  root_pass  = "bogusPassword$"
  
  interface {
    purpose = "vlan"
    label = "my-vlan"
  }
}

data "linode_vlans" "my-vlans" {
  filter {
    name = "label"
    values = ["my-vlan"]
  }
}
```

## Argument Reference

The following arguments are supported

* [`filter`](#filter) - (Optional) A set of filters used to select Linode VLANs that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

## Attributes

Each Linode VLAN will be stored in the `vlans` attribute and will export the following attributes:

* `label` - The unique label of the VLAN.

* `linodes` - The running Linodes currently attached to the VLAN.

* `region` - The region the VLAN is located in.

* `created` - When the VLAN was created.

## Filterable Fields

* `label`

* `region`
