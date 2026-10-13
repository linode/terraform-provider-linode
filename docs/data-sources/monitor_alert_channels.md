---
page_title: "Linode: linode_monitor_alert_channels"
description: |-
  Provides information about Linode Monitor Alert notification channels.
---

# Data Source: linode_monitor_alert_channels

Use this data source to query Linode Monitor Alert notification channels.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-notification-channels).  (**Note: v4beta only.**)

## Example Usage

Retrieve ALL Monitor Alert Channels:

```hcl
data "linode_monitor_alert_channels" "all" {}
```

## Filter Example

```hcl
data "linode_monitor_alert_channels" "system" {
  filter {
    name   = "type"
    values = ["system"]
  }
}
```

## Arguments Reference

The following arguments are supported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`filter`](#filter) - (Optional, Block Set) A set of filters used to select monitor alert channels that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.
* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.
* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Alert Channel will be stored in the `monitor_alert_channels` attribute and will export the following attributes:

* `monitor_alert_channels` - (Nested Attribute List) Monitor alert channels matching the query.

* `id` - The unique ID of the channel.
* `label` - The display label for the channel.
* `type` - The channel type (`system` or `user`).
* `channel_type` - The notification transport type (currently `email`).
* `created` - When the channel was created.
* `updated` - When the channel was last updated.
* `created_by` - The user who created the channel, or `system`.
* `updated_by` - The user who last updated the channel, or `system`.
* `alerts` - (Nested Attribute) Alert linkage metadata for this channel. Referenced directly (e.g. `alerts.alert_count`).
  * `url` - The API URL for associated alerts.
  * `type` - The alert type associated with the channel.
  * `alert_count` - The number of associated alerts.
* `details` - (Nested Attribute) Channel configuration details. Referenced directly (e.g. `details.email`).
  * `email` - (Nested Attribute) Email-specific configuration details. Referenced directly (e.g. `email.recipient_type`).
    * `usernames` - Usernames that receive notifications.
    * `recipient_type` - Recipient selection mode (for example `read_write_users` or `user`).

## Filterable Fields

The following top-level fields can be used with the `filter` block:

* `id`
* `label`
* `type`
* `channel_type`
