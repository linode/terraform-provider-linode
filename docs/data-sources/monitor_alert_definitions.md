---
page_title: "Linode: linode_monitor_alert_definitions"
description: |-
Retrieves Monitor Alert Definitions.
---

# linode\_monitor\_alert\_definitions

Retrieves Monitor Alert Definitions.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-alert-definitions).  (**Note: v4beta only.**)

## Example Usage

Retrieve ALL Monitor Alert Definitions:

```terraform
data "linode_monitor_alert_definitions" "all" {}
```

Retrieve Monitor Alert Definitions for a service type and filter by type:

```terraform
data "linode_monitor_alert_definitions" "all" {
  service_type = "dbaas"
  filter {
    name = type
    values = ["test-alert-definition"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_type` - (Optional) The service type (e.g., dbaas).

* [`filter`](#filter) - (Optional) A set of filters used to select IP addresses that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each alert definition will be stored in the `alert_definitions` attribute and will export the following attributes:

* `service_type` - The service type (e.g., dbaas).
* `id` - The ID of the alert definition.
* `label` - The label for the alert definition.
* `channel_ids` - A list of channel IDs to associate with the alert definition.
* `severity` - The severity level of the alert definition.
* [`rule_criteria`](#rule_criteria) - The criteria expression for the alert.
* [`trigger_conditions`](#trigger_conditions) - The conditions that need to be met to send a notification for the alert.
* `description` - A description for the alert definition.
* `entity_ids` - A list of entity IDs to associate with the alert definition.
* `status` -  The status of the alert definition.
* `type` - The type of alert. This can be either user for an alert specific to the current user, or system for one that applies to all users on your account.
* `has_more_resources` - Whether there are additional entity_ids associated with the alert for which the user doesn't have at least read-only access.
* `created` - The date and time the alert definition was created.
* `updated` - The date and time the alert definition was last updated.
* `created_by` - For a user alert definition, this is the user on your account that created it. For a system alert definition, this is returned as system.
* `updated_by` - For a user alert definition, this is the user on your account that last updated it. For a system alert definition, this is returned as system. If it hasn't been updated, this value is the same as created_by.
* `class` - "The plan type for the Managed Database cluster, either shared or dedicated. This only applies to a system alert for a service_type of dbaas (Managed Databases). For user alerts for dbaas, this is returned as null.",
* [`alert_channels`](#alert_channels) - A list of alert channel objects associated with the alert definition.

### rule_criteria

The following arguments are supported in the `rule_criteria` specification block:

* [`rules`](#rules) -  A list of rule objects defining the criteria for the alert.

#### rules

The following attributes are supported in each `rules` specification block:

* `aggregate_function` - The aggregate function to apply to the metric data.
* [`dimension_filters`](#dimension_filters) - A list of dimension filter objects to filter the metric data.
* `metric` - The metric to query.
* `operator` - The operator to apply to the metric. Allowed values: eq, gt, lt, gte, lte.
* `threshold` - The predefined value or condition that triggers an alert when met or exceeded.
* `label` - The name of the individual rule. This is used for display purposes in Akamai Cloud Manager.
* `unit` - The unit of the metric. This can be values like percent for percentage or GB for gigabyte.

##### dimension_filters

The following attributes are supported in each `dimension_filters` specification block:

* `dimension_label` - The label of the dimension to filter on.
* `operator` - The operator to apply to the dimension. Allowed values: eq, neq, startswith, endswith.
* `value` - The value to compare the dimension_label against.
* `label` - The name of the dimension filter. Used for display purposes.

### trigger_conditions

The following attributes are supported in the `trigger_conditions` specification block:

* `criteria_condition` - The logical operation applied. Currently only 'ALL' allowed.
* `evaluation_period_seconds` - Time period over which data is collected before evaluating the threshold.
* `polling_interval_seconds` - Frequency at which the metric data is polled.
* `trigger_occurrences` -  Number of times the condition must be met before triggering an alert.

### alert_channels

The following attributes are exported in each `alert_channels` block:

* `id` - The unique identifier assigned to the alert channel.
* `label` - The label of the alert channel.
* `type` - The type of alert channel.
* `url` - The URL of the alert channel.
