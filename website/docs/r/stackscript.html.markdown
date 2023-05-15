---
layout: "linode"
page_title: "Linode: linode_stackscript"
sidebar_current: "docs-linode-resource-stackscript"
description: |-
  Manages a Linode StackScript.
---

# linode\_stackscript

Provides a Linode StackScript resource.  This can be used to create, modify, and delete Linode StackScripts.  StackScripts are private or public managed scripts which run within an instance during startup.  StackScripts can include variables whose values are specified when the Instance is created.  

For more information, see [Automate Deployment with StackScripts](https://www.linode.com/docs/platform/stackscripts/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#tag/StackScripts).

The Linode Guide, [Deploy a WordPress Site Using Terraform and Linode StackScripts](https://www.linode.com/docs/applications/configuration-management/deploy-a-wordpress-site-using-terraform-and-linode-stackscripts/), shows how a public StackScript can be used to provision a Linode Instance.   The guide, [Create a Terraform Module](https://www.linode.com/docs/applications/configuration-management/create-terraform-module/), demonstrates StackScript use through a wrapping module.

## Example Usage

The following example shows how one might use this resource to configure a StackScript attached to a Linode Instance.  As shown below, StackScripts must begin with a shebang (`#!`).  The `<UDF ...>` element provided in the Bash comment block defines a variable whose value is provided when creating the Instance (or disk) using the `stackscript_data` field.

```hcl
resource "linode_stackscript" "foo" {
  label = "foo"
  description = "Installs a Package"
  script = <<EOF
#!/bin/bash
# <UDF name="package" label="System Package to Install" example="nginx" default="">
apt-get -q update && apt-get -q -y install $PACKAGE
EOF
  images = ["linode/ubuntu18.04", "linode/ubuntu16.04lts"]
  rev_note = "initial version"
}

resource "linode_instance" "foo" {
  image  = "linode/ubuntu18.04"
  label  = "foo"
  region = "us-east"
  type   = "g6-nanode-1"
  authorized_keys    = ["..."]
  root_pass      = "..."

  stackscript_id = linode_stackscript.foo.id
  stackscript_data = {
    "package" = "nginx"
  }
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The StackScript's label is for display purposes only.

* `script` - (Required) The script to execute when provisioning a new Linode with this StackScript.

* `description` - (Required) A description for the StackScript.

* `images` - (Required) A set of Image IDs representing the Images that this StackScript is compatible for deploying with. `any/all` indicates that all available image distributions, including private images, are accepted. Currently private image IDs are not supported.

- - -

* `rev_note` - (Optional) This field allows you to add notes for the set of revisions made to this StackScript.

* `is_public` - (Optional) This determines whether other users can use your StackScript. Once a StackScript is made public, it cannot be made private. *Changing `is_public` forces the creation of a new StackScript*

## Attributes Reference

This resource exports the following attributes:

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
