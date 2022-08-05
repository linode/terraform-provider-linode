---
layout: "linode"
page_title: "Linode: linode_domain_zonefile"
sidebar_current: "docs-linode-datasource-domain-zonefile"
description: |-
  Provides details about a Linode Domain Zonefile.
---

# Data Source: linode_domain_zonefile

Provides information about a Linode Domain Zonefile.

## Example Usage

The following example shows how one might use this data source to access information about a Linode Domain Zonefile.

```hcl
data "linode_domain_zonefile" "my_zonefile" {
    domain_id = 3150401
}
```

## Argument Reference

The following argument is required:

- `domain_id` - (Required) The associated domain's unique ID.

## Attributes

The Linode Volume resource exports the following attributes:

- `domain_id` - The associated domain's unique ID.

- `zone_file` - Array of strings representing the Domain Zonefile.
