---
layout: "linode"
page_title: "Linode: linode_database_mysql_firewall"
sidebar_current: "docs-linode-database-mysql-firewalloh "
description: |-
  Manages the access control for a Linode MySQL Database.
---

# linode\_database\_mysql\_firewall

**NOTICE:** Managed Databases are currently in beta. Ensure `api_version` is set to `v4beta` in order to use this resource.

Manages the access control for a Linode MySQL Database. Only one `linode_database_mysql_firewall` resource should be defined per-database.

## Example Usage

Grant a Linode access to a MySQL database:

```hcl
resource "linode_database_mysql_firewall" "my-firewall" {
  database_id = linode_database_mysql.my-db.id
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
  engine_id = "mysql/8.0.26"
  region = "us-southeast"
  type = "g6-nanode-1"
}
```

## Argument Reference

The following arguments are supported:

* `database_id` - (Required) The unique ID of the target MySQL database.

* `allow_list` - (Required) A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format.
