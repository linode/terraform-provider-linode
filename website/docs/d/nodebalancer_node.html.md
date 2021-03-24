---
layout: "linode"
page_title: "Linode: linode_nodebalancer_node"
sidebar_current: "docs-linode-datasource-nodebalancer-node"
description: |-
Provides details about a NodeBalancer node.
---

# Data Source: linode\_nodebalancer_node

Provides details about a Linode NodeBalancer node.

## Example Usage

```terraform
data "linode_nodebalancer_node" "my-node" {
    id = 123

    nodebalancer_id = 456
    config_id = 789
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The node's ID.

* `nodebalancer_id` - (Required) The ID of the NodeBalancer that contains the node.

* `config_id` - (Required) The ID of the config that contains the Node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `label` - The label of the Linode NodeBalancer Node. This is for display purposes only.

* `address` - The private IP Address where this backend can be reached.

* `mode` - The mode this NodeBalancer should use when sending traffic to this backend. If set to `accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. If set to `drain` this backend will not receive new traffic, but connections already pinned to it will continue to be routed to it

* `weight` - Used when picking a backend to serve a request and is not pinned to a single backend yet. Nodes with a higher weight will receive more traffic. (1-255).

* `status` - The current status of this node, based on the configured checks of its NodeBalancer Config. (unknown, UP, DOWN).
