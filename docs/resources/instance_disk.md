---
page_title: "Linode: linode_instance_disk"
description: |-
  Manages a Linode Instance Disk.
---

# linode\_instance\_disk

Provides a Linode Instance Disk resource. This can be used to create, modify, and delete Linode Instance Disks.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-add-linode-disk).

**NOTE:** Deleting a disk will shut down the attached instance if the instance is booted. If the disk was not in use by the booted configuration profile, the instance will be automatically rebooted.

## Example Usage

Creating a simple 512 MB Linode Instance Disk:

```hcl
resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = 512
  filesystem = "ext4"
}

resource "linode_instance" "my-instance" {
  label = "my-instance"
  type = "g6-standard-1"
  region = "us-southeast"
}
```

Creating a complex bootable Instance Disk:

```hcl
resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = linode_instance.my-instance.specs.0.disk

  image = "linode/ubuntu22.04"
  root_pass = "myc00lpass!"
  authorized_keys = ["ssh-rsa AAAA...Gw== user@example.local"]
  
  # Optional StackScript to run on first boot
  stackscript_id = 12345
  stackscript_data = {
    "my_var" = "my_value"
  }
}

resource "linode_instance" "my-instance" {
  label = "my-instance"
  type = "g6-standard-1"
  region = "us-southeast"
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to create this Disk under.

* `label` - (Required) The Disk's label for display purposes only.

* `size` - (Required) The size of the Disk in MB. **NOTE:** Resizing a disk will trigger a Linode reboot.

- - -

* `authorized_keys` - (Optional) A list of public SSH keys that will be automatically appended to the root user’s ~/.ssh/authorized_keys file when deploying from an Image. (Requires `image`)

* `authorized_users` - (Optional) A list of usernames. If the usernames have associated SSH keys, the keys will be appended to the root user's ~/.ssh/authorized_keys file. (Requires `image`)

* `filesystem` - (Optional) The filesystem of this disk. (`raw`, `swap`, `ext3`, `ext4`, `initrd`)

* `image` - (Optional) An Image ID to deploy the Linode Disk from.

* `root_pass` - (Optional) The root user’s password on a newly-created Linode Disk when deploying from an Image. (Requires `image`)

* `stackscript_data` - (Optional) An object containing responses to any User Defined Fields present in the StackScript being deployed to this Disk. Only accepted if `stackscript_id` is given. (Requires `image`)

* `stackscript_id` - (Optional) A StackScript ID that will cause the referenced StackScript to be run during deployment of this Disk. (Requires `image`)

## Attributes Reference

This resource exports the following attributes:

* `created` - When this disk was created.

* `disk_encryption` - The disk encryption policy for this disk's parent instance. (`enabled`, `disabled`)

  * **NOTE: Disk encryption may not currently be available to all users.**

* `status` - A brief description of this Disk's current state.

* `updated` - When this disk was last updated.

## Import

Instance Disks can be imported using the `linode_id` followed by the Instance Disk `id` separated by a comma, e.g.

```sh
terraform import linode_instance_disk.my-disk 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
