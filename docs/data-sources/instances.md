---
page_title: "Linode: linode_instances"
description: |-
  Provides information about Linode instances that match a set of filters.
---

# Data Source: linode\_instances

Provides information about Linode instances that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-linode-instances).

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

output "instance_id" {
  value = data.linode_instances.my-instances.instances.0.id
}
```

Get information about all Linode instances associated with the current token:

```hcl
data "linode_instances" "all-instances" {}

output "instance_ids" {
  value = data.linode_instances.all-instances.instances.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode instances that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode instance will be stored in the `instances` attribute and will export the following attributes:

* `id` - The ID of the Linode instance.

* `region` - This is the location where the Linode is deployed. Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc. See all regions [here](https://api.linode.com/v4/regions).

* `type` - The Linode type defines the pricing, CPU, disk, and RAM specs of the instance. Examples are `"g6-nanode-1"`, `"g6-standard-2"`, `"g6-highmem-16"`, `"g6-dedicated-16"`, etc. See all types [here](https://api.linode.com/v4/linode/types).

* `label` - The Linode's label is for display purposes only.
  
* `group` - The display group of the Linode instance.

* `tags` - A list of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

* `private_ip` - If true, the Linode has private networking enabled, allowing use of the 192.168.128.0/17 network within the Linode's region.
  
* `alerts.0.cpu` - The percentage of CPU usage required to trigger an alert. If the average CPU usage over two hours exceeds this value, we'll send you an alert. If this is set to 0, the alert is disabled.

* `alerts.0.network_in` - The amount of incoming traffic, in Mbit/s, required to trigger an alert. If the average incoming traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.network_out` - The amount of outbound traffic, in Mbit/s, required to trigger an alert. If the average outbound traffic over two hours exceeds this value, we'll send you an alert. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.transfer_quota` - The percentage of network transfer that may be used before an alert is triggered. When this value is exceeded, we'll alert you. If this is set to 0 (zero), the alert is disabled.

* `alerts.0.io` - The amount of disk IO operation per second required to trigger an alert. If the average disk IO over two hours exceeds this value, we'll send you an alert. If set to 0, this alert is disabled.

* `watchdog_enabled` - The watchdog, named Lassie, is a Shutdown Watchdog that monitors your Linode and will reboot it if it powers off unexpectedly. It works by issuing a boot job when your Linode powers off without a shutdown job being responsible. To prevent a loop, Lassie will give up if there have been more than 5 boot jobs issued within 15 minutes.

* `image` - An Image ID to deploy the Disk from. Official Linode Images start with linode/, while your Images start with `private/`. See [images](https://api.linode.com/v4/images) for more information on the Images available for you to use. Examples are `linode/debian12`, `linode/fedora39`, `linode/ubuntu22.04`, `linode/arch`, and `private/12345`. See all images [here](https://api.linode.com/v4/linode/images) (Requires a personal access token; docs [here](https://techdocs.akamai.com/linode-api/reference/get-images)). *This value can not be imported.* *Changing `image` forces the creation of a new Linode Instance.*

* `swap_size` - When deploying from an Image, this field is optional with a Linode API default of 512mb, otherwise it is ignored. This is used to set the swap disk size for the newly-created Linode.

* `status` - The status of the instance, indicating the current readiness state. (`running`, `offline`, ...)

* `ip_address` - A string containing the Linode's public IP address.

* `private_ip_address` - This Linode's Private IPv4 Address, if enabled.  The regional private IP address range, 192.168.128.0/17, is shared by all Linode Instances in a region.

* `ipv6` - This Linode's IPv6 SLAAC addresses. This address is specific to a Linode, and may not be shared.  The prefix (`/64`) is included in this attribute.

* `ipv4` - This Linode's IPv4 Addresses. Each Linode is assigned a single public IPv4 address upon creation, and may get a single private IPv4 address if needed. You may need to open a support ticket to get additional IPv4 addresses.

* `has_user_data` - Whether this Instance was created with user-data.

* `disk_encryption` - The disk encryption policy for this instance.

  * **NOTE: Disk encryption may not currently be available to all users.**

* `lke_cluster_id` - If applicable, the ID of the LKE cluster this instance is a part of.

* `specs.0.disk` -  The amount of storage space, in GB. this Linode has access to. A typical Linode will divide this space between a primary disk with an image deployed to it, and a swap disk, usually 512 MB. This is the default configuration created when deploying a Linode with an image through POST /linode/instances.

* `specs.0.memory` - The amount of RAM, in MB, this Linode has access to. Typically a Linode will choose to boot with all of its available RAM, but this can be configured in a Config profile.

* `specs.0.vcpus` - The number of vcpus this Linode has access to. Typically a Linode will choose to boot with all of its available vcpus, but this can be configured in a Config Profile.

* `specs.0.transfer` - The amount of network transfer this Linode is allotted each month.

* [`disk`](#disks) - A list of disks associated with the Linode.

* [`config`](#configs) - A list of configs associated with the Linode.

* [`backups`](#backups) - Information about the Linode's backup status.

* [`placement_group`](#placement-groups) - Information about the Linode's Placement Groups.

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

  * `kernel` - A Kernel ID to boot a Linode with. Default is based on image choice. Examples are `linode/latest-64bit`, `linode/grub2`, `linode/direct-disk`, etc. See all kernels [here](https://api.linode.com/v4/linode/kernels). Note that this is a paginated API endpoint ([docs](https://techdocs.akamai.com/linode-api/reference/get-kernels)).

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

  * [`interface`](#interface) - (Optional) A list of network interfaces to be assigned to the Linode.

### Interface

Interface defines a network interfaces that is exposed to a Linode. See the official [Linode API documentation](https://techdocs.akamai.com/linode-api/reference/post-linode-config-interface) for more details.

Each interface exports the following attributes:

* `purpose` - The type of interface. (`public`, `vlan`, `vpc`)

* `ipam_address` - This Network Interface’s private IP address in Classless Inter-Domain Routing (CIDR) notation. (e.g. `10.0.0.1/24`) This field is only allowed for interfaces with the `vlan` purpose.

* `label` - The name of the VLAN to join. This field is only allowed and required for interfaces with the `vlan` purpose.

* `subnet_id` - The name of the VPC Subnet to join. This field is only allowed and required for interfaces with the `vpc` purpose.

* `primary` - Whether the interface is the primary interface that should have the default route for this Linode. This field is only allowed for interfaces with the `public` or `vpc` purpose.

* [`ipv4`](#ipv4) -The IPv4 configuration of the VPC interface. This field is currently only allowed for interfaces with the `vpc` purpose.

* `vpc_id` - The ID of VPC which this interface is attached to.

* `ip_ranges` - IPv4 CIDR VPC Subnet ranges that are routed to this Interface. IPv6 ranges are also available to select participants in the Beta program.

#### ipv4

The following arguments are available in an `ipv4` configuration block of an `interface` block:

* `vpc` - The IP from the VPC subnet to use for this interface. A random address will be assigned if this is not specified in a VPC interface.

* `nat_1_1` - The public IP that will be used for the one-to-one NAT purpose. If this is `any`, the public IPv4 address assigned to this Linode is used on this interface and will be 1:1 NATted with the VPC IPv4 address.

### Backups

* `backups`

  * `enabled` - If this Linode has the Backup service enabled.

  * `schedule`

    * `day` -  The day of the week that your Linode's weekly Backup is taken. If not set manually, a day will be chosen for you. Backups are taken every day, but backups taken on this day are preferred when selecting backups to retain for a longer period.  If not set manually, then when backups are initially enabled, this may come back as "Scheduling" until the day is automatically selected.

    * `window` - The window ('W0'-'W22') in which your backups will be taken, in UTC. A backups window is a two-hour span of time in which the backup may occur. For example, 'W10' indicates that your backups should be taken between 10:00 and 12:00. If you do not choose a backup window, one will be selected for you automatically.  If not set manually, when backups are initially enabled this may come back as Scheduling until the window is automatically selected.

### Placement Groups

* `placement_group`

  * `id` -  The ID of the Placement Group in the Linode API.

  * `placement_group_type` - The placement group type to use when placing Linodes in this group.

  * `placement_group_policy` - Whether Linodes must be able to become compliant during assignment. (Default `strict`)

  * `label` - The label of the Placement Group. This field can only contain ASCII letters, digits and dashes.

## Filterable Fields

* `group`

* `id`

* `image`

* `label`

* `region`

* `status`

* `tags`

* `type`

* `watchdog_enabled`
