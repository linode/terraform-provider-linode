---
page_title: "Linode: linode_lock"
description: |-
  Provides details about a Linode Lock.
---

# Data Source: linode\_lock

~> **Early Access:** Lock functionality is in early access and may not be available to all users.

~> **Important** Only unrestricted users can view locks. Restricted users cannot access lock information even if they have permissions for the resource.

Provides information about a Linode Lock. Locks prevent accidental deletion, rebuild operations, and service transfers of resources.

For more information, see the Linode APIv4 docs (TBD).

## Example Usage

```hcl
data "linode_lock" "my_lock" {
  id = 123456
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of the Lock.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `entity_id` - The ID of the locked entity.

* `entity_type` - The type of the locked entity.

* `lock_type` - The type of lock.

* `entity_label` - The label of the locked entity.

* `entity_url` - The URL of the locked entity.
