---
page_title: "Linode: linode_token"
description: |-
  Manages a Linode Token.
---

# linode\_token

Provides a Linode Token resource.  This can be used to create, modify, and delete Linode API Personal Access Tokens.  Personal Access Tokens proxy user credentials for Linode API access.  This is necessary for tools, such as Terraform, to interact with Linode services on a user's behalf.

It is common for Terraform itself to be configured with broadly scoped Personal Access Tokens.  Provisioning scripts or tools configured within a Linode Instance should follow the principle of least privilege to afford only the required roles for tools to perform their necessary tasks.  The `linode_token` resource allows for the management of Personal Access Tokens with scopes mirroring or narrowing the scope of the parent token.

For more information, see the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/getTokens).

## Example Usage

The following example shows how one might use this resource to configure a token for use in another tool that needs access to Linode resources.

```hcl
resource "linode_token" "foo" {
  label  = "token"
  scopes = "linodes:read_only"
  expiry = "2100-01-02T03:04:05Z"
}

resource "linode_instance" "foo" {
  # Configure the linode-cli and use it to add other Linode Instances to the hosts file
  provisioner "remote-exec" {
    inline = <<EOF
echo -e "[DEFAULT]\n token = ${linode_token.foo.token}\n region=${self.region}\n type=${self.type}" > ~/.linode-cli
pip install linode-cli
linode-cli linodes list --format "ipv6,label" --text --no-headers >> /etc/hosts
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `label` - A label for the Token.

* `scopes` - The scopes this token was created with. These define what parts of the Account the token can be used to access. Many command-line tools, such as the Linode CLI, require tokens with access to *. Tokens with more restrictive scopes are generally more secure. All scopes can be viewed in [the Linode API documentation](https://www.linode.com/docs/api/#oauth-reference).

* `expiry` - When this token will expire. Personal Access Tokens cannot be renewed, so after this time the token will be completely unusable and a new token will need to be generated. Tokens may be created with 'null' as their expiry and will never expire unless revoked.

## Attributes Reference

This resource exports the following attributes:

* `token` - The token used to access the API.

* `created` - The date this Token was created.

## Import

Linodes Tokens can be imported using the Linode Token `id`, e.g.  The secret token will not be imported.

```sh
terraform import linode_token.mytoken 1234567
```
