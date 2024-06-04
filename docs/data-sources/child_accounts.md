---
page_title: "Linode: linode_child_accounts"
description: |-
  Provides information about Linode account logins that match a set of filters.
---

# linode\_child\_accounts

Provides information about Linode Child Accounts that match a set of filters.

## Example Usage

The following example shows how one might use this data source to access Child Accounts under the current Account.

```hcl
data "linode_child_accounts" "all" {}

data "linode_child_accounts" "filtered" {
  filter {
    name = "email"
    values = ["example@linode.com"]
  }

  filter {
    name = "first_name"
    values = ["John"]
  }

  filter {
    name = "last_name"
    values = ["Smith"]
  }
}

output "all_accounts" {
  value = data.linode_child_accounts.all.child_accounts.*.euuid
}

output "filtered_accounts" {
  value = data.linode_child_accounts.filtered.child_accounts.*.euuid
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Child Accounts that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Child Account will be stored in the `child_accounts` attribute and will export the following attributes:

* `email` - The email address for this Account, for account management communications, and may be used for other communications as configured.

* `first_name` - The first name of the person associated with this Account.

* `last_name` - The last name of the person associated with this Account.

* `company` - The company name associated with this Account.

* `address_1` - First line of this Account's billing address.

* `address_2` - Second line of this Account's billing address.

* `phone` - The phone number associated with this Account.

* `city` - The city for this Account's billing address.

* `state` - If billing address is in the United States, this is the State portion of the Account's billing address. If the address is outside the US, this is the Province associated with the Account's billing address.

* `country` - The two-letter country code of this Account's billing address.

* `zip` - The zip code of this Account's billing address.

* `balance` - This Account's balance, in US dollars.

* `capabilities` - A set containing all the capabilities of this Account.

* `active_since` - When this account was first activated

## Filterable Fields

* `euuid`

* `email`

* `first_name`

* `last_name`

* `company`

* `address_1`

* `address_2`

* `phone`

* `city`

* `state`

* `country`

* `zip`

* `capabilities`

* `active_since`
