---
page_title: "Linode: linode_sshkey"
description: |-
  Provides details about a profile SSH Key
---

# Data Source: linode\_sshkey

`linode_sshkey` provides access to a specifically labeled SSH Key in the Profile of the User identified by the access token.

## Example Usage

The following example shows how the resource might be used to obtain the name of the SSH Key configured on the Linode user profile.

```hcl
data "linode_sshkey" "foo" {
  label = "foo"
}
```

## Argument Reference

- `label` - (Required) The label of the SSH Key to select.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the SSH Key

- `ssh_key` - The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.

- `created` - The date this key was added.
