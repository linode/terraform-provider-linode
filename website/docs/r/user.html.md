---
layout: "linode"
page_title: "Linode: linode_user"
sidebar_current: "docs-linode-resource-user"
description: |-
  Manages a Linode User.
---

# linode\_user

Manages a Linode User.

## Example Usage

```terraform
resource "linode_user" "john" {
    username = "john123"
    email = "john@acme.io"
    restricted = true
}
```

## Argument Reference

The following arguments are supported:

* `username` - (required) The username of the user.

* `email` - (required) The email address of the user.

* `restricted` - (optional) If true, this user will only have explicit permissions granted.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `tfa_enabled` - Whether the user has two-factor-authentication enabled.

* `ssh_keys` - A list of the User's SSH keys.
