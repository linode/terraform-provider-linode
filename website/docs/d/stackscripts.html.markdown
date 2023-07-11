---
layout: "linode"
page_title: "Linode: linode_stackscripts"
sidebar_current: "docs-linode-datasource-stackscripts"
description: |-
  Provides information about Linode StackScripts that match a set of filters.
---

# linode\_stackscripts

Provides information about Linode StackScripts that match a set of filters.

**NOTICE:** Due to the large number of public StackScripts, this data source may time out if `is_public` is not filtered on.

## Example Usage

The following example shows how one might use this data source to access information about a Linode StackScript.

```hcl
data "linode_stackscripts" "specific-stackscripts" {
  filter {
    name = "label"
    values = ["my-cool-stackscript"]
  }

  filter {
    name = "is_public"
    values = [false]
  }
}

output "stackscript_id" {
  value = data.linode_stackscripts.specific-stackscripts.stackscripts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `latest` - (Optional) If true, only the latest StackScript will be returned. StackScripts without a valid `created` field are not included in the result.

* [`filter`](#filter) - (Optional) A set of filters used to select Linode StackScripts that meet certain requirements.

* `order_by` - (Optional) The attribute to order the results by. See the [Filterable Fields section](#filterable-fields) for a list of valid fields.

* `order` - (Optional) The order in which results should be returned. (`asc`, `desc`; default `asc`)

### Filter

* `name` - (Required) The name of the field to filter by. See the [Filterable Fields section](#filterable-fields) for a complete list of filterable fields.

* `values` - (Required) A list of values for the filter to allow. These values should all be in string form.

* `match_by` - (Optional) The method to match the field by. (`exact`, `regex`, `substring`; default `exact`)

## Attributes Reference

Each Linode StackScript will be stored in the `stackscripts` attribute and will export the following attributes:

* `id` - The unique ID of the StackScript.

* `label` - The StackScript's label is for display purposes only.

* `script` - The script to execute when provisioning a new Linode with this StackScript.

* `description` - A description for the StackScript.

* `rev_note` - This field allows you to add notes for the set of revisions made to this StackScript.

* `is_public` - This determines whether other users can use your StackScript. Once a StackScript is made public, it cannot be made private.

* `images` - An array of Image IDs representing the Images that this StackScript is compatible for deploying with.

* `deployments_active` - Count of currently active, deployed Linodes created from this StackScript.

* `user_gravatar_id` - The Gravatar ID for the User who created the StackScript.

* `deployments_total` - The total number of times this StackScript has been deployed.

* `username` - The User who created the StackScript.

* `created` - The date this StackScript was created.

* `updated` - The date this StackScript was updated.

* `user_defined_fields` - This is a list of fields defined with a special syntax inside this StackScript that allow for supplying customized parameters during deployment.

  * `label` - A human-readable label for the field that will serve as the input prompt for entering the value during deployment.

  * `name` - The name of the field.

  * `example` - An example value for the field.

  * `one_of` - A list of acceptable single values for the field.

  * `many_of` - A list of acceptable values for the field in any quantity, combination or order.

  * `default` - The default value. If not specified, this value will be used.

## Filterable Fields

* `deployments_active`

* `deployments_total`

* `description`

* `images`

* `is_public`

* `label`

* `mine`

* `rev_note`

* `username`
