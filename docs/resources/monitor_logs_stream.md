---
page_title: "Linode: linode_monitor_logs_stream"
description: |-
  Manages a Linode Monitor Logs Stream.
---

# linode\_monitor\_logs\_stream

Manages a Linode Monitor Logs Stream.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-logs-stream). (**Note: v4beta only.**)

## Example Usage

Creating a basic audit logs stream:

```terraform
resource "linode_object_storage_key" "key" {
  label   = "my-logs-key"
  regions = ["us-east"]
}

resource "linode_object_storage_bucket" "bucket" {
  region     = "us-east"
  label      = "my-logs-bucket"
  access_key = linode_object_storage_key.key.access_key
  secret_key = linode_object_storage_key.key.secret_key
}

resource "linode_monitor_logs_destination" "destination" {
  label = "my-logs-destination"
  type  = "akamai_object_storage"

  akamai_object_storage_details = {
    access_key_id     = linode_object_storage_key.key.access_key
    access_key_secret = linode_object_storage_key.key.secret_key
    bucket_name       = linode_object_storage_bucket.bucket.label
    host              = linode_object_storage_bucket.bucket.hostname
  }
}

resource "linode_monitor_logs_stream" "example" {
  label        = "my-audit-logs-stream"
  type         = "audit_logs"
  destinations = [linode_monitor_logs_destination.destination.id]
}
```

Creating an LKE audit logs stream with specific cluster IDs:

```terraform
resource "linode_monitor_logs_stream" "lke_example" {
  label        = "my-lke-stream"
  type         = "lke_audit_logs"
  destinations = [linode_monitor_logs_destination.destination.id]

  details = {
    cluster_ids                    = [12345, 67890]
    is_auto_add_all_clusters_enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the logs stream.

* `type` - (Required, Forces New) The type of the logs stream. One of: `audit_logs`, `lke_audit_logs`.

* `destinations` - (Required) The list of logs destination IDs attached to this stream.

* `status` - (Optional) The status of the logs stream. One of: `active`, `inactive`, `provisioning`, `deactivating`.

* [`details`](#details) - (Optional) Additional configuration details. Only applies to `lke_audit_logs` streams.

### details

* `cluster_ids` - (Optional) The list of LKE cluster IDs to include in this stream.

* `is_auto_add_all_clusters_enabled` - (Optional) When true, all LKE clusters are automatically added to this stream.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the logs stream.

* `version` - The version of the logs stream.

* `created` - The date and time when the logs stream was created.

* `updated` - The date and time when the logs stream was last updated.

* `created_by` - The user who created the logs stream.

* `updated_by` - The user who last updated the logs stream.

## Import

Monitor Logs Streams can be imported using the `id`, e.g.

```sh
terraform import linode_monitor_logs_stream.example 12345
```
