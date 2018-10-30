---
layout: "linode"
page_title: "Provider: Linode"
sidebar_current: "docs-linode-index"
description: |-
  The Linode provider is used to interact with Linode services. The provider needs to be configured with the proper credentials before it can be used.
---

# Linode Provider

The Linode provider exposes data sources to interact with [Linode services](https://www.linode.com/).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available data sources.

## Example Usage

```hcl
# Configure the Linode provider
provider "linode" {
  token = "$LINODE_TOKEN"
}

resource "linode_instance" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `token` - (Required) This is your [Linode APIv4 Token](https://developers.linode.com/api/v4#section/Personal-Access-Token).

   The Linode Token can also be specified using the `LINODE_TOKEN` environment variable.

## Debugging

The [Linode APIv4 wrapper](https://github.com/linode/linodego) used by this provider accepts a `LINODE_DEBUG` environment variable.
If this variable is assigned to `1`, the request and response of all Linode API traffic will be reported through [Terraform debugging and logging facilities](/docs/internals/debugging.html).

Use of the `LINODE_DEBUG` variable in production settings is **strongly discouraged** with the `linode_account` datasource.  While Terraform does not directly store sensitive data from this datasource, the Linode Account API endpoint returns **sensitive data** such as the account `tax_id` (VAT) and the credit card `last_four` and `expiry`.  Be very cautious about storing this debug output.