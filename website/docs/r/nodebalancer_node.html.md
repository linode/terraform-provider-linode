---
layout: "linode"
page_title: "Linode: linode_nodebalancer_node"
sidebar_current: "docs-linode-resource-nodebalancer_node"
description: |-
  Manages a Linode NodeBalancer Node.
---

# linode\_nodebalancer\_node

Provides a Linode NodeBalancer Node resource.  This can be used to create, modify, and delete Linodes NodeBalancer Nodes.
For more information, see [Getting Started with NodeBalancers](https://www.linode.com/docs/platform/nodebalancer/getting-started-with-nodebalancers/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createNodeBalancerNode).

The Linode Guide, [Create a NodeBalancer with Terraform](https://www.linode.com/docs/applications/configuration-management/create-a-nodebalancer-with-terraform/), provides step-by-step guidance and additional examples.

## Example Usage

The following example shows how one might use this resource to configure NodeBalancer Nodes attached to Linode instances.

```hcl
resource "linode_instance" "web" {
    count = "3"
    label = "web-${count.index + 1}"
    image = "linode/ubuntu18.04"
    region = "us-east"
    type = "g6-standard-1"
    authorized_keys = ["ssh-rsa AAAA...Gw== user@example.local"]
    root_pass = "terraform-test"

    private_ip = true
}

resource "linode_nodebalancer" "foobar" {
    label = "mynodebalancer"
    region = "us-east"
    client_conn_throttle = 20
}

resource "linode_nodebalancer_config" "foofig" {
    nodebalancer_id = linode_nodebalancer.foobar.id
    port = 80
    protocol = "http"
    check = "http"
    check_path = "/foo"
    check_attempts = 3
    check_timeout = 30
    stickiness = "http_cookie"
    algorithm = "source"
}

resource "linode_nodebalancer_node" "foonode" {
    count = "3"
    nodebalancer_id = linode_nodebalancer.foobar.id
    config_id = linode_nodebalancer_config.foofig.id
    address = "${element(linode_instance.web.*.private_ip_address, count.index)}:80"
    label = "mynodebalancernode"
    weight = 50
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Linode NodeBalancer Node. This is for display purposes only.

* `nodebalancer_id` - (Required) The ID of the NodeBalancer to access.

* `config_id` - (Required) The ID of the NodeBalancerConfig to access.

* `address` - (Required) The private IP Address where this backend can be reached. This must be a private IP address.

- - -

* `mode` - (Optional) The mode this NodeBalancer should use when sending traffic to this backend. If set to `accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. If set to `drain` this backend will not receive new traffic, but connections already pinned to it will continue to be routed to it. (`accept`, `reject`, `drain`, `backup`)

* `weight` - (Optional) Used when picking a backend to serve a request and is not pinned to a single backend yet. Nodes with a higher weight will receive more traffic. (1-255).

## Attributes

This resource exports the following attributes:

* `status` - The current status of this node, based on the configured checks of its NodeBalancer Config. (`unknown`, `UP`, `DOWN`).

* `config_id` - The ID of the NodeBalancerConfig this NodeBalancerNode is attached to.

* `nodebalancer_id` - The ID of the NodeBalancer this NodeBalancerNode is attached to.

## Import

NodeBalancer Nodes can be imported using the NodeBalancer `nodebalancer_id` followed by the NodeBalancer Config `config_id` followed by the NodeBalancer Node `id`, separated by a comma, e.g.

```sh
terraform import linode_nodebalancer_node.https-foobar-1 1234567,7654321,9999999
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for NodeBalancer Nodes and other Linode resource types.
