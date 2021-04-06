---
layout: "linode"
page_title: "Linode: linode_instances"
sidebar_current: "docs-linode-datasource-instances"
description: |-
Provides information about Linode instances that match a set of filters.
---

# Data Source: linode\_instances

Provides information about Linode instances that match a set of filters.

## Example Usage

Get information about all Linode instances with a certain label and tag:

```hcl
data "linode_instances" "my-instances" {
  filter {
    name = "label"
    values = ["my-label", "my-other-label"]
  }

  filter {
    name = "tags"
    values = ["my-tag"]
  }
}
```

Get information about all Linode instances associated with the current token:

```hcl
data "linode_instances" "all-instances" {}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode instances that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

## Attributes

Each Linode instance will be stored in the `instances` attribute and will export the following attributes:

* `region` - This is the location where the Linode is deployed. Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc. See all regions [here](https://api.linode.com/v4/regions).

* `type` - The Linode type defines the pricing, CPU, disk, and RAM specs of the instance. Examples are `"g6-nanode-1"`, `"g6-standard-2"`, `"g6-highmem-16"`, `"g6-dedicated-16"`, etc. See all types [here](https://api.linode.com/v4/linode/types).

* `label` - The Linode's label is for display purposes only.
  
* `group` - The display group of the Linode instance.

* `tags` - A list of tags applied to this object. Tags are for organizational purposes only.

* `private_ip` - If true, the Linode has private networking enabled, allowing use of the 192.168.128.0/17 network within the Linode's region.
  
* `alerts.0.cpu` - The percentage of CPU usage required to trigger an alert. If the average CPU usage over two hours exceeds this value, we'll send you an alert. If this is set to 0, the alert is disabled.

* `alerts.0.network_in` - The amount of incoming traffic, in Mbit/s, required to trigger an alert. If the average incoming traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.network_out` - The amount of outbound traffic, in Mbit/s, required to trigger an alert. If the average outbound traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.transfer_quota` - The percentage of network transfer that may be used before an alert is triggered. When this value is exceeded, we'll alert you. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.io` - The amount of disk IO operation per second required to trigger an alert. If the average disk IO over two hours exceeds this value, we'll send you an alert. If set to 0, this alert is disabled.

* `watchdog_enabled` - The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have been more than 5 boot jobs issued within 15 minutes.

* `image` - An Image ID to deploy the Disk from. Official Linode Images start with linode/, while your Images start with `private/`. See [images](https://api.linode.com/v4/images) for more information on the Images available for you to use. Examples are `linode/debian9`, `linode/fedora28`, `linode/ubuntu16.04lts`, `linode/arch`, and `private/12345`. See all images [here](https://api.linode.com/v4/linode/images) (Requires a personal access token; docs [here](https://developers.linode.com/api/v4/images)). *This value can not be imported.* *Changing `image` forces the creation of a new Linode Instance.*

* `swap_size` - When deploying from an Image, this field is optional with a Linode API default of 512mb, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.

* `status` - The status of the instance, indicating the current readiness state. (`running`, `offline`, ...)

* `ip_address` - A string containing the Linode's public IP address.

* `private_ip_address` - This Linode's Private IPv4 Address, if enabled.  The regional private IP address range, 192.168.128.0/17, is shared by all Linode Instances in a region.

* `ipv6` - This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.  The prefix (`/64`) is included in this attribute.

* `ipv4` - This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon creation, and may get a single private IPv4 address if needed. You may need to open a support ticket to get additional IPv4 addresses.

* `specs.0.disk` -  The amount of storage space, in GB. this Linode has access to. A typical Linode will divide this space between a primary disk with an image deployed to it, and a swap disk, usually 512 MB. This is the default configuration created when deploying a Linode with an image through POST /linode/instances.

* `specs.0.memory` - The amount of RAM, in MB, this Linode has access to. Typically a Linode will choose to boot with all of its available RAM, but this can be configured in a Config profile.

* `specs.0.vcpus` - The number of vcpus this Linode has access to. Typically a Linode will choose to boot with all of its available vcpus, but this can be configured in a Config Profile.

* `specs.0.transfer` - The amount of network transfer this Linode is allotted each month.

* [`disk`](#disks) - A list of disks associated with the Linode.

* [`config`](#configs) - A list of configs associated with the Linode.

* [`backups`](#backups) - Information about the Linode's backup status.

### Disks

* `disk`

  * `label` - The disks label, which acts as an identifier in Terraform.  This must be unique within each Linode Instance.

  * `size` - The size of the Disk in MB.

  * `id` - The ID of the disk in the Linode API.

  * `filesystem` - The Disk filesystem can be one of: `"raw"`, `"swap"`, `"ext3"`, `"ext4"`, or `"initrd"` which has a max size of 32mb and can be used in the config `initrd` (not currently supported in this Terraform Provider).

### Configs

Configuration profiles define the VM settings and boot behavior of the Linode Instance.  Multiple configurations profiles can be provided but their `label` values must be unique.

* `config`

  * `label` - The Config's label for display purposes.  Also used by `boot_config_label`.

  * `kernel` - A Kernel ID to boot a Linode with. Default is based on image choice. Examples are `linode/latest-64bit`, `linode/grub2`, `linode/direct-disk`, etc. See all kernels [here](https://api.linode.com/v4/linode/kernels). Note that this is a paginated API endpoint ([docs](https://developers.linode.com/api/v4/linode-kernels)).

  * `run_level` - Defines the state of your Linode after booting.

  * `virt_mode` - Controls the virtualization mode.

  * `root_device` - The root device to boot.

  * `comments` - Arbitrary user comments about this `config`.

  * `memory_limit` - Defaults to the total RAM of the Linode

  * `helpers` - Helpers enabled when booting to this Linode Config.

    * `updatedb_disabled` -  Disables updatedb cron job to avoid disk thrashing.

    * `distro` -  Controls the behavior of the Linode Config's Distribution Helper setting.

    * `modules_dep` -  Creates a modules dependency file for the Kernel you run.

    * `network` -  Controls the behavior of the Linode Config's Network Helper setting, used to automatically configure additional IP addresses assigned to this instance.

  * `devices` -  A list of `disk` or `volume` attachments for this `config`.  If the `boot_config_label` omits a `devices` block, the Linode will not be booted.

    * `sda` ... `sdh` -  The SDA-SDH slots, represent the Linux block device nodes for the first 8 disks attached to the Linode.  Each device must be suplied sequentially.  The device can be either a Disk or a Volume identified by `disk_label` or `volume_id`. Only one disk identifier is permitted per slot. Devices mapped from `sde` through `sdh` are unavailable in `"fullvirt"` `virt_mode`.

    * `disk_label` -  The `label` of the `disk` to map to this `device` slot.

    * `volume_id` -  The Volume ID to map to this `device` slot.

    * `disk_id` - The Disk ID of the associated `disk_label`, if used

### Backups

* `backups`

  * `enabled` - If this Linode has the Backup service enabled.

  * `schedule`

    * `day` -  The day of the week that your Linode's weekly Backup is taken. If not set manually, a day will be chosen for you. Backups are taken every day, but backups taken on this day are preferred when selecting backups to retain for a longer period.  If not set manually, then when backups are initially enabled, this may come back as "Scheduling" until the day is automatically selected.

    * `window` - The window ('W0'-'W22') in which your backups will be taken, in UTC. A backups window is a two-hour span of time in which the backup may occur. For example, 'W10' indicates that your backups should be taken between 10:00 and 12:00. If you do not choose a backup window, one will be selected for you automatically.  If not set manually, when backups are initially enabled this may come back as Scheduling until the window is automatically selected.
  
## Filterable Fields

* `group`

* `id`

* `image`

* `label`

* `region`

* `tags`
