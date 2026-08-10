---
page_title: "Linode: linode_account_transfer"
description: |-
  Provides information about Linode account network transfer utilization for the current month.
---

# Data Source: linode_account_transfer

Provides information about Linode account network transfer utilization for the current month.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-transfer).

## Example Usage

The following example shows how one might use this data source to access account network transfer details.

```hcl
data "linode_account_transfer" "transfer" {}
```

## Argument Reference

There are no supported arguments because the provider `token` can only access the associated account.

## Attributes Reference

The Linode Account Transfer data source exports the following attributes:

* `billable` - The amount of your transfer pool that is billable this billing cycle.

* `quota` - The amount of network usage allowed this billing cycle.

* `used` - The amount of network usage you have used this billing cycle.

* `region_transfers` - A list of network utilization details for regions with separate utilization quotas and rates.

  * `id` - The Region ID for this network utilization data.

  * `billable` - The amount of your transfer pool that is billable this billing cycle for this Region.

  * `quota` - The amount of network usage allowed this billing cycle for this Region.

  * `used` - The amount of network usage you have used this billing cycle for this Region.
