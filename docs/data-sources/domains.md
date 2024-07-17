---
page_title: "Linode: linode_domains"
description: |-
  Provides information about Linode Cloud Domains that match a set of filters.
---

# Data Source: linode\_domains

Provides information about Linode Domains that match a set of filters.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-domains).

## Example Usage

Get information about all Linode Cloud Domains with a specific tag:

```hcl
data "linode_domains" "specific" {
  filter {
    name = "tags"
    values = ["test-tag"]
  }
}

output "domain" {
  value = data.linode_domains.specific.domains.0.domain
}
```

## Argument Reference

The following arguments are supported:

* [`filter`](#filter) - (Optional) A set of filters used to select Linode Cloud Domains that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode Domain will be stored in the `domains` attribute and will export the following attributes:

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

* `tags` - An array of tags applied to this object. Tags are case-insensitive and are for organizational purposes only.

## Filterable Fields

* `group`

* `tags`

* `domain`

* `type`

* `status`

* `description`

* `master_ips`

* `axfr_ips`

* `ttl_sec`

* `retry_sec`

* `expire_sec`

* `refresh_sec`

* `soa_email`
