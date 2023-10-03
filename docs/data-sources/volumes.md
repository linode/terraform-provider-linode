---
page_title: "Linode: linode_volumes"
description: |-
  Lists Volumes on your Account.
---

# linode\_volumes

Provides information about Linode volumes that match a set of filters.

```hcl
data "linode_volumes" "filtered-volumes" {
  filter {
    name = "label"
    values = ["test-volume"]
  }
}

output "volumes" {
  value = data.linode_volumes.filtered-volumes.volumes
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode volumes that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode volume will be stored in the `volumes` attribute and will export the following attributes:

* `id` - The unique ID of this Volume.

* `created` - When this Volume was created.

* `status` - The current status of the Volume. (`creating`, `active`, `resizing`, `contact_support`)

* `label` - This Volume's label is for display purposes only.

* `tags` - An array of tags applied to this object.

* `size` - The Volume's size, in GiB.

* `region` - The datacenter in which this Volume is located. See all regions [here](https://api.linode.com/v4/regions).

* `updated` - When this Volume was last updated.

* `linode_id` - If a Volume is attached to a specific Linode, the ID of that Linode will be displayed here. If the Volume is unattached, this value will be null.

* `filesystem_path` - The full filesystem path for the Volume based on the Volume's label. Path is /dev/disk/by-id/scsi-0LinodeVolume + Volume label.

## Filterable Fields

* `label`

* `tags`
