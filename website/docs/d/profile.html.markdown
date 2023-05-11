---
layout: "linode"
page_title: "Linode: linode_profile"
sidebar_current: "docs-linode-datasource-profile"
description: |-
  Provides details about a Linode profile.
---

# Data Source: linode\_profile

Provides information about a Linode profile.

## Example Usage

The following example shows how one might use this data source to access profile details.

```hcl
data "linode_profile" "profile" {}
```

## Argument Reference

There are no supported arguments because the provider `token` can only access the associated profile.

## Attributes Reference

The Linode Profile resource exports the following attributes:

* `email` - The profile email address. This address will be used for communication with Linode as necessary.

* `timezone` - The profile's preferred timezone. This is not used by the API, and is for the benefit of clients only. All times the API returns are in UTC.

* `email_notifications` - If true, email notifications will be sent about account activity. If false, when false business-critical communications may still be sent through email.

* `username` - The username for logging in to Linode services.

* `ip_whitelist_enabled` - If true, logins for the user will only be allowed from whitelisted IPs. This setting is currently deprecated, and cannot be enabled.

* `lish_auth_method` - The methods of authentication allowed when connecting via Lish. 'keys_only' is the most secure with the intent to use Lish, and 'disabled' is recommended for users that will not use Lish at all.

* `authorized_keys` - The list of SSH Keys authorized to use Lish for this user. This value is ignored if lish_auth_method is 'disabled'.

* `two_factor_auth` - If true, logins from untrusted computers will require Two Factor Authentication.

* `restricted` - If true, the user has restrictions on what can be accessed on the Account.

* `referrals` - Credit Card information associated with this Account.

* `referrals.total` - The number of users who have signed up with the referral code.

* `referrals.credit` - The amount of account credit in US Dollars issued to the account through the referral program.

* `referrals.completed` - The number of completed signups with the referral code.

* `referrals.pending` - The number of pending signups for the referral code. To receive credit the signups must be completed.

* `referrals.code` - The Profile referral code.  If new accounts use this when signing up for Linode, referring account will receive credit.

* `referrals.url` - The referral URL.
