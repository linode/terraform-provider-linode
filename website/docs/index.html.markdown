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
