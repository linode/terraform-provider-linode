---
page_title: "Linode: linode_monitor_logs_destination"
description: |-
  Manages a Linode Monitor Logs Destination.
---

# linode\_monitor\_logs\_destination

Manages a Linode Monitor Logs Destination.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/post-logs-destination). (**Note: v4beta only.**)

## Example Usage

Creating a logs destination using Akamai Object Storage:

```terraform
resource "linode_object_storage_key" "key" {
  label   = "my-logs-key"
  regions = ["us-east"]
}

resource "linode_object_storage_bucket" "bucket" {
  region     = "us-east"
  label      = "my-logs-bucket"
  access_key = linode_object_storage_key.key.access_key
  secret_key = linode_object_storage_key.key.secret_key
}

resource "linode_monitor_logs_destination" "example" {
  label = "my-logs-destination"
  type  = "akamai_object_storage"

  akamai_object_storage_details = {
    access_key_id     = linode_object_storage_key.key.access_key
    access_key_secret = linode_object_storage_key.key.secret_key
    bucket_name       = linode_object_storage_bucket.bucket.label
    host              = linode_object_storage_bucket.bucket.hostname
  }
}
```

Creating a logs destination using a custom HTTPS endpoint:

```terraform
resource "linode_monitor_logs_destination" "https_example" {
  label = "my-https-destination"
  type  = "custom_https"

  custom_https_details = {
    endpoint_url     = "https://logs.example.com/ingest"
    content_type     = "application/json"
    data_compression = "gzip"

    authentication = {
      type     = "basic"
      username = "myuser"
      password = "mypassword"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label for this logs destination.

* `type` - (Required, Forces New) The type of this logs destination. One of: `akamai_object_storage`, `custom_https`.

* [`akamai_object_storage_details`](#akamai_object_storage_details) - (Optional) Details for an `akamai_object_storage` logs destination. Exactly one of `akamai_object_storage_details` or `custom_https_details` must be specified.

* [`custom_https_details`](#custom_https_details) - (Optional) Details for a `custom_https` logs destination. Exactly one of `akamai_object_storage_details` or `custom_https_details` must be specified.

### akamai_object_storage_details

* `access_key_id` - (Required) The access key ID for the object storage bucket.

* `access_key_secret` - (Required, Sensitive) The access key secret for the object storage bucket. This value is write-only and will not be returned by the API.

* `bucket_name` - (Required) The name of the object storage bucket.

* `host` - (Required) The hostname of the object storage endpoint.

* `path` - (Optional) The path within the bucket where logs will be stored.

### custom_https_details

* `endpoint_url` - (Required) The HTTPS endpoint URL to send logs to.

* `content_type` - (Required) The content type of the log data. One of: `application/json`, `application/json; charset=utf-8`.

* `data_compression` - (Required) The compression format for log data. One of: `none`, `gzip`.

* [`authentication`](#authentication) - (Required) Authentication configuration for the HTTPS endpoint.

* [`client_certificate_details`](#client_certificate_details) - (Optional) TLS client certificate configuration.

* [`custom_headers`](#custom_headers) - (Optional) Custom HTTP headers to include in log delivery requests.

#### authentication

* `type` - (Required) The authentication type. One of: `basic`, `none`.

* `username` - (Optional, Sensitive) The username for basic authentication. This value is write-only and will not be returned by the API.

* `password` - (Optional, Sensitive) The password for basic authentication. This value is write-only and will not be returned by the API.

#### client_certificate_details

* `tls_hostname` - (Required) The TLS hostname for certificate verification.

* `client_ca_certificate` - (Required, Sensitive) The client CA certificate. This value is write-only and will not be returned by the API.

* `client_certificate` - (Required, Sensitive) The client certificate. This value is write-only and will not be returned by the API.

* `client_private_key` - (Required, Sensitive) The client private key. This value is write-only and will not be returned by the API.

#### custom_headers

* `name` - (Required) The name of the HTTP header.

* `value` - (Required, Sensitive) The value of the HTTP header. This value is write-only and will not be returned by the API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of this logs destination.

* `status` - The status of this logs destination.

* `created_by` - The user who created this logs destination.

* `updated_by` - The user who last updated this logs destination.

* `created` - When this logs destination was created.

* `updated` - When this logs destination was last updated.

* `version` - The version of this logs destination.

## Import

Monitor Logs Destinations can be imported using the `id`, e.g.

```sh
terraform import linode_monitor_logs_destination.example 12345
```
