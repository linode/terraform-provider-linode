---
page_title: "Linode: linode_instance_config"
description: |-
  Manages a Linode Instance Config.
---

# linode\_instance\_config

Provides a Linode Instance Config resource. This can be used to create, modify, and delete Linode Instance Configs.

~> **NOTICE:** If a VPC interface is defined in your `linode_instance_config` resource and the config is currently booted with the Linode, then the Linode is required to be powered off during the update operation. The Terraform provider will try to implicitly shutdown you Linode instance during the update and restart it when it's finished. Unless you explicitly config the `booted` attribute in the resource or explicitly set `skip_implicit_reboots` to `false` in the Terraform provider config.

**NOTE:** Changes to a config that is currently booted will trigger a reboot, in order to skip this explicitly set `skip_implicit_reboots` to `true` in the Terraform provider config.

**NOTE:** Deleting a config will shut down the attached instance if the config is in use.

## Example Usage

Creating a simple bootable Linode Instance Configuration Profile:

```hcl
resource "linode_instance_config" "my-config" {
  linode_id = linode_instance.my-instance.id
  label = "my-config"

  device {
    device_name = "sda"
    disk_id = linode_instance_disk.boot.id
  }

  booted = true
}

resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = linode_instance.my-instance.specs.0.disk

  image = "linode/ubuntu22.04"
  root_pass = "myc00lpass!"
}

resource "linode_instance" "my-instance" {
  label = "my-instance"
  type = "g6-standard-1"
  region = "us-southeast"
}
```

Creating a complex bootable Instance Configuration Profile with a VPC:

```hcl
resource "linode_instance_config" "my-config" {
  linode_id = linode_instance.my-instance.id
  label = "my-config"

  device {
    device_name = "sda"
    disk_id = linode_instance_disk.boot.id
  }

  device {
    device_name = "sdb"
    disk_id = linode_instance_disk.swap.id
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

  # VPC networking on eth2
  interface {
    purpose = "vpc"
    subnet_id = linode_vpc_subnet.foobar.id
    ipv4 {
      vpc = "10.0.4.250"
    }
  }
  
  booted = true

  // Run a remote-exec provisioner
  connection {
    host        = linode_instance.my-instance.ip_address
    user        = "root"
    password    = "myc00lpass!ciuw23asxbviwuc"
  }

  provisioner "remote-exec" {
    inline = [
      "echo 'Hello World!'"
    ]
  }
}

# Create a VPC and a subnet
resource "linode_vpc" "foobar" {
    label = "my-vpc"
    region = "us-mia"
    description = "test description"
}

resource "linode_vpc_subnet" "foobar" {
    vpc_id = linode_vpc.foobar.id
    label = "my-subnet"
    ipv4 = "10.0.4.0/24"
}

# Create a boot disk
resource "linode_instance_disk" "boot" {
  label = "boot"
  linode_id = linode_instance.my-instance.id
  size = linode_instance.my-instance.specs.0.disk - 512

  image = "linode/ubuntu22.04"
  root_pass = "myc00lpass!ciuw23asxbviwuc"
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
  region = "us-mia"
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

* `kernel` - (Optional) A Kernel ID to boot a Linode with. Default is `linode/latest-64bit`. Examples are `linode/latest-64bit`, `linode/grub2`, `linode/direct-disk`, etc. See all kernels [here](https://api.linode.com/v4/linode/kernels). Note that this is a paginated API endpoint ([docs](https://developers.linode.com/api/v4/linode-kernels)).

* `memory_limit` - (Optional) The memory limit of the Config. Defaults to the total ram of the Linode.

* `root_device` - (Optional) The root device to boot. (default `/dev/sda`)

* `run_level` - (Optional) Defines the state of your Linode after booting. (`default`, `single`, `binbash`)

* `virt_mode` - (Optional) Controls the virtualization mode. (`paravirt`, `fullvirt`)

### devices and device

#### devices (deprecated)

The following attributes are available on devices:

* `sda` ... `sdh` - (Optional) The SDA-SDH slots, represent the Linux block device nodes for the first 8 disks attached to the Linode.  Each device must be suplied sequentially.  The device can be either a Disk or a Volume identified by `disk_id` or `volume_id`. Only one disk identifier is permitted per slot. Devices mapped from `sde` through `sdh` are unavailable in `"fullvirt"` `virt_mode`.

  * `volume_id` - (Optional) The Volume ID to map to this `device` slot.

  * `disk_id` - (Optional) The Disk ID to map to this `device` slot

#### device (recommended)

An assignment between a disk and a configuration profile device. This block supersedes the `devices` block.

Compared with `devices`, `sda` ... `sdh` is now in the `device_name` attribute in a device block, and the block itself becomes unnamed.

```terraform
device {
  device_name = "sda"
  volume_id = 1234
}

device {
  device_name = "sdb"
  disk_id = 5678
}
```

### helpers

The following attributes are available on helpers:

* `devtmpfs_automount` - (Optional) Populates the /dev directory early during boot without udev. (default `true`)

* `distro` - (Optional) Helps maintain correct inittab/upstart console device. (default `true`)

* `modules_dep` - (Optional) Creates a modules dependency file for the Kernel you run. (default `true`)

* `network` - (Optional) Automatically configures static networking. (default `true`)

* `updatedb_disabled` - (Optional) Disables updatedb cron job to avoid disk thrashing. (default `true`)

### interface

A Linode must have a public interface in the first/eth0 position to be reachable via the public internet
upon boot without additional system configuration. If no public interface is configured, the Linode
is not directly reachable via the public internet. In this case, access can only be established via
LISH or other Linodes connected to the same VLAN.

Only one public interface per Linode can be defined.

The Linode’s default public IPv4 address is assigned to the public interface.

The following arguments are available in an interface:

* `purpose` - (Required) The type of interface. (`public`, `vlan`, `vpc`)

* `ipam_address` - (Optional) This Network Interface’s private IP address in Classless Inter-Domain Routing (CIDR) notation. (e.g. `10.0.0.1/24`) This field is only allowed for interfaces with the `vlan` purpose.

* `label` - (Optional) The name of the VLAN to join. This field is only allowed and required for interfaces with the `vlan` purpose.

* `subnet_id` - (Optional) The name of the VPC Subnet to join. This field is only allowed and required for interfaces with the `vpc` purpose.

* `primary` - (Optional) Whether the interface is the primary interface that should have the default route for this Linode. This field is only allowed for interfaces with the `public` or `vpc` purpose.

* [`ipv4`](#ipv4) - (Optional) The IPv4 configuration of the VPC interface. This field is currently only allowed for interfaces with the `vpc` purpose.

The following computed attribute is available in a VPC interface:

* `vpc_id` - The ID of VPC which this interface is attached to.

* `ip_ranges` - (Optional) IPv4 CIDR VPC Subnet ranges that are routed to this Interface. IPv6 ranges are also available to select participants in the Beta program.

#### ipv4

The following arguments are available in an `ipv4` configuration block of an `interface` block:

* `vpc` - (Optional) The IP from the VPC subnet to use for this interface. A random address will be assigned if this is not specified in a VPC interface.

* `nat_1_1` - (Optional) The public IP that will be used for the one-to-one NAT purpose. If this is `any`, the public IPv4 address assigned to this Linode is used on this interface and will be 1:1 NATted with the VPC IPv4 address.

## Import

Instance Configs can be imported using the `linode_id` followed by the Instance Config `id` separated by a comma, e.g.

```sh
terraform import linode_instance_config.my-config 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
