---
page_title: "Linode: linode_object_storage_global_quotas"
description: |-
  Provides details about Object Storage global quotas on your account.
---

# linode\_object\_storage\_global\_quotas

Provides details about Object Storage global quotas on your account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-global-quotas).

## Example Usage

The following example shows how one might use this data source to list and filter Object Storage global quotas.

```hcl
data "linode_object_storage_global_quotas" "keys" {
  filter {
    name   = "quota_id"
    values = ["keys"]
  }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Object Storage global quotas that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Object Storage global quota will be stored in the `quotas` attribute and will export the following attributes:

* `quota_id` - The ID of the Object Storage global quota.

* `quota_name` - The name of the Object Storage global quota.

* `description` - The description of the Object Storage global quota.

* `quota_limit` - The maximum quantity of the `resource_metric` allowed by the quota.

* `resource_metric` - The specific Object Storage resource for the quota.

* `quota_type` - The type of the Object Storage global quota.

* `has_usage` - Whether usage data is available for the Object Storage global quota.

## Filterable Fields

* `quota_id`

* `quota_name`

* `description`

* `quota_limit`

* `resource_metric`

* `quota_type`

* `has_usage`
