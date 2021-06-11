---
layout: "linode"
page_title: "Linode: linode_domain"
sidebar_current: "docs-linode-datasource-domain"
description: |-
  Provides details about a Linode domain.
---

# Data Source: linode\_domain

Provides information about a Linode domain.

## Example Usage

The following example shows how one might use this data source to access information about a Linode domain.

```hcl
data "linode_domain" "foo" {
    id = "1234567"
}

data "linode_domain" "bar" {
    domain = "bar.example.com"
}
```

## Argument Reference

The following arguments are supported, at least one is required:

* `id` - (Optional) The unique numeric ID of the Domain record to query.

* `domain` - (Optional) The unique domain name of the Domain record to query.

## Attributes

The Linode Domain resource exports the following attributes:

* `id` - The unique ID of this Domain.

* `domain` - The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain

* `type` - If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave) (`master`, `slave`)

* `group` - The group this Domain belongs to.

* `status` - Used to control whether this Domain is currently being rendered. (`disabled`, `active`)

* `description` - A description for this Domain.

* `master_ips` - The IP addresses representing the master DNS for this Domain.

* `axfr_ips` - The list of IPs that may perform a zone transfer for this Domain.

* `ttl_sec` - 'Time to Live'-the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers.

* `retry_sec` - The interval, in seconds, at which a failed refresh should be retried.

* `expire_sec` - The amount of time in seconds that may pass before this Domain is no longer authoritative.

* `refresh_sec` - The amount of time in seconds before this Domain should be refreshed.

* `soa_email` - Start of Authority email address.

* `tags` - An array of tags applied to this object.
