---
page_title: "Linode: linode_instance_backups"
description: |-
  Provides details about the backups of an Instance.
---

# Data Source: linode\_instance_backups

Provides details about the backups of an Instance.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-backups).

## Example Usage

```terraform
data "linode_instance_backups" "my-backups" {
    linode_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The Linode instance's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

**NOTE:** Nested fields are tagged as either **Block** (declared as `field { ... }`) or **Nested Attribute** (declared as `field = { ... }`). See the [Blocks vs. Nested Attributes](../guides/blocks_vs_nested_attributes.md) guide for details.

* [`automatic`](#backup) - (Read-Only Object List) A list of backups or snapshots for a Linode. Referenced with an index (e.g. `automatic.0.id`).

* [`current`](#backup) - (Read-Only Object List) The current Backup for a Linode. Referenced with an index (e.g. `current.0.id`).

* [`in_progress`](#backup) - (Read-Only Object List) The in-progress Backup for a Linode. Referenced with an index (e.g. `in_progress.0.id`).

### Backup

The following attributes are available for each Backup:

* `id` - The unique ID of this Backup.

* `label` - A label for Backups that are of type snapshot.

* `status` - The current state of a specific Backup. (`paused`, `pending`, `running`, `needsPostProcessing`, `successful`, `failed`, `userAborted`)

* `type` - This indicates whether the Backup is an automatic Backup or manual snapshot taken by the User at a specific point in time. (`auto`, `snapshot`)

* `created` - The date the Backup was taken.

* `updated` - The date the Backup was most recently updated.

* `finished` - The date the Backup completed.

* `configs` - A list of the labels of the Configuration profiles that are part of the Backup.

* [`disks`](#disk) - (Read-Only Object List) A list of the disks that are part of the Backup. Referenced with an index (e.g. `disks.0.label`).

### Disk

The following attributes are available for each disk:

* `label` - The label of this disk.

* `size` - The size of this disk.

* `filesystem` - The filesystem of this disk.
