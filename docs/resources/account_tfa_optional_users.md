---
page_title: "Linode: linode_account_tfa_optional_users"
description: |-
  Manages the complete account-wide set of users exempt from two-factor authentication enforcement.
---

# linode\_account\_tfa\_optional\_users

Manages the **complete** set of users for whom two-factor authentication (2FA) is optional.
Only one resource should manage this set per account, including across Terraform states.
Applying this resource replaces all existing exemptions, including those managed outside Terraform.
See the [API documentation](https://techdocs.akamai.com/linode-api/reference/put-tfa-optional-users).

Exemptions persist regardless of the account's `tfa_enforced` policy, but only affect access when enforcement is enabled.
Removing an exemption immediately restores the account policy and can prevent users without 2FA configured from logging in.

## Example Usage

```hcl
resource "linode_account_tfa_optional_users" "example" {
  usernames = ["automation-user"]
}

resource "linode_account_settings" "example" {
  tfa_enforced = true
  depends_on   = [linode_account_tfa_optional_users.example]
}
```

The explicit dependency sets exemptions before enabling enforcement.

## Argument Reference

* `usernames` - (Required, Set of String) The complete set of exempt usernames. Set to `[]` to remove all exemptions. Ordering is ignored.

## Attributes Reference

* `id` - The singleton identifier, `account`.

## Lifecycle

Create and update replace the entire list, then retrieve all pages of the saved list.
The PUT response contains only a count and is not used as the list in state.
Refresh detects external additions, removals, and deleted or renamed users; the next apply restores the configured set.
If a configured user no longer exists, update the configuration before applying.
An empty exemption set remains a managed resource.

The resource's delete action makes no API requests and removes it from Terraform state without changing the exemptions in Linode.
Terraform may still refresh the resource before a destroy.
To remove all exemptions, first set `usernames = []` and apply the change before destroying the resource.
Removing exemptions restores the account policy and may prevent users without 2FA configured from logging in.

The provider needs `account:read_write` scope and the `list_tfa_optional_users` and `update_tfa_optional_users` permissions.
If a write succeeds but its follow-up read fails, Terraform reports an error and retains the resource for recovery on refresh.

## Import

Import the existing account-wide set using `account`. Import does not change exemptions.

```sh
terraform import linode_account_tfa_optional_users.example account
```
