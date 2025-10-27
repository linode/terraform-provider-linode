---
page_title: "Linode: linode_producer_image_share_group_member"
description: |-
  Provides details about a Member of an Image Share Group.
---

# Data Source: linode\_producer\_image\_share\_group\_member

`linode_producer_image_share_group_member` provides details about a Member of an Image Share Group.
For more information, see the [Linode APIv4 docs](TODO).


## Example Usage

The following example shows how the datasource might be used to obtain additional information about a member of an Image Share Group.

```hcl
data "linode_producer_image_share_group_member" "member" {
  sharegroup_id = 12345
  token_uuid = "db58ab2e-3021-4b08-9426-8e456f6dd268
}
```

## Argument Reference

* `sharegroup_id` - (Required) The ID of the Image Share Group the member belongs to.

* `token_uuid` - (Required) The UUID of member's token.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the member.

* `status` - The status of the member.

* `created` - When the member was created.

* `updated` - When the member was last updated.

* `expiry` - When the member will expire.
