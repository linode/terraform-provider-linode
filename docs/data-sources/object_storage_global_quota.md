---
page_title: "Linode: linode_object_storage_global_quota"
description: |-
  Provides details about an Object Storage global quota on your account.
---

# linode\_object\_storage\_global\_quota

Provides details about an Object Storage global quota on your account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-global-quota).

## Example Usage

The following example shows how one might use this data source to access information about an Object Storage global quota.

```hcl
data "linode_object_storage_global_quota" "my_quota" {
  quota_id = "keys"
}
```

## Argument Reference

The following arguments are supported:

* `quota_id` - (Required) The Object Storage global quota ID.

## Attributes Reference

The Linode Object Storage global quota data source exports the following attributes:

* `quota_name` - The name of the Object Storage global quota.

* `description` - The description of the Object Storage global quota.

* `quota_limit` - The maximum quantity of the `resource_metric` allowed by the quota.

* `resource_metric` - The specific Object Storage resource for the quota.

* `quota_type` - The type of the Object Storage global quota.

* `has_usage` - Whether usage data is available for the Object Storage global quota.

* `quota_usage` - The usage data for a specific global Object Storage quota on your account. This value is `null` when `has_usage` is `false`. For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-global-quota-usage).

  * `quota_limit` - The maximum quantity allowed by the quota.

  * `usage` - The quantity of the Object Storage resource currently in use.

* `id` - The unique ID of the Object Storage global quota data source.
