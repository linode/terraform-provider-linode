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

## Argument Reference

The following arguments are supported:

* `username` - (Required) The unique username of this User.

## Attributes Reference

The Linode User resource exports the following attributes:

* `ssh_keys` - A list of SSH Key labels added by this User. These are the keys that will be deployed if this User is included in the authorized_users field of a create Linode, rebuild Linode, or create Disk request.

* `email` - The email address for this User, for account management communications, and may be used for other communications as configured.

* `restricted` - If true, this User must be granted access to perform actions or access entities on this Account.
