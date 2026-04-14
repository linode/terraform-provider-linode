---
page_title: "Linode: linode_tag"
description: |-
  Provides details about a Linode Tag.
---

# Data Source: linode\_tag

Provides information about a Linode Tag, including the objects associated with it.

## Example Usage

```hcl
data "linode_tag" "example" {
  label = "my-tag"
}
```

## Argument Reference

* `label` - (Required) The label of the tag to look up.

## Attributes Reference

* `id` - The label of the tag.

* `objects` - A list of objects associated with this tag. Each object has the following attributes:

  * `type` - The type of the tagged object (e.g. `linode`, `domain`, `volume`, `nodebalancer`, `reserved_ipv4_address`).

  * `id` - The ID of the tagged object. For `reserved_ipv4_address` objects, this is the IP address string.
