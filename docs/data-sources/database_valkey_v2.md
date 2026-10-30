---
page_title: "Linode: linode_database_valkey_v2 data source"
description: |-
  Reads a Linode Valkey Managed Database.
---

# linode_database_valkey_v2 data source

Reads the configuration and status of an existing Valkey Managed Database.
The database must be accessible to the configured Linode API token.

## Example Usage

```hcl
data "linode_database_valkey_v2" "cache" {
  id = "12345"
}

output "valkey_host" {
  value = data.linode_database_valkey_v2.cache.host_primary
}
```

## Argument Reference

- `id` (Required) The numeric ID of the Valkey database.

## Attribute Reference

The data source returns the resource's computed attributes, including
`engine_id`, `label`, `region`, `type`, `status`, `version`, connection
hosts, `port`, `allow_list`, `private_network`, engine configuration,
`fork_source`, and `fork_restore_time`.

`available_restore_times` is the list of discrete snapshot timestamps that
can be used to fork this database. It is not a continuous point-in-time
recovery range.

`root_username`, `root_password`, and `ca_cert` are marked sensitive in
Terraform output. Terraform still stores them in state; protect and restrict
access to the state file and its backups.
