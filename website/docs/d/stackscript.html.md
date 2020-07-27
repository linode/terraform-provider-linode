---
layout: "linode"
page_title: "Linode: linode_stackscript"
sidebar_current: "docs-linode-datasource-stackscript"
description: |-
  Provides details about a Linode StackScript.
---

# linode\_stackscript

Provides details about a specific Linode StackScript.

## Example Usage

The following example shows how one might use this data source to access information about a Linode StackScript.

```hcl
data "linode_stackscript" "my_stackscript" {
    id = 355872
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique numeric ID of the StackScript to query.

## Attributes

This resource exports the following attributes:

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

## Import

Linodes StackScripts can be imported using the Linode StackScript `id`, e.g.

```sh
terraform import linode_stackscript.mystackscript 1234567
```
