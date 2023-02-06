---
layout: "linode"
page_title: "Linode: linode_account_settings"
sidebar_current: "docs-linode-datasource-account-settings"
description: |-
Provides information about Linode account settings.
---

# linode\_account\_settings

Provides information about Linode account settings.

## Example Usage

The following example shows how one might use this data source to access information about Linode account settings.

```hcl
data "linode_account_settings" "example" {}
```

## Attributes Reference

* `backups_enabled` - Account-wide backups default.

* `longview_subscription` - The Longview Pro tier you are currently subscribed to.

* `managed` - Enables monitoring for connectivity, response, and total request time.

* `network_helper` - Enables network helper across all users by default for new Linodes and Linode Configs.

* `object_storage` - A string describing the status of this account’s Object Storage service enrollment.
