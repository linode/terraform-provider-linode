---
page_title: "Linode: linode_child_account"
description: |-
  Provides details about a Linode Child Account.
---

# Data Source: linode\_child\_account

Provides information about a Linode Child Account.

Due to the sensitive nature of the data exposed by this data source, it should not be used in conjunction with the `LINODE_DEBUG` option.  See the [debugging notes](/providers/linode/linode/latest/docs#debugging) for more details.

## Example Usage

The following example shows how one might use this data source to access child account details.

```hcl
data "linode_child_account" "account" {
  euuid = "FFFFFFFF-FFFF-FFFF-FFFFFFFFFFFFFFFF"
}
```

## Argument Reference

The following arguments are supported:

* `euuid` - (Required) The unique EUUID of this Child Account.

## Attributes Reference

The Linode Account resource exports the following attributes:

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

* `active_since` - When this account was first activated.
