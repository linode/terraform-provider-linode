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

Create an unrestricted user:

```terraform
resource "linode_user" "john" {
    username = "john123"
    email = "john@acme.io"
}
```

Create a restricted user with grants:

```terraform
resource "linode_user" "fooser" {
    username = "cooluser123"
    email = "cool@acme.io"
    restricted = true

    global_grants {
        add_linodes = true
        add_images = true
    }

    linode_grant {
        id = 12345
        permissions = "read_write"
    }
}
```

## Argument Reference

The following arguments are supported:

* `username` - (required) The username of the user.

* `email` - (required) The email address of the user.

* `restricted` - (optional) If true, this user will only have explicit permissions granted.

* [`global_grants`](#global-grants) - (optional) A structure containing the Account-level grants a User has.

The following arguments are sets of [entity grants](#entity-grants):

* `domain_grant` - (optional) The domains the user has permissions access to.

* `image_grant` - (optional) The images the user has permissions access to.

* `linode_grant` - (optional) The Linodes the user has permissions access to.

* `longview_grant` - (optional) The longview the user has permissions access to.

* `nodebalancer_grant` - (optional) The NodeBalancers the user has permissions access to.

* `stackscript_grant` - (optional) The StackScripts the user has permissions access to.

* `volume_grant` - (optional) The volumes the user has permissions access to.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

* `tfa_enabled` - Whether the user has two-factor-authentication enabled.

* `ssh_keys` - A list of the User's SSH keys.

## Global Grants

* `account-access` - (optional) The level of access this User has to Account-level actions, like billing information. (`read_only`, `read_write`)

* `add_domains` - (optional) If true, this User may add Domains.

* `add_images` - (optional) If true, this User may add Images.

* `add_linodes` - (optional) If true, this User may create Linodes.

* `add_longview` - (optional) If true, this User may create Longview clients and view the current plan.

* `add_nodebalancers` - (optional) If true, this User may add NodeBalancers.

* `add_stackscripts` - (optional) If true, this User may add StackScripts.

* `cancel_account` - (optional) If true, this User may cancel the entire Account.

* `longview_subscription` - (optional) If true, this User may manage the Account’s Longview subscription.

## Entity Grants

* `id` - (required) The ID of the entity this grant applies to.

* `permissions` - (required) The level of access this User has to this entity. (`read_only`, `read_write`)
