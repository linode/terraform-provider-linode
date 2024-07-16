---
page_title: "Linode: linode_account_availability"
description: |-
  Provides details about service availability in a region to an account specifically. 
---

# linode\_account\_availability

Provides details about service availability in a region to an account specifically.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-account-availability).

## Example Usage

The following example shows how one might use this data source to access information about a Linode account availability.

```hcl
data "linode_account_availability" "my_account_availability" {
    region = "us-east"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region ID.

## Attributes Reference

The Linode Account Availability data source exports the following attributes:

* `region` - The region ID.

* `unavailable` - A set of services which are unavailable to the account in a region.

* `available` - A set of services which are available to the account in a region.
