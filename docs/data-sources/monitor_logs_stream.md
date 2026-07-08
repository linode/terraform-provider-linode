---
page_title: "Linode: linode_monitor_logs_stream"
description: |-
  Provides details about a Linode Monitor Logs Stream.
---

# linode\_monitor\_logs\_stream

Provides details about a Linode Monitor Logs Stream.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-logs-stream). (**Note: v4beta only.**)

## Example Usage

```terraform
data "linode_monitor_logs_stream" "example" {
  id = "12345"
}

output "stream_label" {
  value = data.linode_monitor_logs_stream.example.label
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the logs stream.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the logs stream.

* `type` - The type of the logs stream. One of: `audit_logs`, `lke_audit_logs`.

* `status` - The status of the logs stream.

* `version` - The version of the logs stream.

* `destinations` - The list of logs destination IDs attached to this stream.

* `created` - The date and time when the logs stream was created.

* `updated` - The date and time when the logs stream was last updated.

* `created_by` - The user who created the logs stream.

* `updated_by` - The user who last updated the logs stream.

* [`details`](#details) - Additional configuration details. Only populated for `lke_audit_logs` streams.

### details

* `cluster_ids` - The list of LKE cluster IDs included in this stream.

* `is_auto_add_all_clusters_enabled` - When true, all LKE clusters are automatically added to this stream.
