---
layout: "linode"
page_title: "Linode: linode_volume"
sidebar_current: "docs-linode-resource-volume"
description: |-
  Manages a Linode Volume.
---

# linode\_volume

Provides a Linode volume resource.  This can be used to create,
modify, and delete Linodes Volumes. For more information, see [How to Use Block Storage with Your Linode](https://www.linode.com/docs/platform/block-storage/how-to-use-block-storage-with-your-linode/)
and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createVolume).

## Example Usage

The following example shows how one might use this resource to configure a Volume attached to a Linode instance.

```hcl
resource "linode_instance" "foobaz" {
    root_pass = "3X4mp13"
    type = "g6-nanode-1"
    region = "us-west"
}

resource "linode_volume" "foobar" {
    label = "foo-volume"
    region = "${linode_instance.foobaz.region}"
    linode_id = "${linode_instance.foobaz.id}"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Linode Volume

* `region` - (Required) The region where this volume will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode Volume.*.

- - -

* `size` - (Optional) Size of the Volume in GB.

* `linode_id` - (Optional) The ID of a Linode Instance where the the Volume should be attached.

## Attributes

This resource exports the following attributes:

* `status` - The label of the Linode Volume.

* `filesystem_path` - The full filesystem path for the Volume based on the Volume's label. The path is "/dev/disk/by-id/scsi-0Linode_Volume_" + the Volume label

## Import

Linodes Volumes can be imported using the Linode Volume `id`, e.g.

```sh
terraform import linode_volume.myvolume 1234567
```
