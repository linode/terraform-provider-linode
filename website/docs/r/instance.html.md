---
layout: "linode"
page_title: "Linode: linode_instance"
sidebar_current: "docs-linode-resource-instance"
description: |-
  Manages a Linode instance.
---

# linode\_instance

Provides a Linode Instance resource.  This can be used to create, modify, and delete Linodes.
For more information, see [Getting Started with Linode](https://linode.com/docs/getting-started/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createLinodeInstance).

The Linode Guide, [Use Terraform to Provision Linode Environments](https://www.linode.com/docs/applications/configuration-management/how-to-build-your-infrastructure-using-terraform-and-linode/), provides step-by-step guidance and additional examples.

Linode Instances can also use [provisioners](/docs/provisioners/index.html).

## Example Usage

### Simple Linode Instance

The following example shows how one might use this resource to configure a Linode instance.

```hcl
resource "linode_instance" "web" {
    label = "simple_instance"
    image = "linode/ubuntu18.04"
    region = "us-central"
    type = "g6-standard-1"
    authorized_keys = ["ssh-rsa AAAA...Gw== user@example.local"]
    root_pass = "terr4form-test"

    group = "foo"
    tags = [ "foo" ]
    swap_size = 256
    private_ip = true
}
```

### Linode Instance with explicit Configs and Disks

Using explicit Instance Configs and Disks it is possible to create a more elaborate Linode instance.  This can be used to provision multiple disks and volumes during Instance creation.

```hcl
data "linode_profile" "me" {}

resource "linode_instance" "web" {
  label      = "complex_instance"
  group      = "foo"
  tags = [ "foo" ]
  region     = "us-central"
  type       = "g6-nanode-1"
  private_ip = true

  disk {
    label = "boot"
    size = 3000
    image  = "linode/ubuntu18.04"

    # Any of authorized_keys, authorized_users, and root_pass
    # can be used for provisioning.
    authorized_keys = [ "ssh-rsa AAAA...Gw== user@example.local" ]
    authorized_users = [ "${data.linode_profile.me.username}" ]
    root_pass = "terr4form-test"
  }

  config {
    label = "boot_config"
    kernel = "linode/latest-64bit"
    devices {
      sda {
        disk_label = "boot"
      }
      sdb {
        volume_id = "${linode_volume.web_volume.id}"
      }
    }
    root_device = "/dev/sda"
  }

  boot_config_label = "boot_config"
}

resource "linode_volume" "web_volume" {
  label = "web_volume"
  size = 20
  region = "us-central"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) This is the location where the Linode is deployed. Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode Instance.*.

* `type` - (Required) The Linode type defines the pricing, CPU, disk, and RAM specs of the instance.  Examples are `"g6-nanode-1"`, `"g6-standard-2"`, `"g6-highmem-16"`, `"g6-dedicated-16"`, etc.

* `label` - (Optional) The Linode's label is for display purposes only. If no label is provided for a Linode, a default will be assigned.

* `group` - (Optional) The display group of the Linode instance.

* `tags` - (Optional) A list of tags applied to this object. Tags are for organizational purposes only.

* `private_ip` - (Optional) If true, the created Linode will have private networking enabled, allowing use of the 192.168.128.0/17 network within the Linode's region. It can be enabled on an existing Linode but it can't be disabled.

* `alerts.0.cpu` - (Optional) The percentage of CPU usage required to trigger an alert. If the average CPU usage over two hours exceeds this value, we'll send you an alert. If this is set to 0, the alert is disabled.

* `alerts.0.network_in` - (Optional) The amount of incoming traffic, in Mbit/s, required to trigger an alert. If the average incoming traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.network_out` - (Optional) The amount of outbound traffic, in Mbit/s, required to trigger an alert. If the average outbound traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.transfer_quota` - (Optional) The percentage of network transfer that may be used before an alert is triggered. When this value is exceeded, we'll alert you. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.io` - (Optional) The amount of disk IO operation per second required to trigger an alert. If the average disk IO over two hours exceeds this value, we'll send you an alert. If set to 0, this alert is disabled.

* `backups_enabled` - (Optional) If this field is set to true, the created Linode will automatically be enrolled in the Linode Backup service. This will incur an additional charge. The cost for the Backup service is dependent on the Type of Linode deployed.

* `watchdog_enabled` - (Optional) The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have been more than 5 boot jobs issued within 15 minutes.

### Simplified Resource Arguments

Just as the Linode API provides, these fields are for the most common provisioning use case, a single data disk, a single swap disk, and a single config.  These arguments are not compatible with `disk` and `config` fields, described later.

* `authorized_keys` - (Optional with `image`) A list of SSH public keys to deploy for the root user on the newly created Linode. *This value can not be imported.* *Changing `authorized_keys` forces the creation of a new Linode Instance.*

* `authorized_users` - (Optional with `image`) A list of Linode usernames. If the usernames have associated SSH keys, the keys will be appended to the `root` user's `~/.ssh/authorized_keys` file automatically. *This value can not be imported.* *Changing `authorized_users` forces the creation of a new Linode Instance.*

* `root_pass` - (Optional) The initial password for the `root` user account. *This value can not be imported.* *Changing `root_pass` forces the creation of a new Linode Instance.* *If omitted, a random password will be generated but will not be stored in Terraform state.*

* `image` - (Optional) An Image ID to deploy the Disk from. Official Linode Images start with linode/, while your Images start with `private/`. See [images](https://api.linode.com/v4/images) for more information on the Images available for you to use. Examples are `linode/debian9`, `linode/fedora28`, `linode/ubuntu16.04lts`, `linode/arch`, and `private/12345`. *This value can not be imported.* *Changing `image` forces the creation of a new Linode Instance.*

* `stackscript_id` - (Optional) The StackScript to deploy to the newly created Linode. If provided, 'image' must also be provided, and must be an Image that is compatible with this StackScript. *This value can not be imported.* *Changing `stackscript_id` forces the creation of a new Linode Instance.*

* `stackscript_data` - (Optional) An object containing responses to any User Defined Fields present in the StackScript being deployed to this Linode. Only accepted if 'stackscript_id' is given. The required values depend on the StackScript being deployed.  *This value can not be imported.* *Changing `stackscript_data` forces the creation of a new Linode Instance.*

* `swap_size` - (Optional) When deploying from an Image, this field is optional with a Linode API default of 512mb, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.

* `backup_id` - (Optional) A Backup ID from another Linode's available backups. Your User must have read_write access to that Linode, the Backup must have a status of successful, and the Linode must be deployed to the same region as the Backup. See /linode/instances/{linodeId}/backups for a Linode's available backups. This field and the image field are mutually exclusive. *This value can not be imported.* *Changing `backup_id` forces the creation of a new Linode Instance.*

### Disk and Config Arguments

By specifying the `disk` and `config` fields for a Linode instance, it is possible to use non-standard kernels, boot with and provision multiple disks, and modify the boot behaviors (`helpers`) of the Linode.

* `boot_config_label` - (Optional) The Label of the Instance Config that should be used to boot the Linode instance.  If there is only one `config`, the `label` of that `config` will be used as the `boot_config_label`. *This value can not be imported.*

#### Disks

* `disk`

  * `label` - (Required) The disks label, which acts as an identifier in Terraform.  This must be unique within each Linode Instance.

  * `size` - (Required) The size of the Disk in MB.

  * `id` - (Computed) The ID of the disk in the Linode API.

  * `filesystem` - (Optional) The Disk filesystem can be one of: `"raw"`, `"swap"`, `"ext3"`, `"ext4"`, or `"initrd"` which has a max size of 32mb and can be used in the config `initrd` (not currently supported in this Terraform Provider).

  * `readonly` - (Optional) If true, this Disk is read-only.

  * `image` - (Optional) An Image ID to deploy the Disk from. Official Linode Images start with linode/, while your Images start with private/. See /images for more information on the Images available for you to use. Examples are `linode/debian9`, `linode/fedora28`, `linode/ubuntu16.04lts`, `linode/arch`, and `private/12345`. *Changing `image` forces the creation of a new Linode Instance.*

  * `authorized_keys` - (Optional with `image`) A list of SSH public keys to deploy for the root user on the newly created Linode. Only accepted if `image` is provided. *This value can not be imported.* *Changing `authorized_keys` forces the creation of a new Linode Instance.*

  * `authorized_users` - (Optional with `image`) A list of Linode usernames. If the usernames have associated SSH keys, the keys will be appended to the `root` user's `~/.ssh/authorized_keys` file automatically. *This value can not be imported.* *Changing `authorized_users` forces the creation of a new Linode Instance.*

  * `root_pass` - (Optional with `image`) The initial password for the `root` user account. *This value can not be imported.* *Changing `root_pass` forces the creation of a new Linode Instance.* *If omitted, a random password will be generated but will not be stored in Terraform state.*

  * `stackscript_id` - (Optional with `image`) The StackScript to deploy to the newly created Linode. If provided, 'image' must also be provided, and must be an Image that is compatible with this StackScript. *This value can not be imported.* *Changing `stackscript_id` forces the creation of a new Linode Instance.*

  * `stackscript_data` - (Optional with `image`) An object containing responses to any User Defined Fields present in the StackScript being deployed to this Linode. Only accepted if 'stackscript_id' is given. The required values depend on the StackScript being deployed.  *This value can not be imported.* *Changing `stackscript_data` forces the creation of a new Linode Instance.*

#### Configs

Configuration profiles define the VM settings and boot behavior of the Linode Instance.  Multiple configurations profiles can be provided but their `label` values must be unique.

* `config`

  * `label` - (Required) The Config's label for display purposes.  Also used by `boot_config_label`.

  * `helpers` - (Options) Helpers enabled when booting to this Linode Config.

    * `updatedb_disabled` - (Optional) Disables updatedb cron job to avoid disk thrashing.

    * `distro` - (Optional) Controls the behavior of the Linode Config's Distribution Helper setting.

    * `modules_dep` - (Optional) Creates a modules dependency file for the Kernel you run.

    * `network` - (Optional) Controls the behavior of the Linode Config's Network Helper setting, used to automatically configure additional IP addresses assigned to this instance.

  * `devices` - (Optional) A list of `disk` or `volume` attachments for this `config`.  If the `boot_config_label` omits a `devices` block, the Linode will not be booted.

    * `sda` ... `sdh` - (Optional) The SDA-SDH slots, represent the Linux block device nodes for the first 8 disks attached to the Linode.  Each device must be suplied sequentially.  The device can be either a Disk or a Volume identified by `disk_label` or `volume_id`. Only one disk identifier is permitted per slot. Devices mapped from `sde` through `sdh` are unavailable in `"fullvirt"` `virt_mode`.

      * `disk_label` - (Optional) The `label` of the `disk` to map to this `device` slot.

      * `volume_id` - (Optional) The Volume ID to map to this `device` slot.

      * `disk_id` - (Computed) The Disk ID of the associated `disk_label`, if used.

    * `kernel` - (Optional) - A Kernel ID to boot a Linode with. Default is based on image choice. (examples: linode/latest-64bit, linode/grub2, linode/direct-disk)

    * `run_level` - (Optional) - Defines the state of your Linode after booting. Defaults to `"default"`.

    * `virt_mode` - (Optional) - Controls the virtualization mode. Defaults to `"paravirt"`.

    * `root_device` - (Optional) - The root device to boot. The corresponding disk must be attached to a `device` slot.  Example: `"/dev/sda"`

    * `comments` - (Optional) - Arbitrary user comments about this `config`.

    * `memory_limit` - (Optional) - Defaults to the total RAM of the Linode

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when launching the instance (until it reaches the initial `running` state)
* `update` - (Defaults to 20 mins) Used when stopping and starting the instance when necessary during update - e.g. when changing instance type
* `delete` - (Defaults to 10 mins) Used when terminating the instance

## Attributes

This Linode Instance resource exports the following attributes:

* `status` - The status of the instance, indicating the current readiness state. (`running`, `offline`, ...)

* `ip_address` - A string containing the Linode's public IP address.

* `private_ip_address` - This Linode's Private IPv4 Address, if enabled.  The regional private IP address range, 192.168.128.0/17, is shared by all Linode Instances in a region.

* `ipv6` - This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.  The prefix (`/64`) is included in this attribute.

* `ipv4` - This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon creation, and may get a single private IPv4 address if needed. You may need to open a support ticket to get additional IPv4 addresses.

* `specs.0.disk` -  The amount of storage space, in GB. this Linode has access to. A typical Linode will divide this space between a primary disk with an image deployed to it, and a swap disk, usually 512 MB. This is the default configuration created when deploying a Linode with an image through POST /linode/instances.

* `specs.0.memory` - The amount of RAM, in MB, this Linode has access to. Typically a Linode will choose to boot with all of its available RAM, but this can be configured in a Config profile.

* `specs.0.vcpus` - The number of vcpus this Linode has access to. Typically a Linode will choose to boot with all of its available vcpus, but this can be configured in a Config Profile.

* `specs.0.transfer` - The amount of network transfer this Linode is allotted each month.

* `backups` - Information about this Linode's backups status.

  * `enabled` - If this Linode has the Backup service enabled.

  * `schedule`

    * `day` -  The day of the week that your Linode's weekly Backup is taken. If not set manually, a day will be chosen for you. Backups are taken every day, but backups taken on this day are preferred when selecting backups to retain for a longer period.  If not set manually, then when backups are initially enabled, this may come back as "Scheduling" until the day is automatically selected.

    * `window` - The window ('W0'-'W22') in which your backups will be taken, in UTC. A backups window is a two-hour span of time in which the backup may occur. For example, 'W10' indicates that your backups should be taken between 10:00 and 12:00. If you do not choose a backup window, one will be selected for you automatically.  If not set manually, when backups are initially enabled this may come back as Scheduling until the window is automatically selected.

## Import

Linodes Instances can be imported using the Linode `id`, e.g.

```sh
terraform import linode_instance.mylinode 1234567
```

When importing an instance, all `disk` and `config` values must be represented.

Imported disks must include their `label` value.  **Any disk that is not precisely represented may be removed resulting in data loss.**

Imported configs should include all `devices`, and must include `label`, `kernel`, and the `root_device`.  The instance must include a `boot_config_label` referring to the correct configuration profile.

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for Instances and other Linode resource types.
