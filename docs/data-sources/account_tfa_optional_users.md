---
page_title: "Linode: linode_account_tfa_optional_users"
description: |-
  Lists all users exempt from account-level two-factor authentication enforcement.
---

# Data Source: linode\_account\_tfa\_optional\_users

Lists all users for whom two-factor authentication is optional, regardless of the account enforcement policy.
All result pages are retrieved. An account with no exemptions returns an empty set, not null.
Requires `account:read_only` scope and the `list_tfa_optional_users` permission.
See the [API documentation](https://techdocs.akamai.com/linode-api/reference/get-account-tfa-optional-users).

## Example Usage

```hcl
data "linode_account_tfa_optional_users" "example" {}

output "exempt_users" {
  value = data.linode_account_tfa_optional_users.example.usernames
}
```

## Argument Reference

There are no supported arguments.

## Attributes Reference

* `id` - The singleton identifier, `account`.

* `usernames` - (Set of String) The complete set of exempt usernames.
