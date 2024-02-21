---
page_title: "Linode: linode_firewall_device"
description: |-
  Manages a Linode Firewall Device.
---

# linode\_firewall\_device

Manages a Linode Firewall Device.

**NOTICE:** Attaching a Linode Firewall Device to a `linode_firewall` resource with user-defined `linodes` may cause device conflicts.

## Example Usage

```terraform
resource "linode_firewall_device" "my_device" {
  firewall_id = linode_firewall.my_firewall.id
  entity_id = linode_instance.my_instance.id
}

resource "linode_firewall" "my_firewall" {
  label = "my_firewall"

  inbound {
    label    = "http"
    action = "ACCEPT"
    protocol  = "TCP"
    ports     = "80"
    ipv4 = ["0.0.0.0/0"]
    ipv6 = ["::/0"]
  }
  
  inbound_policy = "DROP"
  outbound_policy = "ACCEPT"
}

resource "linode_instance" "my_instance" {
  label      = "my_instance"
  region     = "us-southeast"
  type       = "g6-standard-1"
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Required) The unique ID of the target Firewall.

* `entity_id` - (Required) The unique ID of the entity to attach.

* `entity_type` - (Optional) The type of the entity to attach. (default: `linode`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created` - When the Firewall Device was last created.

* `updated` - When the Firewall Device was last updated.

## Import

Firewall Device can be imported using the `firewall_id` followed by the Firewall Device `id` separated by a comma, e.g.

```sh
terraform import linode_firewall_device.my_device_duplicated 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for various Linode resource types.
