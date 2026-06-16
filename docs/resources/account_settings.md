---
page_title: "Linode: linode_account_settings"
description: |-
  Manages the settings of a Linode account.
---

# linode\_account\_settings

Manages the settings of a Linode account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-account-settings).

## Example Usage

The following example shows how one might use this resource to change their Linode account settings.

```hcl
resource "linode_account_settings" "myaccount" {
    backups_enabled = "true"
}
```

## Argument Reference

The following arguments are supported:

* `backups_enabled` - (Optional) The account-wide backups default. If true, all Linodes created will automatically be enrolled in the Backups service. If false, Linodes will not be enrolled by default, but may still be enrolled on creation or later.

* `network_helper` - (Optional) Enables network helper across all users by default for new Linodes and Linode Configs.

* `interfaces_for_new_linodes` - (Optional) Type of interfaces for new Linode instances. Available values are `"legacy_config_only"`, `"legacy_config_default_but_linode_allowed"`, `"linode_default_but_legacy_config_allowed"`, and `"linode_only"`.

* `maintenance_policy` - (Optional) The default maintenance policy for this account. Examples are `"linode/migrate"` and `"linode/power_off_on"`. Defaults to `"linode/migrate"`.

## Additional Results

* `managed` - Enables monitoring for connectivity, response, and total request time.

* `longview_subscription` - (Deprecated) The Longview Pro tier you are currently subscribed to.

* `object_storage` - A string describing the status of this account’s Object Storage service enrollment.
