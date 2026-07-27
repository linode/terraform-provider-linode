---
page_title: "Linode: linode_monitor_logs_destinations"
description: |-
  Provides a list of Linode Monitor Logs Destinations.
---

# linode\_monitor\_logs\_destinations

Provides a list of Linode Monitor Logs Destinations.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-logs-destinations). (**Note: v4beta only.**)

## Example Usage

Retrieve all logs destinations:

```terraform
data "linode_monitor_logs_destinations" "all" {}
```

Retrieve logs destinations filtered by type:

```terraform
data "linode_monitor_logs_destinations" "object_storage" {
  filter {
    name   = "type"
    values = ["akamai_object_storage"]
  }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select logs destinations that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Filterable Fields

* `label`
* `type`
* `status`

## Attributes Reference

Each logs destination will be stored in the `destinations` attribute and will export the following attributes:

* `id` - The unique ID of this logs destination.

* `label` - The label for this logs destination.

* `type` - The type of this logs destination.

* `status` - The status of this logs destination.

* `created_by` - The user who created this logs destination.

* `updated_by` - The user who last updated this logs destination.

* `created` - When this logs destination was created.

* `updated` - When this logs destination was last updated.

* `version` - The version of this logs destination.
