---
page_title: "Linode: linode_volume"
description: |-
  Provides details about a Linode Volume.
---

# Data Source: linode_volume

Provides information about a Linode Volume.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-volume).

## Example Usage

The following example shows how one might use this data source to access information about a Linode Volume.

```hcl
data "linode_volume" "foo" {
    id = "1234567"
}
```

## Argument Reference

The following argument is required:

- `id` - The unique numeric ID of the Volume record to query.

## Attributes Reference

The Linode Volume resource exports the following attributes:

- `id` - The unique ID of this Volume.

- `created` - When this Volume was created.

- `status` - The current status of the Volume. (`creating`, `active`, `resizing`, `contact_support`)

- `label` - This Volume's label is for display purposes only.

- `tags` - An array of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

- `size` - The Volume's size, in GiB.

- `region` - The datacenter in which this Volume is located. See all regions [here](https://api.linode.com/v4/regions).

- `updated` - When this Volume was last updated.

- `linode_id` - If a Volume is attached to a specific Linode, the ID of that Linode will be displayed here. If the Volume is unattached, this value will be null.

- `filesystem_path` - The full filesystem path for the Volume based on the Volume's label. Path is /dev/disk/by-id/scsi-0LinodeVolume + Volume label.

- `encryption` - Whether Block Storage Disk Encryption is enabled or disabled on this Volume. Note: Block Storage Disk Encryption is not currently available to all users.
