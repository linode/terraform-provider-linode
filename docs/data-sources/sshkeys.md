---
page_title: "Linode: linode_sshkeys"
description: |-
    Provides details about Linode SSH Keys that match a set of filters.
---

# Data Source: linode\_sshkeys

`linode_sshkey` provides access to a filtered list of SSH Keys in the Profile of the User identified by the access token.

## Example Usage

The following example shows how the resource might be used to obtain the names of the SSH Keys configured on the Linode user profile.

The following example shows how one might use this data source to access information about a Linode Kernel.

```hcl
data "linode_sshkeys" "filtered_ssh" {
    filter {
        name = "label"
        values = ["my-ssh"]
    }
    filter {
        name = "ssh_key"
        values = ["RSA-6522525"]
    }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode SSH Keys that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode SSH Key will be stored in the `sshkeys` attribute and will export the following attributes:

* `id` - The ID of the SSH Key.

* `label` - The label of the SSH Key.

* `ssh_key` - The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.

* `created` - The date this key was added.

## Filterable Fields

* `id`

* `label`

* `ssh_key`
