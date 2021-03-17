---
layout: "linode"
page_title: "Linode: linode_sshkey"
sidebar_current: "docs-linode-resource-sshkey"
description: |-
  Manages a Linode SSH Key.
---

# linode\_sshkey

Provides a Linode SSH Key resource.  This can be used to create, modify, and delete Linodes SSH Keys.  Managed SSH Keys allow instances to be created with a list of Linode usernames, whose SSH keys will be automatically applied to the root account's `~/.ssh/authorized_keys` file.
For more information, see the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/getSSHKeys).

## Example Usage

The following example shows how one might use this resource to configure a SSH Key for access to a Linode Instance.

```hcl
resource "linode_sshkey" "foo" {
  label = "foo"
  ssh_key = chomp(file("~/.ssh/id_rsa.pub"))
}

resource "linode_instance" "foo" {
  image  = "linode/ubuntu18.04"
  label  = "foo"
  region = "us-east"
  type   = "g6-nanode-1"
  authorized_keys    = [linode_sshkey.foo.ssh_key]
  root_pass      = "..."
}
```

## Argument Reference

The following arguments are supported:

* `label` - A label for the SSH Key.

* `ssh_key` - The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.

## Attributes

This resource exports the following attributes:

* `created` - The date this SSH Key was created.

## Import

Linodes SSH Keys can be imported using the Linode SSH Key `id`, e.g.

```sh
terraform import linode_sshkey.mysshkey 1234567
```
