---
page_title: "Linode: linode_monitor_alert_channels"
description: |-
  Provides information about Linode Monitor Alert notification channels.
---

# Data Source: linode_monitor_alert_channels

Use this data source to query Linode Monitor Alert notification channels.

## Example Usage

```hcl
data "linode_monitor_alert_channels" "all" {}

output "first_channel_label" {
  value = data.linode_monitor_alert_channels.all.monitor_alert_channels[0].label
}
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

- `filter` - (Optional) One or more filter blocks to select channels.
  - `name` - (Required) The name of the field to filter by.

## Attributes Reference

Each Alert Channel will be stored in the `monitor_alert_channels` attribute and will export the following attributes:

- `id` - The unique ID of the channel.
- `label` - The display label for the channel.
- `type` - The channel type (`system` or `user`).
- `channel_type` - The notification transport type (currently `email`).
- `created` - When the channel was created.
- `updated` - When the channel was last updated.
- `created_by` - The user who created the channel, or `system`.
- `updated_by` - The user who last updated the channel, or `system`.
- `alerts` - Alert linkage metadata for this channel.
  - `url` - The API URL for associated alerts.
  - `type` - The alert type associated with the channel.
  - `alert_count` - The number of associated alerts.
- `content` - (Deprecated) Legacy read-only channel content.
  - `email` - Legacy email content values.
    - `email_addresses` - Legacy email recipients for system channels.
- `details` - Channel configuration details.
  - `email` - Email-specific configuration details.
    - `usernames` - Usernames that receive notifications.
    - `recipient_type` - Recipient selection mode (for example `read_write_users` or `user`).

## Filterable Fields

The following top-level fields can be used with the `filter` block:

  - `id`
  - `label`
  - `type`
  - `channel_type`