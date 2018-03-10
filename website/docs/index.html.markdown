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
  key = "$LINODE_API_KEY"
}

resource "linode_linode" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `key` - (Required) This is your [Linode APIv3 Key](https://linode.com/docs/platform/api/api-key/).

   The Linode API key can also be specified using the `LINODE_API_KEY` environment variable.

