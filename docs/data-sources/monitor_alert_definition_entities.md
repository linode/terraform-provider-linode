---
page_title: "Linode: linode_monitor_alert_definition_entities"
description: |-
  Retrieves entities associated with a Monitor Alert Definition.
---

# linode\_monitor\_alert\_definition\_entities

Retrieves the entities associated with a specific Monitor Alert Definition.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-alert-definition-entities).  (**Note: v4beta only.**)

## Example Usage

Retrieve entities for a specific alert definition:

```terraform
data "linode_monitor_alert_definition_entities" "test" {
  service_type = "dbaas"
  alert_id     = 123
}
```

Retrieve entities filtered by type:

```terraform
data "linode_monitor_alert_definition_entities" "test" {
  service_type = "dbaas"
  alert_id     = 123

  filter {
    name   = "type"
    values = ["dbaas"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_type` - (Required) The service type for the alert definition (e.g., `dbaas`).

* `alert_id` - (Required) The unique identifier for the alert definition.

* [`filter`](#filter) - (Optional) A set of filters used to select entities that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each entity will be stored in the `entities` attribute and will export the following attributes:

* `id` - The unique identifier for this entity.

* `label` - The label of this entity.

* `url` - The URL for this entity.

* `type` - The type of this entity.

## Filterable Fields

* `id`

* `label`

* `type`

* `url`

