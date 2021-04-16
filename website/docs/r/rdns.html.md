---
layout: "linode"
page_title: "Linode: linode_rdns"
sidebar_current: "docs-linode-resource-rdns"
description: |-
  Manages the RDNS / PTR record for the IP Address associated with a Linode Instance.
---

# linode\_rdns

Provides a Linode RDNS resource.  This can be used to create and modify RDNS records.

Linode RDNS names must have a matching address value in an A or AAAA record.  This A or AAAA name must be resolvable at the time the RDNS resource is being associated.

For more information, see the [Linode APIv4 docs](https://developers.linode.com/api/v4/networking-ips-address/#put) and the [Configure your Linode for Reverse DNS](https://www.linode.com/docs/networking/dns/configure-your-linode-for-reverse-dns-classic-manager/) guide.

## Example Usage

The following example shows how one might use this resource to configure an RDNS address for an IP address.

```hcl
resource "linode_rdns" "foo" {
  address = linode_instance.foo.ip_address
  rdns = "${linode_instance.foo.ip_address}.nip.io"
}

resource "linode_instance" "foo" {
   image = "linode/alpine3.9"
   region = "ca-east"
   type = "g6-dedicated-2"
}
```

The following example shows how one might use this resource to configure RDNS for multiple IP addresses.

```hcl
resource "linode_instance" "my_instance" {
  count = 3

  label = "simple_instance-${count.index + 1}"
  image = "linode/ubuntu18.04"
  region = "us-central"
  type = "g6-standard-1"
  root_pass = "terr4form-test"
}

resource "linode_rdns" "my_rdns" {
  count = length(linode_instance.my_instance)

  address = linode_instance.my_instance[count.index].ip_address
  rdns = "${linode_instance.my_instance[count.index].ip_address}.nip.io"
}
```

## Argument Reference

The following arguments are supported:

* `address` - The Public IPv4 or IPv6 address that will receive the `PTR` record.  A matching `A` or `AAAA` record must exist.

* `rdns` - The name of the RDNS address.

## Import

Linodes RDNS resources can be imported using the address as the `id`.

```sh
terraform import linode_rdns.foo 123.123.123.123
```
