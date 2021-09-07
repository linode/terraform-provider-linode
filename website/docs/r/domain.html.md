---
layout: "linode"
page_title: "Linode: linode_domain_record"
sidebar_current: "docs-linode-resource-domain-record"
description: |-
  Manages a Linode Domain Record.
---

# linode\_domain

Provides a Linode Domain resource.  This can be used to create, modify, and delete Linode Domains through Linode's managed DNS service.
For more information, see [DNS Manager](https://www.linode.com/docs/platform/manager/dns-manager/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createDomain).

The Linode Guide, [Deploy a WordPress Site Using Terraform and Linode StackScripts](https://www.linode.com/docs/applications/configuration-management/deploy-a-wordpress-site-using-terraform-and-linode-stackscripts/), demonstrates the management of Linode Domain resources in the context of Linode Instance running WordPress.

## Example Usage

The following example shows how one might use this resource to configure a Domain Record attached to a Linode Domain.

```hcl
resource "linode_domain" "foobar" {
    type = "master"
    domain = "foobar.example"
    soa_email = "example@foobar.example"
    tags = ["foo", "bar"]
}

resource "linode_domain_record" "foobar" {
    domain_id = linode_domain.foobar.id
    name = "www"
    record_type = "CNAME"
    target = "foobar.example"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.

* `type` - (Required) If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).

* `soa_email` - (Required) Start of Authority email address. This is required for master Domains.

* `master_ips` - (Required for type="slave") The IP addresses representing the master DNS for this Domain.

- - -

* `status` - (Optional) Used to control whether this Domain is currently being rendered (defaults to "active").

* `description` - (Optional) A description for this Domain. This is for display purposes only.

* `group` - (Optional) The group this Domain belongs to. This is for display purposes only.

* `ttl_sec` - (Optional) 'Time to Live' - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.

* `retry_sec` - (Optional) The interval, in seconds, at which a failed refresh should be retried. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.

* `expire_sec` - (Optional) The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.

* `refresh_sec` - (Optional) The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.

* `axfr_ips` - (Optional) The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.

* `tags` - (Optional) A list of tags applied to this object. Tags are for organizational purposes only.

## Attributes

This resource exports no additional attributes, however `status` may reflect degraded states.

## Import

Linodes Domains can be imported using the Linode Domain `id`, e.g.

```sh
terraform import linode_domain.foobar 1234567
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for Domains and other Linode resource types.
