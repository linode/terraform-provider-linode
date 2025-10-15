---
page_title: "Linode: linode_interface"
description: |-
  Manages a Linode interface configuration.
---

# linode\_interface

Provides a Linode Interface resource that can be used to create, modify, and delete network interfaces for Linode instances. Interfaces allow you to configure public, VLAN, and VPC networking for your Linode instances.

This resource is specifically for Linode interfaces. If you are interested in deploying a Linode instance with a legacy config interface, please refer to the `linode_instance_config` resource documentation for details.

This resource is designed to work with explicitly defined disk and config resources for the Linode instance. See the [Complete Example with Linode](#complete-example-with-linode) section below for details.

For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-linode-instance-interface).

## Example Usage

### Public Interface Example

The following example shows how to create a public interface with specific IPv4 and IPv6 configurations.

```hcl
resource "linode_interface" "public" {
  linode_id = linode_instance.my-instance.id

  public = {
    ipv4 = {
      addresses = [
        {
          address = "auto",
          primary = true,
        }
      ]
    }
    ipv6 = {
      ranges = [
        {
          range = "/64"
        }
      ]
    }
  }
}
```

### IPv6-Only Public Interface Example

The following example shows how to create an IPv6-only public interface. Note that you must explicitly set `addresses = []` to prevent the automatic creation of an IPv4 address.

```hcl
resource "linode_interface" "ipv6_only" {
  linode_id = linode_instance.my-instance.id

  public = {
    ipv4 = {
      addresses = []  # Empty list prevents auto-creation of IPv4 address
    }
    ipv6 = {
      ranges = [
        {
          range = "/64"
        }
      ]
    }
  }
}
```

### VPC Interface Example

The following example shows how to create a VPC interface with custom IPv4 configuration and 1:1 NAT.

```hcl
resource "linode_interface" "vpc" {
  linode_id   = linode_instance.my-instance.id

  vpc = {
    subnet_id = 240213
    ipv4 = {
      addresses = [
        {
          address = "auto"
        }
      ]
      ranges = [
        {
          range = "/32"
        }
      ]
    }
  }
}
```

### VLAN Interface Example

The following example shows how to create a VLAN interface.

```hcl
resource "linode_interface" "vlan" {
  linode_id = linode_instance.web.id

  vlan = {
    vlan_label   = "web-vlan"
    ipam_address = "192.168.200.5/24"
  }
}
```

### Complete Example with Linode

```hcl
resource "linode_instance" "my-instance" {
  label                = "my-instance"
  region               = "us-mia"
  type                 = "g6-standard-1"
  interface_generation = "linode"
}

resource "linode_instance_config" "my-config" {

  # This is necessary to ensure the interface is created
  # before the config is booted with the Linode instance
  depends_on = [linode_interface.public]

  linode_id = linode_instance.my-instance.id
  label     = "my-config"

  device {
    device_name = "sda"
    disk_id     = linode_instance_disk.boot.id
  }

  booted = true
}

resource "linode_instance_disk" "boot" {
  label     = "boot"
  linode_id = linode_instance.my-instance.id
  size      = linode_instance.my-instance.specs.0.disk

  image     = "linode/debian12"
  root_pass = "this-is-NOT-a-safe-password"
}

resource "linode_interface" "public" {
  linode_id = linode_instance.my-instance.id
  public = {
    ipv4 = {
      addresses = [
        {
          address = "auto",
          primary = true,
        }
      ]
    }
    ipv6 = {
      ranges = [
        {
          range = "/64"
        }
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to assign this interface to.

* `firewall_id` - (Optional) The ID of an enabled firewall to secure a VPC or public interface. Not allowed for VLAN interfaces.

* `default_route` - (Optional) Indicates whether the interface serves as the default route when multiple interfaces are eligible for this role.

  * `ipv4` - (Optional) When set to true, the interface is used for the IPv4 default route.

  * `ipv6` - (Optional) When set to true, the interface is used for the IPv6 default route.

* `public` - (Optional) Configuration for a Linode public interface. Exactly one of `public`, `vlan`, or `vpc` must be specified.

  * `ipv4` - (Optional) IPv4 configuration for this interface.

    * `addresses` - (Optional) IPv4 addresses configured for this Linode interface. Each object in this list supports:

      * `address` - (Optional) The IPv4 address. Defaults to "auto" for automatic assignment.

      * `primary` - (Optional) Whether this address is the primary address for the interface.

  * `ipv6` - (Optional) IPv6 configuration for this interface.

    * `ranges` - (Optional) IPv6 ranges in CIDR notation (2600:0db8::1/64) or prefix-only (/64). Each object in this list supports:

      * `range` - (Required) The IPv6 range.

      * `route_target` - (Optional) The public IPv6 address that the range is routed to.

* `vlan` - (Optional) Nested attributes object for a Linode VLAN interface. Exactly one of `public`, `vlan`, or `vpc` must be specified.

  * `ipam_address` - (Optional) The VLAN interface's private IPv4 address in CIDR notation.

  * `vlan_label` - (Required) The VLAN's unique label. Must be between 1 and 64 characters.

* `vpc` - (Optional) Nested attributes object for a Linode VPC interface. Exactly one of `public`, `vlan`, or `vpc` must be specified.

  * `subnet_id` - (Required) The VPC subnet identifier for this interface.

  * `ipv4` - (Optional) IPv4 configuration for the VPC interface.

    * `addresses` - (Optional) Specifies the IPv4 addresses to use in the VPC subnet. Each object in this list supports:

      * `address` - (Optional) The IPv4 address. Defaults to "auto" for automatic assignment.

      * `primary` - (Optional) Whether this address is the primary address for the interface.

      * `nat_1_1_address` - (Optional) The 1:1 NAT IPv4 address used to associate a public IPv4 address with the interface's VPC subnet IPv4 address.

    * `ranges` - (Optional) IPv4 ranges in CIDR notation (1.2.3.4/24) or prefix-only format (/24). Each object in this list supports:

      * `range` - (Required) The IPv4 range.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for this interface.

* `public` - When a public interface is configured, the following computed attributes are available:

  * `ipv4` - IPv4 configuration for the public interface:

    * `assigned_addresses` - (Computed) The IPv4 addresses exclusively assigned to this Linode interface. Each object in this set supports:

      * `address` - The assigned IPv4 address.

      * `primary` - Whether this address is the primary address for the interface.

    * `shared` - (Computed) The IPv4 addresses assigned to this Linode interface that are also shared with another Linode. Each object in this set supports:

      * `address` - The shared IPv4 address.

      * `linode_id` - The ID of the Linode that this address is shared with.

  * `ipv6` - IPv6 configuration for the public interface:

    * `assigned_ranges` - (Computed) The IPv6 ranges exclusively assigned to this Linode interface. Each object in this set supports:

      * `range` - The assigned IPv6 range.

      * `route_target` - The public IPv6 address that the range is routed to.

    * `shared` - (Computed) The IPv6 ranges assigned to this Linode interface that are also shared with another Linode. Each object in this set supports:

      * `range` - The shared IPv6 range.

      * `route_target` - The public IPv6 address that the range is routed to.

    * `slaac` - (Computed) The public SLAAC and subnet prefix settings for this public interface. Each object in this set supports:

      * `address` - The SLAAC IPv6 address.

      * `prefix` - The subnet prefix length.

* `vpc` - When a VPC interface is configured, the following computed attributes are available:

  * `ipv4` - IPv4 configuration for the VPC interface:

    * `assigned_addresses` - (Computed) The IPv4 addresses assigned for use in the VPC subnet, calculated from the `addresses` input. Each object in this set supports:

      * `address` - The assigned IPv4 address.

      * `primary` - Whether this address is the primary address for the interface.

      * `nat_1_1_address` - The assigned 1:1 NAT IPv4 address used to associate a public IPv4 address with the interface's VPC subnet IPv4 address.

    * `assigned_ranges` - (Computed) The IPv4 ranges assigned for use in the VPC subnet, calculated from the `ranges` input. Each object in this set supports:

      * `range` - The assigned IPv4 range.

## Import

Interfaces can be imported using a Linode ID followed by an Interface ID, separated by a comma, e.g.

```sh
terraform import linode_interface.example 67890,12345
```

## Notes

* Each Linode instance can have up to 3 network interfaces.
* VLAN interfaces cannot be updated after creation and require recreation.
* VPC subnet IDs cannot be changed after interface creation.
* Firewall IDs are only supported for public and VPC interfaces, not for VLAN interfaces.
* When configuring multiple interfaces, use the `default_route` setting to specify which interface should handle default routing.
