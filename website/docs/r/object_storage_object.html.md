---
layout: "linode"
page_title: "Linode: linode_object_storage_object"
sidebar_current: "docs-linode-resource-object-storage-object"
description: |-
  Manages a Linode Object Storage Object.
---

# linode\_object\_storage\_object

Provides a Linode Object Storage Object resource. This can be used to create, modify, and delete Linodes Object Storage Objects for Buckets.

## Example Usage

### Uploading a file to a bucket

```hcl
resource "linode_object_storage_object" "object" {
    bucket  = "my-bucket"
    cluster = "us-east-1"
    key     = "my-object"

    secret_key = linode_object_storage_key.my_key.secret_key
    access_key = linode_object_storage_key.my_key.access_key

    source = pathexpand("~/files/log.txt")
}

```

### Uploading plaintext to a bucket

```hcl
resource "linode_object_storage_object" "object" {
    bucket  = "my-bucket"
    cluster = "us-east-1"
    key     = "my-object"

    secret_key = linode_object_storage_key.my_key.secret_key
    access_key = linode_object_storage_key.my_key.access_key

    content          = "This is the content of the Object..."
    content_type     = "text/plain"
    content_language = "en"
}

```

## Argument Reference

-> **Note:** If you specify `content_encoding` you are responsible for encoding the body appropriately. `source`, `content`, and `content_base64` all expect already encoded/compressed bytes.

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to put the object in.

* `cluster` - (Required) The cluster the bucket is in.

* `key` - (Required) They name of the object once it is in the bucket.

* `secret_key` - (Required) The secret key to authenitcate with.

* `access_key` - (Required) The access key to authenticate with.

* `source` - (Optional, conflicts with `content` and `content_base64`) The path to a file that will be read and uploaded as raw bytes for the object content. The path must either be relative to the root module or absolute.

* `content` - (Optional, conflicts with `source` and `content_base64`) Literal string value to use as the object content, which will be uploaded as UTF-8-encoded text.

* `content_base64` - (Optional, conflicts with `source` and `content`) Base64-encoded data that will be decoded and uploaded as raw bytes for the object content. This allows safely uploading non-UTF8 binary data, but is recommended only for small content such as the result of the `gzipbase64` function with small text strings. For larger objects, use `source` to stream the content from a disk file.

* `acl` - (Optional) The canned ACL to apply. Can be one of `private`, `public-read`, `authenticated-read`, `public-read-write`, and `custom` (defaults to `private`).

* `cache_control` - (Optional) Specifies caching behavior along the request/reply chain Read [w3c cache_control](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9) for further details.

* `content_disposition` - (Optional) Specifies presentational information for the object. Read [w3c content_disposition](http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1) for further information.

* `content_encoding` - (Optional) Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [w3c content encoding](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11) for further information.

* `content_language` - (Optional) The language the content is in e.g. en-US or en-GB.

* `content_type` - (Optional) A standard MIME type describing the format of the object data, e.g. application/octet-stream. All Valid MIME Types are valid for this input.

* `website_redirect` - (Optional) Specifies a target URL for website redirect.

* `etag` - (Optional) Used to trigger updates. The only meaningful value is `${filemd5("path/to/file")}` (Terraform 0.11.12 or later) or `${md5(file("path/to/file"))}` (Terraform 0.11.11 or earlier).

* `metadata` - (Optional) A map of keys/values to provision metadata.

* `force_destroy` - (Optional) Allow the object to be deleted regardless of any legal hold or object lock (defaults to `false`).

## Attributes Reference

The following attributes are exported

* `version_id` - A unique version ID value for the object.
