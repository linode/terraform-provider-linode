---
page_title: "Linode: linode_iam_user"
description: |-
  Provides IAM details about a Linode user.
---

# Data Source: linode\_iam\_user

Provides IAM information about a Linode user
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-iam-users-role-permissions).

## Example Usage

The following example shows how one might use this data source to access IAM information about a Linode user.

```hcl
data "linode_iam_user" "foo" {
    username = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The unique username of this User.

## Attributes Reference

The Linode IAM User datasource exports the following attributes:

* `account_access` - A list of account level roles the user currently has.

* [`entity_access`](#entity-access) - A list of specific entities the user has specific roles for.

### Entity Access

* `id` - The unique ID for the entity.

* `type` - The type of product for the entity. (eg. Volume)

* `roles` - A list of the roles for this entity and specific user.
