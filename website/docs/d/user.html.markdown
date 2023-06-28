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

* [`global_grants`](#global-grants) - The Account-level grants a User has.

* [`database_grant`](#global-grants) - The grants this User has pertaining to Databases on this Account.

* [`domain_grant`](#grant) - The grants this User has pertaining to Domains on this Account.

* [`firewall_grant`](#grant) - The grants this User has pertaining to Firewalls on this Account.

* [`image_grant`](#grant) - The grants this User has pertaining to Images on this Account.

* [`linode_grant`](#grant) - The grants this User has pertaining to Linodes on this Account.

* [`longview_grant`](#grant) - The grants this User has pertaining to Longview Clients on this Account.

* [`nodebalancer_grant`](#grant) - The grants this User has pertaining to NodeBalancers on this Account.

* [`stackscript_grant`](#grant) - The grants this User has pertaining to StackScripts on this Account.

* [`volume_grant`](#grant) - The grants this User has pertaining to Volumes on this Account.

* `id` - The unique identifier for this DataSource.

* `password_created` - The date and time when this User’s current password was created. User passwords are first created during the Account sign-up process, and updated using the Reset Password webpage. null if this User has not created a password yet.

* `tfa_enabled` - A boolean value indicating if the User has Two Factor Authentication (TFA) enabled.

* `verified_phone_number` - The phone number verified for this User Profile with the Phone Number Verify command. null if this User Profile has no verified phone number.

### Global Grants

* `account_access` - The level of access this User has to Account-level actions, like billing information. A restricted User will never be able to manage users. (`read_only`, `read_write`)

* `add_databases` - If true, this User may add Managed Databases.

* `add_domains` - If true, this User may add Domains.

* `add_firewalls` - If true, this User may add Firewalls.

* `add_images` - If true, this User may add Images.

* `add_linodes` - If true, this User may create Linodes.

* `add_longview` - If true, this User may create Longview clients and view the current plan.

* `add_nodebalancers` - If true, this User may add NodeBalancers.

* `add_stackscritps` - If true, this User may add StackScripts.

* `add_volumes` - If true, this User may add Volumes.

* `cancel_account` - If true, this User may cancel the entire Account.

* `longview_subscription` - If true, this User may manage the Account’s Longview subscription.

### Grant

* `id` - The ID of entity this grant applies to.

* `label` - The current label of the entity this grant applies to, for display purposes.

* `permissions` - The level of access this User has to this entity. If null, this User has no access. (`read_only`, `read_write`)
