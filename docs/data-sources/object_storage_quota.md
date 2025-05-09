---
page_title: "Linode: linode_object_storage_quota"
description: |-
  Provides details about Object Storage quota information on your account.
---

# linode\_object\_storage\_quota

Provides details about Object Storage quota information on your account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-quota).

## Example Usage

The following example shows how one might use this data source to access information about an Object Storage quota.

```hcl
data "linode_object_storage_quota" "my_quota" {
  quota_id = "obj-buckets-br-gru-1.linodeobjects.com"
}
```

## Argument Reference

The following arguments are supported:

* `quota_id` - (Required) The Object Storage quota ID.

## Attributes Reference

The Linode Object Storage quota data source exports the following attributes:

* `quota_name` - The name of the Object Storage quota.

* `endpoint_type` - The type of the S3 endpoint of the Object Storage.

* `s3_endpoint` - The S3 endpoint URL of the Object Storage, based on the `endpoint_type` and `region`.

* `description` - The description of the Object Storage quota.

* `quota_limit` - The maximum quantity of the `resource_metric` allowed by the quota.

* `resource_metric` - The specific Object Storage resource for the quota.

* `quota_usage` - The usage data for a specific Object Storage related quota on your account. For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-quota-usage).

  * `quota_limit` - The maximum quantity allowed by the quota.
  
  * `usage` - The quantity of the Object Storage resource currently in use.

* `id` - The unique ID of the Object Storage quota data source.
