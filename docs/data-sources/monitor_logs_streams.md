---
page_title: "Linode: linode_monitor_logs_streams"
description: |-
  Provides a list of Linode Monitor Logs Streams.
---

# linode\_monitor\_logs\_streams

Provides a list of Linode Monitor Logs Streams.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-logs-streams). (**Note: v4beta only.**)

## Example Usage

Retrieve all logs streams:

```terraform
data "linode_monitor_logs_streams" "all" {}
```

Retrieve active audit logs streams:

```terraform
data "linode_monitor_logs_streams" "active_audit" {
  filter {
    name   = "type"
    values = ["audit_logs"]
  }

  filter {
    name   = "status"
    values = ["active"]
  }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select logs streams that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Filterable Fields

* `label`
* `type`
* `status`

## Attributes Reference

Each logs stream will be stored in the `streams` attribute and will export the following attributes:

* `id` - The unique ID of this logs stream.

* `label` - The label for this logs stream.

* `type` - The type of this logs stream.

* `status` - The status of this logs stream.

* `created_by` - The user who created this logs stream.

* `updated_by` - The user who last updated this logs stream.

* `created` - When this logs stream was created.

* `updated` - When this logs stream was last updated.

* `version` - The version of this logs stream.
