---
layout: "linode"
page_title: "Linode: linode_user"
sidebar_current: "docs-linode-datasource-user"
description: |-
  Provides details about a Linode user.
---

# Data Source: linode\_user

Provides information about a Linode user

## Example Usage

The following example shows how one might use this data source to access information about a Linode user.

```hcl
data "linode_user" "foo" {
    username = "foo"
}
```

The following example shows a sample grant.

```hcl
"domain": [
  {
    "id": 123,
    "label": "example-entity",
    "permissions": "read_only"
  }
]
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The unique username of this User.

## Attributes Reference

The Linode User resource exports the following attributes:

* `ssh_keys` - A list of SSH Key labels added by this User. These are the keys that will be deployed if this User is included in the authorized_users field of a create Linode, rebuild Linode, or create Disk request.

* `email` - The email address for this User, for account management communications, and may be used for other communications as configured.

* `restricted` - If true, this User must be granted access to perform actions or access entities on this Account.

* `global_grants` - The Account-level grants a User has.

* `database_grant` - The grants this User has pertaining to Databases on this Account.

* `domain_grant` - The grants this User has pertaining to Domains on this Account.

* `firewall_grant` - The grants this User has pertaining to Firewalls on this Account.

* `image_grant` - The grants this User has pertaining to Images on this Account.

* `linode_grant` - The grants this User has pertaining to Linodes on this Account.

* `longview_grant` - The grants this User has pertaining to Longview Clients on this Account.

* `nodebalancer_grant` - The grants this User has pertaining to NodeBalancers on this Account.

* `stackscript_grant` - The grants this User has pertaining to StackScripts on this Account.

* `volume_grant` - The grants this User has pertaining to Volumes on this Account.
