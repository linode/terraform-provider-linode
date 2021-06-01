---
layout: "linode"
page_title: "Linode: linode_daomin_record"
sidebar_current: "docs-linode-datasource-domain-record"
description: |-
  Provides details about a Linode Domain Record.
---

# Data Source: linode_domain_record

Provides information about a Linode Domain Record.

## Example Usage

The following example shows how one might use this data source to access information about a Linode Domain Record.

```hcl
data "linode_domain_record" "my_record" {
    id = "14950401"
    domain_id = "3150401"
}

data "linode_domain_record" "my_www_record" {
    name = "www"
    domain_id = "3150401"
}
```

## Argument Reference

The following argument is required:

- `id` (Optional) - The unique ID of the Domain Record.

- `name` (Optional) - The name of the Record.

- `domain_id` - (Required) The associated domain's unique ID.

## Attributes

The Linode Volume resource exports the following attributes:

- `id` - The unique ID of the Domain Record.

- `name` - The name of the Record.

- `domain_id` - The associated domain's unique ID.

- `type` - The type of Record this is in the DNS system. See all record types [here](https://www.linode.com/docs/api/domains/#domain-records-list__responses).

- `ttl_sec` - The amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers.

- `target` - The target for this Record. This field's actual usage depends on the type of record this represents. For A and AAAA records, this is the address the named Domain should resolve to.

- `priority` - The priority of the target host. Lower values are preferred.

- `weight` - The relative weight of this Record. Higher values are preferred.

- `port` - The port this Record points to.

- `protocol` - The protocol this Record's service communicates with. Only valid for SRV records.

- `service` - The service this Record identified. Only valid for SRV records.

- `tag` - The tag portion of a CAA record.
