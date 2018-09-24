---
layout: "linode"
page_title: "Linode: linode_instance"
sidebar_current: "docs-linode-datasource-instance-type"
description: |-
  Provides details about a Linode instance type.
---

# Data Source: instance_type

Provides information about a Linode instance type

## Example Usage

The following example shows how one might use this data source to define a Linode instance.

```hcl
data "linode_instance_type" "default" {
    id = "g6-standard-2"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) Label used to identify instance type 

## Attributes

This Linode Instance resource exports the following attributes:

* `id` - The ID representing the Linode Type

* `label` - The Linode Type's label is for display purposes only

* `class` - The class of the Linode Type

* `disk` - The Disk size, in MB, of the Linode Type

* `price.0.hourly` -  Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `addons.0.backups.0.price.0.hourly` - The cost (in US dollars) per hour to add Backups service.

* `addons.0.backups.0.price.0.monthly` - The cost (in US dollars) per month to add Backups service.
