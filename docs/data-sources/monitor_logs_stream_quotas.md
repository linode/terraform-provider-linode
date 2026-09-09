---
page_title: "Linode: linode_monitor_logs_stream_quotas"
description: |-
  Provides details about logs stream quotas.
---

# linode\_monitor\_logs\_stream\_quotas

Provides details about logs stream quotas available.

## Example Usage

The following example shows how one might use this data source to list logs stream quotas.

```hcl
data "linode_monitor_logs_stream_quotas" "all" {}
```

## Attributes Reference

The returned logs stream quotas are stored in the `quotas` attribute and export the following attributes:

* `quotas` - (Nested Attribute List) The list of logs stream quotas available.

* `quota_id` - The ID of the logs stream quota.

* `quota_name` - The name of the logs stream quota.

* `description` - The description of the logs stream quota.

* `quota_limit` - The maximum number allowed by this quota.

* `quota_type` - The type of the logs stream quota.
