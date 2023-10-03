---
page_title: "Linode: linode_vlans"
description: |-
  Provides details about Linode VLANs.
---

# Data Source: linode\_vlans

~> **Beta Notice** VLANs are currently available through early access.
To use early access resources, the `api_version` provider argument must be set to `v4beta`.
To learn more, see the [early access documentation](../..#early-access).

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

output "vlan_linodes" {
  value = data.linode_vlans.my-vlans.vlans.0.linodes
}
```

## Argument Reference

The following arguments are supported

* [`filter`](#filter) - (Optional) A set of filters used to select Linode VLANs that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode VLAN will be stored in the `vlans` attribute and will export the following attributes:

* `label` - The unique label of the VLAN.

* `linodes` - The running Linodes currently attached to the VLAN.

* `region` - The region the VLAN is located in. See all regions [here](https://api.linode.com/v4/regions).

* `created` - When the VLAN was created.

## Filterable Fields

* `label`

* `region`
