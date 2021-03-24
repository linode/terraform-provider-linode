---
layout: "linode"
page_title: "Linode: linode_nodebalancer_config"
sidebar_current: "docs-linode-datasource-nodebalancer-config"
description: |-
Provides details about a NodeBalancer config.
---

# Data Source: linode\_nodebalancer_config

Provides details about a Linode NodeBalancer Config.

## Example Usage

```terraform
data "linode_nodebalancer_config" "my-config" {
    id = 123
    nodebalancer_id = 456
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The config's ID.

* `nodebalancer_id` - (Required) The ID of the NodeBalancer that contains the config.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region` - The region where this nodebalancer_config will be deployed.  Examples are `"us-east"`, `"us-west"`, `"ap-south"`, etc
  
* `protocol` - The protocol this port is configured to serve. If this is set to https you must include an ssl_cert and an ssl_key. (Defaults to "http")

* `proxy_protocol` - The version of ProxyProtocol to use for the underlying NodeBalancer. This requires protocol to be `tcp`. Valid values are `none`, `v1`, and `v2`. (Defaults to `none`)

* `port` - The TCP port this Config is for.
  
* `algorithm` - What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, leastconn, source

* `stickiness` - Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'

* `check` - The type of check to perform against backends to ensure they are serving requests. This is used to determine if backends are up or down. If none no check is performed. connection requires only a connection to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned matches what is expected.

* `check_interval` - How often, in seconds, to check that backends are up and serving requests.

* `check_timeout` - How long, in seconds, to wait for a check attempt before considering it failed. (1-30)

* `check_attempts` - How many times to attempt a check before considering a backend to be down. (1-30)

* `check_path` - The URL path to check on each backend. If the backend does not respond to this request it is considered to be down.

* `check_passive` - If true, any response from this backend with a 5xx status code will be enough for it to be considered unhealthy and taken out of rotation.

* `cipher_suite` - What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is considered insecure and should only be used if necessary.

* `ssl_commonname` - The read-only common name automatically derived from the SSL certificate assigned to this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned to your NodeBalancerConfig.

* `ssl_fingerprint` - The read-only fingerprint automatically derived from the SSL certificate assigned to this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned to your NodeBalancerConfig.

* [`node_status`](#node_status) - The status of the attached nodes.

### node_status

The following attributes are available on node_status:

* `up` - The number of backends considered to be 'UP' and healthy, and that are serving requests.

* `down` - The number of backends considered to be 'DOWN' and unhealthy. These are not in rotation, and not serving requests.
