---
page_title: "Linode: linode_object_storage_quotas"
description: |-
  Provides details about a list of Object Storage quotas information on your account.
---

# linode\_object\_storage\_quotas

Provides details about a list of Object Storage quotas information on your account.
For more information, see the [Linode APIv4 docs](TBD).

## Example Usage

The following example shows how one might use this data source to list and filter information about Object Storage quotas.

```hcl
data "linode_object_storage_quotas" "max_buckets_quotas" {
  filter {
    name = "endpoint_type"
    values = ["E0"]
  }
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode account availabilities that meet certain requirements.

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Object Storage quota will be stored in the `quotas` attribute and will export the following attributes:

* `id` - The ID of the Object Storage quota.

* `quota_name` - The name of the Object Storage quota.

* `endpoint_type` - The type of the S3 endpoint of the Object Storage.

* `s3_endpoint` - The S3 endpoint URL of the Object Storage, based on the `endpoint_type` and `region`.

* `description` - The description of the Object Storage quota.

* `quota_limit` - The maximum quantity of the `resource_metric` allowed by the quota.

* `resource_metric` - The specific Object Storage resource for the quota.

## Filterable Fields

* `id`

* `quota_name`

* `endpoint_type`

* `s3_endpoint`

* `description`

* `quota_limit`

* `resource_metric`
