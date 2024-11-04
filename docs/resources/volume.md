---
page_title: "Linode: linode_volume"
description: |-
  Manages a Linode Volume.
---

# linode\_volume

Provides a Linode Volume resource.  This can be used to create, modify, and delete Linodes Block Storage Volumes.  Block Storage Volumes are removable storage disks that persist outside the life-cycle of Linode Instances. These volumes can be attached to and detached from Linode instances throughout a region.

For more information, see [How to Use Block Storage with Your Linode](https://www.linode.com/docs/platform/block-storage/how-to-use-block-storage-with-your-linode/) and the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-volume).

## Example Usage

The following example shows how one might use this resource to configure a Block Storage Volume attached to a Linode Instance.

```hcl
resource "linode_instance" "foobaz" {
    root_pass = "3X4mp13"
    type = "g6-nanode-1"
    region = "us-west"
    tags = ["foobaz"]

}

resource "linode_volume" "foobar" {
    label = "foo-volume"
    region = linode_instance.foobaz.region
    linode_id = linode_instance.foobaz.id
}
```

Volumes can also be attached using the Linode Instance config device map.

```hcl
resource "linode_instance" "foo" {
  region             = "us-east"
  type               = "g6-nanode-1"
}

resource "linode_instance_config" "foo" {
  linode_id = linode_instance.foo.id
  label = "boot-existing-volume"
  kernel = "linode/grub2"

  device {
    device_name = "sda"
    volume_id = 12345
  }

  booted = true
}
```

Volumes may also be cloned from existing volumes.

```hcl
resource "linode_volume" "foobar" {
  label = "my-cloned-volume"
  source_volume_id = 12345
  
  # Region is optional when cloning a volume
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Linode Volume

* `region` - (Required) The region where this volume will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc. See all regions [here](https://api.linode.com/v4/regions). This field is optional for cloned volumes. *Changing `region` forces the creation of a new Linode Volume.*.

- - -

* `source_volume_id` - (Optional) The ID of a Linode Volume to clone. NOTE: Cloned volumes must be in the same region as the source volume.

* `size` - (Optional) Size of the Volume in GB.

* `linode_id` - (Optional) The ID of a Linode Instance where the Volume should be attached.

* `tags` - (Optional) A list of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

* `encryption` - (Optional) Whether Block Storage Disk Encryption is enabled or disabled on this Volume. Note: Block Storage Disk Encryption is not currently available to all users.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the volume (until the volume is reaches the initial `active` state)
* `update` - (Defaults to 20 mins) Used when updating the volume when necessary during update - e.g. when resizing the volume
* `delete` - (Defaults to 10 mins) Used when deleting the volume

## Attributes Reference

This resource exports the following attributes:

* `status` - The status of the Linode Volume. (`creating`, `active`, `resizing`, `contact_support`)

* `filesystem_path` - The full filesystem path for the Volume based on the Volume's label. The path is "/dev/disk/by-id/scsi-0Linode_Volume_" + the Volume label

## Import

Linodes Volumes can be imported using the Linode Volume `id`, e.g.

```sh
terraform import linode_volume.myvolume 1234567
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for Block Storage Volumes and other Linode resource types.
