---
page_title: "Linode: linode_database_access_controls"
description: |-
  Manages the access controls for a Linode Database.
---

# linode\_database\_access_controls

Manages the access control for a Linode Database. Only one `linode_database_access_controls` resource should be defined per-database.

## Example Usage

Grant a Linode access to a database:

```hcl
resource "linode_database_access_controls" "my-access" {
  database_id = linode_database_mysql.my-db.id
  database_type = "mysql"
  
  allow_list = [linode_instance.my-instance.ip_address]
}

resource "linode_instance" "my-instance" {
  label = "myinstance"
  region = "us-southeast"
  type = "g6-nanode-1"
  image = "linode/alpine3.14"
}

resource "linode_database_mysql" "my-db" {
  label = "mydatabase"
  engine_id = "mysql/8.0.30"
  region = "us-southeast"
  type = "g6-nanode-1"
}
```

## Argument Reference

The following arguments are supported:

* `database_id` - (Required) The unique ID of the target database.

* `database_type` - (Required) The unique type of the target database. (`mysql`, `mongodb`, `postgresql`)

* `allow_list` - (Required) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format.
