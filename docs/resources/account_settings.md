---
page_title: "Linode: linode_account_settings"
description: |-
  Manages the settings of a Linode account.
---

# linode\_account\_settings

Manages the settings of a Linode account.

## Example Usage

The following example shows how one might use this resource to change their Linode account settings.

```hcl
resource "linode_account_settings" "myaccount" {
    longview_subscription = "longview-40"
    backups_enabled = "true"
}
```

## Argument Reference

The following arguments are supported:

* `backups_enabled` - (Optional) The account-wide backups default. If true, all Linodes created will automatically be enrolled in the Backups service. If false, Linodes will not be enrolled by default, but may still be enrolled on creation or later.

* `network_helper` - (Optional) Enables network helper across all users by default for new Linodes and Linode Configs.

* `longview_subscription` - (Optional) The Longview Pro tier you are currently subscribed to. The value must be a [Longview Subscription](https://www.linode.com/docs/api/longview/#longview-subscriptions-list) ID or null for Longview Free.

## Additional Results

* `managed` - Enables monitoring for connectivity, response, and total request time.

* `object_storage` - A string describing the status of this account’s Object Storage service enrollment.
