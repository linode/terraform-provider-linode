---
layout: "linode"
page_title: "Linode: linode_account_logins"
sidebar_current: "docs-linode-datasource-account-logins"
description: |-
  Provides information about Linode account logins that match a set of filters.
---

# linode\_account\_logins

Provides information about Linode account logins that match a set of filters.

## Example Usage

The following example shows how one might use this data source to access information about a Linode account login.

```hcl
data "linode_account_logins" "filtered-account-logins" {
  filter {
    name = "restricted"
    values = ["true"]
  }

  filter {
    name = "username"
    values = ["myUsername"]
  }
}

output "login_ids" {
  value = data.linode_account_logins.filtered-account-logins.logins.*.id
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode account logins that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode account login will be stored in the `logins` attribute and will export the following attributes:

* `id` - The unique ID of this login object.

* `ip` - The remote IP address that requested the login.

* `datetime` - When the login was initiated.

* `username` - The username of the User that was logged into.

* `restricted` -  True if the User that was logged into was a restricted User, false otherwise.

## Filterable Fields

* `ip`

* `restricted`

* `username`
