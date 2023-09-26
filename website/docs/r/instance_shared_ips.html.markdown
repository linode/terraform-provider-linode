---
layout: "linode"
page_title: "Linode: linode_instance_shared_ips"
sidebar_current: "docs-linode-instance-shared-ips"
description: |-
  Manages IP addresses shared to a Linode.
---

# linode\_instance\_shared\_ips

~> **Beta Notice** IPv6 sharing is currently available through early access.
To use early access resources, the `api_version` provider argument must be set to `v4beta`.
To learn more, see the [early access documentation](../..#early-access).

~> **Notice** This resource should only be defined once per-instance and should not be used alongside the `shared_ipv4` field in `linode_instance`.

Manages IPs shared to a Linode instance.

## Example Usage

Share in IPv4 address between two instances:

```terraform
# Share the IP with the secondary node
resource "linode_instance_shared_ips" "share-primary" {
  linode_id = linode_instance.secondary.id
  addresses = [linode_instance_ip.primary.address]
}

# Allocate an IP under the primary node
resource "linode_instance_ip" "primary" {
  linode_id = linode_instance.primary.id
}

# Create a single primary node
resource "linode_instance" "primary" {
  label = "node-primary"
  type = "g6-nanode-1"
  region = "eu-central"
}

# Create a secondary node
resource "linode_instance" "secondary" {
  label = "node-secondary"
  type = "g6-nanode-1"
  region = "eu-central"
}
```

Share an IPv6 address among a primary node and its replicas:

```terraform
# Share with primary node
resource "linode_instance_shared_ips" "share-primary" {
  linode_id = linode_instance.primary.id
  addresses = [linode_ipv6_range.range.range]
}

# Share with secondary nodes
resource "linode_instance_shared_ips" "share-secondary" {
  count = var.number_replicas

  # Ranges must be shared with their primary node before being shared with a secondary
  depends_on = [linode_instance_shared_ips.share-primary]

  linode_id = linode_instance.secondary[count.index].id
  addresses = [linode_ipv6_range.range.range]
}

# Allocate an IPv6 range pointing at the primary node
resource "linode_ipv6_range" "range" {
  prefix_length = 64
  linode_id = linode_instance.primary.id
}

# Create a single primary node
resource "linode_instance" "primary" {
  label = "node-primary"
  type = "g6-nanode-1"
  region = "eu-central"
}

# Create two secondary nodes
resource "linode_instance" "secondary" {
  count = var.number_replicas

  label = "node-secondary-${count.index}"
  type = "g6-nanode-1"
  region = "eu-central"
}

variable "number_replicas" {
  default = 2
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The ID of the Linode to share the IPs to.

* `addresses` - (Required) The set of IPs to share with the Linode.
