---
layout: "linode"
page_title: "Linode: linode_account"
sidebar_current: "docs-linode-datasource-account"
description: |-
  Provides details about a Linode account.
---

# Data Source: linode\_account

Provides information about a Linode account.

This data source should not be used in conjuction with the `LINODE_DEBUG` option.  See the [debugging notes](/providers/linode/linode/latest/docs#debugging) for more details.

## Example Usage

The following example shows how one might use this data source to access account details.

```hcl
data "linode_account" "account" {}
```

## Argument Reference

There are no supported arguments because the provider `token` can only access the associated account.

## Attributes

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
