---
page_title: "Linode: linode_instance"
description: |-
  Provides details about a Linode instance type.
---

# Data Source: linode\_instance\_type

Provides information about a Linode instance type
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-linode-type).

## Example Usage

The following example shows how one might use this data source to access information about a Linode Instance type.

```hcl
data "linode_instance_type" "default" {
    id = "g6-standard-2"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) Label used to identify instance type

## Attributes Reference

The Linode Instance Type resource exports the following attributes:

* `id` - The ID representing the Linode Type

* `label` - The Linode Type's label is for display purposes only

* `class` - The class of the Linode Type. See all classes [here](https://techdocs.akamai.com/linode-api/reference/get-linode-type).

* `disk` - The Disk size, in MB, of the Linode Type

* `price.0.hourly` -  Cost (in US dollars) per hour.

* `price.0.monthly` - Cost (in US dollars) per month.

* `addons.0.backups.0.price.0.hourly` - The cost (in US dollars) per hour to add Backups service.

* `addons.0.backups.0.price.0.monthly` - The cost (in US dollars) per month to add Backups service.

* `region_prices.*.id` - The name of the Linode Region this entry corresponds to.

* `region_prices.*.hourly` - The hourly price of this plan for the given region.

* `region_prices.*.monthly` - The monthly price of this plan for the given region.

* `network_out` - The Mbits outbound bandwidth allocation.

* `memory` - The amount of RAM included in this Linode Type.

* `transfer` - The monthly outbound transfer amount, in MB.

* `vcpus` - The number of VCPU cores this Linode Type offers.

* `accelerated_devices` - The number of VPUs this Linode Type offers.
