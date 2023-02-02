---
layout: "linode"
page_title: "Linode: linode_account_login"
sidebar_current: "docs-linode-datasource-account-login"
description: |-
  Provides details about a Linode account login.
---

# linode\_account\_login

Provides details about a specific Linode account login.

## Example Usage

The following example shows how one might use this data source to access information about a Linode account login.

```hcl
data "linode_account_login" "my_account_login" {
    id = 123456
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this login object.

## Attributes Reference

The Linode Account Login resource exports the following attributes:

* `id` - The unique ID of this login object.

* `ip` - The remote IP address that requested the login.

* `datetime` - When the login was initiated.

* `username` - The username of the User that was logged into.

* `restricted` -  True if the User that was logged into was a restricted User, false otherwise.
