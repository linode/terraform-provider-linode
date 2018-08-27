---
layout: "linode"
page_title: "Linode: linode_instance"
sidebar_current: "docs-linode-resource-instance"
description: |-
  Manages a Linode instance.
---

# linode\_instance

Provides a Linode instance resource.  This can be used to create,
modify, and delete Linodes. For more information, see [Getting Started with Linode](https://linode.com/docs/getting-started/)
and [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createLinodeInstance).

Linodes also support `[provisioning](/docs/provisioners/index.html).

## Example Usage

The following example shows how one might use this resource to configure a Linode instance.

```hcl
resource "linode_instance" "web" {
    image = "linode/ubuntu18.04"
    kernel = "linode/latest-64"
    region = "us-central"
    type = "g6-standard-1"
    ssh_key = "ssh-rsa AAAA...Gw== user@example.local"
    root_password = "terraform-test"

    label = "foobaz"
    group = "integration"
    status = "on"
    swap_size = 256
    private_networking = true

    // ip_address = "8.8.8.8"
    // plan_storage = 24576
    // plan_storage_utilized = 24576
    // private_ip_address = "192.168.10.50"
}
```

## Argument Reference

The following arguments are supported:

* `image` - (Required) The image to use when creating the Linode's disks. Examples are `"linode/debian9"`, `"linode/fedora28"`, and `"linode/arch"`. *Changing `image` forces the creation of a new Linode Instance.*

* `kernel` - (Required) The kernel to start the linode with. Specify `"linode/latest-64bit"` or `"linode/latest-32bit""` for the most recent Linode provided kernel. "linode/direct-disk" can be used to boot the raw disk and "linode/grub2" will boot to the Grub config on the disk.

* `region` - (Required) The region that the linode will be created in.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode Instance.*.

* `type` - (Required) The Linode type defines the pricing, CPU, disk, and RAM specs of the instance.  Examples are `"g6-nanode-1"`, `"g6-standard-2"`, `"g6-highmem-16"`, etc.

* `ssh_key` - (Required) The full text of the public key to add to the root user. *Changing `ssh_key` forces the creation of a new Linode Instance.*

* `root_password` - (Required) The initial password for the `root` user account. *Changing `ssh_key` forces the creation of a new Linode Instance.*

  A `root_password` is required by the Linode API. You'll likely want to modify this on the server during provisioning and then disable password logins in favor of SSH keys.

- - -

* `label` - (Optional) The label of the Linode.

* `group` - (Optional) The group of the Linode.

* `private_networking` - (Optional) A boolean controlling whether or not to enable private networking. It can be enabled on an existing Linode but it can't be disabled.

* `helper_distro` - (Optional) A boolean used to enable the Distro Filesystem helper.   This corrects fstab and inittab/upstart entries depending on the distribution or kernel being booted. You want this unless you're providing your own kernel.

* `manage_private_ip_automatically` - (Optional) A boolean used to enable the Network Helper.  This automatically creates network configuration files for your distro and places them into your filesystem. Enabling this in a change will reboot your Linode.

* `disk_expansion` - (Optional) A boolean that when true will automatically expand the root volume if the size of the Linode plan is increased.  Setting this value will prevent downsizing without manually shrinking the volume prior to decreasing the size.

* `swap_size` - (Optional) Sets the size of the swap partition on a Linode in MB.  At this time, this cannot be modified by Terraform after initial provisioning.  If manually modified via the Web GUI, this value will reflect such modification.  This value can be set to 0 to create a Linode without a swap partition.  Defaults to 512.

## Attributes

This resource exports the following attributes:

* `status` - A string representing the power status of the Linode (`"on"`, `"off"`)

* `ip_address` - A string containing the Linode's public IP address.

* `private_ip_address` - A string containing the Linode's private IP address if private networking is enabled.

* `plan_storage` - An integer reflecting the size of the Linode's storage capacity in MB, based on the Linode plan.

* `plan_storage_utilized` - An integer sum of the size of all the Linode's disks, given in MB.

## Import

Linodes Instances can be imported using the Linode `id`, e.g.

```sh
terraform import linode_instance.mylinode 1234567
```
