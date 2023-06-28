---
layout: "linode"
page_title: "Linode: linode_users"
sidebar_current: "docs-linode-datasource-users"
description: |-
Lists Users on your Account.

Users may access all or part of your Account based on their restricted status and grants. An unrestricted User may access everything on the account, whereas restricted User may only access entities or perform actions they’ve been given specific grants to.
---

# linode\_users

Provides information about Linode users that match a set of filters.

```hcl
data "linode_users" "filtered-users" {
  filter {
    name = "username"
    values = ["test-user"]
  }
}

output "users" {
  value = data.linode_users.filtered-users.users
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode users that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode user will be stored in the `users` attribute and will export the following attributes:

* `username` - This User's username. This is used for logging in, and may also be displayed alongside actions the User performs (for example, in Events or public StackScripts).

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

* `permissions` - The level of access this User has to this entity. If null, this User has no access.

## Filterable Fields

* `username`

* `email`

* `restricted`

* `password_created`

* `tfa_enabled`

* `verfied_phone_number`
