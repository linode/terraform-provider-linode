---
layout: "linode"
page_title: "Linode: linode_kernel"
sidebar_current: "docs-linode-datasource-kernel"
description: |-
Provides details about a Linode kernel.
---

# Data Source: linode\_kernel

Provides information about a Linode kernel

## Example Usage

The following example shows how one might use this data source to access information about a Linode kernel.

```hcl
data "linode_kernel" "latest" {
    id = "linode/latest-64bit"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique ID of this Kernel.

## Attributes

The Linode Kernel resource exports the following attributes:

* `architecture` - The architecture of this Kernel.

* `deprecated` - Whether or not this Kernel is deprecated.

* `kvm` - If this Kernel is suitable for KVM Linodes.

* `label` - The friendly name of this Kernel.

* `pvops` - If this Kernel is suitable for paravirtualized operations.

* `version` - Linux Kernel version

* `xen` - If this Kernel is suitable for Xen Linodes.
