---
page_title: "Linode: linode_prefix_list"
description: |-
  Provides details about a Prefix List.
---

# Data Source: linode\_prefix\_list

Provides details about a Linode Prefix List.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-prefix-list).

## Example Usage

```terraform
data "linode_prefix_list" "example" {
  id = "12345"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Prefix List.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the Prefix List (e.g. `pl:system:object-storage:us-iad`, `pl::customer:my-list`).

* `description` - A description of the Prefix List.

* `visibility` - The visibility of the Prefix List. (`account`, `restricted`)

* `source_prefixlist_id` - The ID of the source prefix list, if this is a derived list.

* `ipv4` - A list of IPv4 addresses or networks in CIDR format contained in this prefix list.

* `ipv6` - A list of IPv6 addresses or networks in CIDR format contained in this prefix list.

* `version` - The version number of this Prefix List.

* `created` - When this Prefix List was created.

* `updated` - When this Prefix List was last updated.
