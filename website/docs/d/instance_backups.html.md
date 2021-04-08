---
layout: "linode"
page_title: "Linode: linode_instance_backups"
sidebar_current: "docs-linode-datasource-instance-backups"
description: |-
Provides details about the backups of an Instance.
---

# Data Source: linode\_instance_backups

Provides details about the backups of an Instance.

## Example Usage

```terraform
data "linode_instance_backups" "my-backups" {
    id = 123
}
```

## Argument Reference

The following arguments are supported:

* `linode_id` - (Required) The Linode instance's ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* [`automatic`](#backup) - A list of backups or snapshots for a Linode.

* [`current`](#backup) - The current Backup for a Linode.

* [`in_progress`](#backup) - The in-progress Backup for a Linode.

### Backup

The following attributes are available for each Backup:

* `id` - The unique ID of this Backup.

* `label` - A label for Backups that are of type snapshot.

* `status` - The current state of a specific Backup.

* `type` - This indicates whether the Backup is an automatic Backup or manual snapshot taken by the User at a specific point in time.

* `created` - The date the Backup was taken.

* `updated` - The date the Backup was most recently updated.

* `finished` - The date the Backup completed.

* `configs` - A list of the labels of the Configuration profiles that are part of the Backup.

* [`disks`](#disk) - A list of the disks that are part of the Backup.

### Disk

The following attributes are available for each disk:

* `label` - The label of this disk.

* `size` - The size of this disk.

* `filesystem` - The filesystem of this disk.
