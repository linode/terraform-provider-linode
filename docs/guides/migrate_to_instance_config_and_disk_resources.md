# How to Migrate Embedded Attributes to Terraform-Managed Resources

In this guide, we demonstrate how to migrate the `config` and `disk`
attributes from the `linode_instance` resource to the `linode_instance_config`
and `linode_instance_disk` resources, respectively.

-> An experimental feature, generating configuration, will be used in this guide.
For details, please check
[Terraform document site](https://developer.hashicorp.com/terraform/language/import/generating-configuration).

## Create a simple instance with embedded config and disk

This is a sample Terraform config file for a Linode instance containing deprecated
and embedded `config` and `disk` fields, let's say `main.tf`.

This doc will explain how to migrate them to be
`linode_instance_config` and `linode_instance_disk`
resources.

main.tf

```terraform
terraform {
  required_providers {
    linode = {
      source = "linode/linode"
    }
  }
}

provider "linode" {
}

resource "linode_instance" "test_instance" {
  label = "simple_instance"
  region = "us-southeast"
  type   = "g6-nanode-1"
  private_ip = true

  # the embedded config block attribute that needs to be migrated
  config {
    label  = "test-config-1"
    kernel = "linode/grub2"

    interface {
      purpose = "public"
    }

    interface {
      purpose = "vlan"
      label = "test-vlan1"
    }
    devices {
      sda {
        disk_label = "test-disk"
      }

      sdb {
        disk_label = "test-swap"
      }
    }
  }

  # the embedded disk block attribute that needs to be migrated
  disk {
    label      = "test-swap"
    size       = 1024
    filesystem = "swap"
  }

  # the embedded disk block attribute that needs to be migrated
  disk {
    label            = "test-disk"
    size             = 24576
    filesystem       = "ext4"
    authorized_users = ["zliang27"]
    root_pass        = "this_is_not_a_safe_password"
    image            = "linode/ubuntu22.04"
  }
}
```

## Import linode_instance_config and generate the config file

Add import statements for the `linode_instance_config` resources
then remove the embedded config from the `linode_instance` resource.

Now we can run `terraform plan -generate-config-out=generated.tf`
command to generate the config files.

This step will result the following Terraform config files.

main.tf

```terraform
terraform {
  required_providers {
    linode = {
      source = "linode/linode"
    }
  }
}

provider "linode" {
}

resource "linode_instance" "test_instance" {
  label = "simple_instance"
  region = "us-southeast"
  type   = "g6-nanode-1"
  private_ip = true
}

import {
  to = linode_instance_config.test_imported_config_1
  id = "49667786,52599108"
}

import {
  to = linode_instance_disk.test_imported_disk
  id = "49667786,98757139"
}

import {
  to = linode_instance_disk.test_imported_disk_swap
  id = "49667786,98757126"
}

```

generated.tf

```terraform
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "49667786,98757139"
resource "linode_instance_disk" "test_imported_disk" {
  authorized_keys  = null
  authorized_users = null
  filesystem       = "ext4"
  image            = null
  label            = "test-disk"
  linode_id        = 49667786
  root_pass        = null # sensitive
  size             = 24576
  stackscript_data = null # sensitive
  stackscript_id   = null
}

# __generated__ by Terraform from "49667786,98757126"
resource "linode_instance_disk" "test_imported_disk_swap" {
  authorized_keys  = null
  authorized_users = null
  filesystem       = "swap"
  image            = null
  label            = "test-swap"
  linode_id        = 49667786
  root_pass        = null # sensitive
  size             = 1024
  stackscript_data = null # sensitive
  stackscript_id   = null
}

# __generated__ by Terraform from "49667786,52599108"
resource "linode_instance_config" "test_imported_config_1" {
  booted       = true
  comments     = null
  kernel       = "linode/grub2"
  label        = "test-config-1"
  linode_id    = 49667786
  memory_limit = 0
  root_device  = "/dev/sda"
  run_level    = "default"
  virt_mode    = "paravirt"
  device {
    device_name = "sda"
    disk_id     = 98757139
    volume_id   = 0
  }
  device {
    device_name = "sdb"
    disk_id     = 98757126
    volume_id   = 0
  }
  helpers {
    devtmpfs_automount = true
    distro             = true
    modules_dep        = true
    network            = true
    updatedb_disabled  = true
  }
  interface {
    ipam_address = null
    label        = null
    purpose      = "public"
  }
  interface {
    ipam_address = null
    label        = "test-vlan1"
    purpose      = "vlan"
  }
}
```

shell output

```
Terraform has generated configuration and written it to generated.tf. Please review the configuration and edit it as necessary before adding it to version control.
```

## Apply

Finally, execute `terraform apply` to incorporate the imported configurations into the state.

Don't forget to double check the reference to the sensitive values
because Terraform may not be able to generate them.

```
Apply complete! Resources: 3 imported, 0 added, 0 changed, 0 destroyed.
```

## Future

Ideally, we aim to import resources with IDs stored in the state.
Although this feature is not currently supported by Terraform, it
is planned for release in version 1.6.

We also plan to store Linode config ID in the embedded `config` field of `linode_instance`
resource.

Since then, we will be able to use dynamic ID values to import and generate
config files for the new resources.

For example:

```terraform
import {
  to = linode_instance_config.test_imported_config_1
  id = join(",", [linode_instance.test_instance.id, linode_instance.test_instance.config.id])
}

import {
  to = linode_instance_disk.test_imported_disk
  id = join(",", [linode_instance.test_instance.id, linode_instance.test_instance.disk.1.id])
}

import {
  to = linode_instance_disk.test_imported_disk_swap
  id = join(",", [linode_instance.test_instance.id, linode_instance.test_instance.disk.0.id])
}
```
