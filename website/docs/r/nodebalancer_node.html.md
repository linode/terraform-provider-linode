---
layout: "linode"
page_title: "Linode: linode_nodebalancer_node"
sidebar_current: "docs-linode-resource-nodebalancer_node"
description: |-
  Manages a Linode NodeBalancer Node.
---

# linode\_nodebalancer_node

Provides a Linode nodebalancer_node resource.  This can be used to create,
modify, and delete Linodes NodeBalancer Nodes. For more information, see [Getting Started with NodeBalancers](https://www.linode.com/docs/platform/nodebalancer/getting-started-with-nodebalancers/)
and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createNodeBalancerNode).

## Example Usage

The following example shows how one might use this resource to configure NodeBalancer Nodes attached to Linode instances.

```hcl
resource "linode_instance" "web" {
    count = "3"
    label = "web-${count.index + 1}"
    image = "linode/ubuntu18.04"
    kernel = "linode/latest-64"
    region = "us-east"
    type = "g6-standard-1"
    ssh_key = "ssh-rsa AAAA...Gw== user@example.local"
    root_password = "terraform-test"

    private_networking = true
}

resource "linode_nodebalancer" "foobar" {
    label = "mynodebalancer"
    region = "us-east"
    client_conn_throttle = 20
}

resource "linode_nodebalancer_config" "foofig" {
    nodebalancer_id = "${linode_nodebalancer.foobar.id}"
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
    nodebalancer_id = "${linode_nodebalancer.foobar.id}"
    config_id = "${linode_nodebalancer_config.foofig.id}"
    address = "${linode_instance.web.*.private_ip_address}:80"
    label = "mynodebalancernode"
    weight = 50
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Linode NodeBalancer Node. This is for display purposes only.

* `region` - (Required) The region where this nodebalancer_node will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode NodeBalancer Node.*.
* `address` - (Required) The private IP Address where this backend can be reached. This must be a private IP address.

- - -

* `mode` - (Optional) The mode this NodeBalancer should use when sending traffic to this backend. If set to `accept` this backend is accepting traffic. If set to `reject` this backend will not receive traffic. If set to `drain` this backend will not receive new traffic, but connections already pinned to it will continue to be routed to it

* `weight` - (Optional) Used when picking a backend to serve a request and is not pinned to a single backend yet. Nodes with a higher weight will receive more traffic. (1-255).

## Attributes

This resource exports the following attributes:

* `status` - The current status of this node, based on the configured checks of its NodeBalancer Config. (unknown, UP, DOWN).

* `config_id` - The ID of the NodeBalancerConfig this NodeBalancerNode is attached to.

* `nodebalancer_id` - The ID of the NodeBalancer this NodeBalancerNode is attached to.