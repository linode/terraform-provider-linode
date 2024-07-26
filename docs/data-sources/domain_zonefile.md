---
page_title: "Linode: linode_domain_zonefile"
description: |-
  Provides details about a Linode Domain Zonefile.
---

# Data Source: linode_domain_zonefile

Provides information about a Linode Domain Zonefile.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-domain-zone).

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

## Attributes Reference

The Linode Volume resource exports the following attributes:

- `domain_id` - The associated domain's unique ID.

- `zone_file` - Array of strings representing the Domain Zonefile.
