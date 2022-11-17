---
layout: "linode"
page_title: "Linode: linode_instance_config"
sidebar_current: "docs-linode-resource-instance_config"
description: |-
Manages a Linode Instance Config.
---

# linode\_instance\_config

Provides a Linode Instance Config resource. This can be used to create, modify, and delete Linode Instance Configs.

**NOTE:** Deleting a config will shut down the attached instance if the config is in use.

## Example Usage

Creating a simple bootable Linode Instance Configuration Profile:

```hcl
resource "linode_instance_config" "my-config" {
  linode_id = linode_instance.my-instance.id
  label = "my-config"

  devices {
    sda {
      disk_id = linode_instance_disk.boot.id
    }
  }

  booted = true
}

resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = linode_instance.my-instance.specs.0.disk

  image = "linode/ubuntu20.04"
  root_pass = "myc00lpass!"
}

resource "linode_instance" "my-instance" {
  label = "my-instance"
  type = "g6-standard-1"
  region = "us-southeast"
}
```

Creating a complex bootable Instance Configuration Profile:

```hcl
resource "linode_instance_config" "my-config" {
  linode_id = linode_instance.my-instance.id
  label = "my-config"

  devices {
    sda {
      disk_id = linode_instance_disk.boot.id
    }

    sdb {
      disk_id = linode_instance_disk.swap.id
    }
  }
  
  helpers {
    # Disable the updatedb helper
    updatedb_disabled = false
  }
  
  # Public networking on eth0
  interface {
    purpose = "public"
  }
  
  # VLAN networking on eth1
  interface {
    purpose = "vlan"
    label = "my-vlan"
    ipam_address = "10.0.0.2/24"
  }
  
  booted = true
}

# Create a boot disk
resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = linode_instance.my-instance.specs.0.disk - 512

  image = "linode/ubuntu20.04"
  root_pass = "myc00lpass!"
}

# Create a swap disk
resource "linode_instance_disk" "swap" {
  label = "swap"
  linode_id = linode_instance.my-instance.id
  size = 512
  filesystem = "swap"
}

resource "linode_instance" "my-instance" {
  label = "my-instance"
  type = "g6-standard-1"
  region = "us-southeast"
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to create this configuration profile under.

* `label` - (Required) The Config’s label for display purposes only.

- - -

* `booted` - (Optional) If true, the Linode will be booted into this config. If another config is booted, the Linode will be rebooted into this config. If false, the Linode will be shutdown only if it is currently booted into this config. If undefined, the config will alter the boot status of the Linode.

* `comments` - (Optional) Optional field for arbitrary User comments on this Config.

* [`devices`](#devices) - (Optional) A dictionary of device disks to use as a device map in a Linode’s configuration profile.

* [`helpers`](#helpers) - (Optional) Helpers enabled when booting to this Linode Config.

* [`interface`](#interface) - (Optional) An array of Network Interfaces to use for this Configuration Profile.

* `kernel` - (Optional) A Kernel ID to boot a Linode with. (default `linode/latest-64bit`)

* `memory_limit` - (Optional) The memory limit of the Config. Defaults to the total ram of the Linode.

* `root_device` - (Optional) The root device to boot. (default `/dev/sda`)

* `run_level` - (Optional) Defines the state of your Linode after booting. (`default`, `single`, `binbash`)

* `virt_mode` - (Optional) Controls the virtualization mode. (`paravirt`, `fullvirt`)

### devices

The following attributes are available on devices:

* `sda` ... `sdh` - (Optional) The SDA-SDH slots, represent the Linux block device nodes for the first 8 disks attached to the Linode.  Each device must be suplied sequentially.  The device can be either a Disk or a Volume identified by `disk_id` or `volume_id`. Only one disk identifier is permitted per slot. Devices mapped from `sde` through `sdh` are unavailable in `"fullvirt"` `virt_mode`.

  * `volume_id` - (Optional) The Volume ID to map to this `device` slot.

  * `disk_id` - (Optional) The Disk ID to map to this `device` slot

### helpers

The following attributes are available on helpers:

* `devtmpfs_automount` - (Optional) Populates the /dev directory early during boot without udev. (default `true`)

* `distro` - (Optional) Helps maintain correct inittab/upstart console device. (default `true`)

* `modules_dep` - (Optional) Creates a modules dependency file for the Kernel you run. (default `true`)

* `network` - (Optional) Automatically configures static networking. (default `true`)

* `updatedb_disabled` - (Optional) Disables updatedb cron job to avoid disk thrashing. (default `true`)

### interface

The following attributes are available on interface:

* `purpose` - (Required) The type of interface. (`public`, `vlan`)

* `ipam_address` - (Optional) This Network Interface’s private IP address in Classless Inter-Domain Routing (CIDR) notation. (e.g. `10.0.0.1/24`)

* `label` - (Optional) The name of the VLAN to join. This field is only allowed for interfaces with the `vlan` purpose.

## Import

Instance Configs can be imported using the `linode_id` followed by the Instance Config `id` separated by a comma, e.g.

```sh
terraform import linode_instance_config.my-config 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
