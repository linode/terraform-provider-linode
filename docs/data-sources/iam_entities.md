---
page_title: "Linode: linode_iam_entities"
description: |-
  Lists all entities on the account.
---

# linode\_iam\_entities

Provides a list of all Entities on this Account.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-entities).

```hcl
data "linode_iam_entities" "entities" {
}
```

## Argument Reference

The following arguments are supported:

*While the filtering system is in place currently this endpoint does not have filter support for any field.*

* [`filter`](#filter) - (Optional) A set of filters used to select Linode users that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

*While the filtering system is in place currently this endpoint does not have filter support for any field.*

## Attributes Reference

Each Linode entity will be stored in the `entities` attribute and will export the following attributes:

* `id` - A unique identifier for each entity.

* `label` - A unique label for each entity.

* `type` - The type for each entity. (eg. Volume)

## Filterable Fields

*While the filtering system is in place currently this endpoint does not have filter support for any field.*
