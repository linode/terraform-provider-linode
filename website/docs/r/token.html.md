---
layout: "linode"
page_title: "Linode: linode_token"
sidebar_current: "docs-linode-resource-token"
description: |-
  Manages a Linode Token.
---

# linode\_token

Provides a Linode Token resource.  This can be used to create, modify, and delete Linodes tokens.
For more information, see the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/getTokens).

## Example Usage

The following example shows how one might use this resource to configure a Token for access to Linode resources.

```hcl
resource "linode_token" "foo" {
    label = "token"
    scopes = "linodes:read_only"
    expiry = "2100-01-02T03:04:05Z"
}
```

## Argument Reference

The following arguments are supported:

* `label` - A label for the Token.

* `scopes` - The scopes this token was created with. These define what parts of the Account the token can be used to access. Many command-line tools, such as the Linode CLI, require tokens with access to *. Tokens with more restrictive scopes are generally more secure.

* `expiry` - When this token will expire. Personal Access Tokens cannot be renewed, so after this time the token will be completely unusable and a new token will need to be generated. Tokens may be created with 'null' as their expiry and will never expire unless revoked.

## Attributes

This resource exports the following attributes:

* `token` - The token used to access the API.

* `created` - The date this Token was created.

## Import

Linodes Tokens can be imported using the Linode Token `id`, e.g.  The secret token will not be imported.

```sh
terraform import linode_token.mytoken 1234567
```
