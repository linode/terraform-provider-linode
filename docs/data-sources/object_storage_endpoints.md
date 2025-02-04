---
page_title: "Linode: linode_object_storage_endpoints"
description: |-
  Provides information about Linode Object Storage endpoints available to the user.
---

# Data Source: linode_object_storage_endpoints

Provides information about Linode Object Storage endpoints available to the user.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-object-storage-endpoints).

## Example Usage

Get an endpoint of E3 type (highest performance and capacity) of Linode Object Storage services:

```hcl
data "linode_object_storage_endpoints" "test" {
    filter {
        name = "endpoint_type"
        values = ["E3"]
    }
}

output "high-performance-obj-endpoint" {
  value = data.linode_object_storage_endpoints.test.endpoints[0].s3_endpoint
}
```

Get a list of all available endpoints of Linode Object Storage services.

```hcl
data "linode_object_storage_endpoints" "test" {}

output "available-endpoints" {
  value = data.linode_object_storage_endpoints.test.endpoints
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Object Storage endpoints that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Object Storage endpoint type will export the following attributes:

* `endpoint_type` - The type of `s3_endpoint` available to the active `user`. See [Endpoint types](https://techdocs.akamai.com/cloud-computing/docs/object-storage#endpoint-type) for more information.

* `region` - The Akamai cloud computing region, represented by its slug value. The [list regions](https://techdocs.akamai.com/linode-api/reference/get-regions) API is available to see all regions available.

* `s3_endpoint` -  Your s3 endpoint URL, based on the `endpoint_type` and `region`. Output as null if you haven't assigned an endpoint for your user in this region with the specific endpoint type.

## Filterable Fields

* `endpoint_type`

* `region`

* `s3_endpoint`
