---
page_title: "Linode: linode_monitor_logs_destination"
description: |-
  Provides details about a Linode Monitor Logs Destination.
---

# linode\_monitor\_logs\_destination

Provides details about a Linode Monitor Logs Destination.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-logs-destination). (**Note: v4beta only.**)

## Example Usage

```terraform
data "linode_monitor_logs_destination" "example" {
  id = 12345
}

output "destination_label" {
  value = data.linode_monitor_logs_destination.example.label
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of the logs destination.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label for this logs destination.

* `type` - The type of this logs destination. One of: `akamai_object_storage`, `custom_https`.

* `status` - The status of this logs destination.

* `created_by` - The user who created this logs destination.

* `updated_by` - The user who last updated this logs destination.

* `created` - When this logs destination was created.

* `updated` - When this logs destination was last updated.

* `version` - The version of this logs destination.

* [`details`](#details) - The details returned by the API. Write-only fields are not included.

### details

* `access_key_id` - The access key ID (applies to `akamai_object_storage` type).

* `bucket_name` - The bucket name (applies to `akamai_object_storage` type).

* `host` - The storage endpoint hostname (applies to `akamai_object_storage` type).

* `path` - The path within the bucket (applies to `akamai_object_storage` type).

* `endpoint_url` - The HTTPS endpoint URL (applies to `custom_https` type).

* `content_type` - The content type of log data. One of: `application/json`, `application/json; charset=utf-8` (applies to `custom_https` type).

* `data_compression` - The compression format (applies to `custom_https` type).

* `authentication_type` - The authentication type (applies to `custom_https` type).

* `tls_hostname` - The TLS hostname (applies to `custom_https` type).
