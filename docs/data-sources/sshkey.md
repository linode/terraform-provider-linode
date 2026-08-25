---
page_title: "Linode: linode_sshkey"
description: |-
  Provides details about a profile SSH Key
---

# Data Source: linode\_sshkey

`linode_sshkey` provides access to a specifically identified SSH Key in the Profile of the User identified by the access token.
For more information, see the [Linode APIv4 docs](https://techdocs.akamai.com/linode-api/reference/get-ssh-key).

## Example Usage

The following example shows how one might use this data source to access information about an SSH Key configured on the Linode user profile.

```hcl
data "linode_sshkey" "foo" {
  label = "foo"
}

data "linode_sshkey" "bar" {
  id = "1234567"
}
```

## Argument Reference

The following arguments are supported, exactly one is required:

- `id` - (Optional) The ID of the SSH Key to select. When set, `label` is computed from the API response.

- `label` - (Optional) The label of the SSH Key to select. When set, `id` is computed from the API response.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the SSH Key. Computed when `label` is used as the selector.

- `label` - The label of the SSH Key. Computed when `id` is used as the selector.

- `ssh_key` - The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.

- `created` - The date this key was added.
