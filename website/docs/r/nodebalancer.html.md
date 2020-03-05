---
layout: "linode"
page_title: "Linode: linode_nodebalancer"
sidebar_current: "docs-linode-resource-nodebalancer"
description: |-
  Manages a Linode NodeBalancer.
---

# linode\_nodebalancer

Provides a Linode NodeBalancer resource.  This can be used to create, modify, and delete Linodes NodeBalancers in Linode's managed load balancer service.
For more information, see [Getting Started with NodeBalancers](https://www.linode.com/docs/platform/nodebalancer/getting-started-with-nodebalancers/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createNodeBalancer).

The Linode Guide, [Create a NodeBalancer with Terraform](https://www.linode.com/docs/applications/configuration-management/create-a-nodebalancer-with-terraform/), provides step-by-step guidance and additional examples.

## Example Usage

The following example shows how one might use this resource to configure a NodeBalancer.

```hcl
resource "linode_nodebalancer" "foobar" {
    label = "mynodebalancer"
    region = "us-east"
    client_conn_throttle = 20
    tags = ["foobar"]

}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region where this NodeBalancer will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode NodeBalancer.*.

- - -

* `label` - (Optional) The label of the Linode NodeBalancer

* `client_conn_throttle` - (Optional) Throttle connections per second (0-20). Set to 0 (default) to disable throttling.

* `linode_id` - (Optional) The ID of a Linode Instance where the the NodeBalancer should be attached.

* `tags` - (Optional) A list of tags applied to this object. Tags are for organizational purposes only.

## Attributes

This resource exports the following attributes:

* `hostname` - This NodeBalancer's hostname, ending with .nodebalancer.linode.com

* `ipv4` - The Public IPv4 Address of this NodeBalancer

* `ipv6` - The Public IPv6 Address of this NodeBalancer

## Import

Linodes NodeBalancers can be imported using the Linode NodeBalancer `id`, e.g.

```sh
terraform import linode_nodebalancer.mynodebalancer 1234567
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for NodeBalancers and other Linode resource types.
