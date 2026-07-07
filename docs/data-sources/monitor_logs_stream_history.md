---
page_title: "Linode: linode_monitor_logs_stream_history"
description: |-
  Provides the historical versions of a Linode Monitor Logs Stream.
---

# linode\_monitor\_logs\_stream\_history

Provides the historical versions of a Linode Monitor Logs Stream.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-logs-stream-history). (**Note: v4beta only.**)

## Example Usage

```terraform
data "linode_monitor_logs_stream_history" "example" {
  stream_id = 12345
}

output "history_count" {
  value = length(data.linode_monitor_logs_stream_history.example.streams)
}
```

## Argument Reference

The following arguments are supported:

* `stream_id` - (Required) The ID of the logs stream to retrieve history for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `streams` - The historical versions of the logs stream. Each entry exports the following attributes:

  * `id` - The ID of the logs stream version.

  * `label` - The label of the logs stream at this version.

  * `type` - The type of the logs stream at this version. One of: `audit_logs`, `lke_audit_logs`.

  * `status` - The status of the logs stream at this version.

  * `version` - The version number of this history entry.

  * `destinations` - The destination IDs configured at this version.

  * `created` - The date and time when this version was created.

  * `updated` - The date and time when this version was last updated.

  * `created_by` - The user who created this stream version.

  * `updated_by` - The user who last updated this stream version.

  * [`details`](#details) - Additional configuration details at this version.

### details

* `cluster_ids` - The LKE cluster IDs included in this stream version.

* `is_auto_add_all_clusters_enabled` - Whether all clusters were auto-added at this version.
