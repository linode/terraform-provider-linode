---
layout: "linode"
page_title: "Linode: linode_nodebalancer_config"
sidebar_current: "docs-linode-resource-nodebalancer_config"
description: |-
  Manages a Linode NodeBalancer Config.
---

# linode\_nodebalancer_config

Provides a Linode NodeBalancer Config resource.  This can be used to create, modify, and delete Linodes NodeBalancer Configs.
For more information, see [Getting Started with NodeBalancers](https://www.linode.com/docs/platform/nodebalancer/getting-started-with-nodebalancers/) and the [Linode APIv4 docs](https://developers.linode.com/api/v4#operation/createNodeBalancerConfig).

The Linode Guide, [Create a NodeBalancer with Terraform](https://www.linode.com/docs/applications/configuration-management/create-a-nodebalancer-with-terraform/), provides step-by-step guidance and additional examples.

## Example Usage

The following example shows how one might use this resource to configure a NodeBalancer Config attached to a Linode instance.

```hcl
resource "linode_nodebalancer" "foobar" {
    label = "mynodebalancer"
    region = "us-east"
    client_conn_throttle = 20
}

resource "linode_nodebalancer_config" "foofig" {
    nodebalancer_id = "${linode_nodebalancer.foobar.id}"
    port = 8088
    protocol = "http"
    check = "http"
    check_path = "/foo"
    check_attempts = 3
    check_timeout = 30
    stickiness = "http_cookie"
    algorithm = "source"
}

```

## Argument Reference

The following arguments are supported:

* `nodebalancer_id` - (Required) The ID of the NodeBalancer to access.

* `region` - (Required) The region where this nodebalancer_config will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc.  *Changing `region` forces the creation of a new Linode NodeBalancer Config.*.

- - -

* `protocol` - (Optional) The protocol this port is configured to serve. If this is set to https you must include an ssl_cert and an ssl_key. (Defaults to "http")

* `port` - (Optional) The TCP port this Config is for. These values must be unique across configs on a single NodeBalancer (you can't have two configs for port 80, for example). While some ports imply some protocols, no enforcement is done and you may configure your NodeBalancer however is useful to you. For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have a NodeBalancer listening on port 443. (Defaults to 80)

* `algorithm` - (Optional) What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, leastconn, source

* `stickiness` - (Optional) Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'

* `check` - (Optional) The type of check to perform against backends to ensure they are serving requests. This is used to determine if backends are up or down. If none no check is performed. connection requires only a connection to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned matches what is expected.

* `check_interval` - (Optional) How often, in seconds, to check that backends are up and serving requests.

* `check_timeout` - (Optional) How long, in seconds, to wait for a check attempt before considering it failed. (1-30)

* `check_attempts` - (Optional) How many times to attempt a check before considering a backend to be down. (1-30)

* `check_path` - (Optional) The URL path to check on each backend. If the backend does not respond to this request it is considered to be down.

* `check_passive` - (Optional) If true, any response from this backend with a 5xx status code will be enough for it to be considered unhealthy and taken out of rotation.

* `cipher_suite` - (Optional) What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is considered insecure and should only be used if necessary.

* `ssl_cert` - (Optional) The certificate this port is serving. This is not returned. If set, this field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.

* `ssl_key` - (Optional) The private key corresponding to this port's certificate. This is not returned. If set, this field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.

## Attributes

This resource exports the following attributes:

* `ssl_commonname` - The common name for the SSL certification this port is serving if this port is not configured to use SSL.

* `ssl_fingerprint` - The fingerprint for the SSL certification this port is serving if this port is not configured to use SSL.

* `node_status_up` - The number of backends considered to be 'UP' and healthy, and that are serving requests.

* `node_status_down` - The number of backends considered to be 'DOWN' and unhealthy. These are not in rotation, and not serving requests.

## Import

NodeBalancer Configs can be imported using the NodeBalancer `nodebalancer_id` followed by the NodeBalancer Config `id` separated by a comma, e.g.

```sh
terraform import linode_nodebalancer_config.http-foobar 1234567,7654321
```

The Linode Guide, [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/), offers resource importing examples for NodeBalancer Configs and other Linode resource types.
