---
page_title: "Provider: Linode"
description: |-
  The Linode provider is used to interact with Linode services. The provider needs to be configured with the proper credentials before it can be used.
---

# Linode Provider

The Linode provider exposes resources and data sources to interact with [Linode](https://www.linode.com/) services.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available data sources.

## Example Usage

Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    linode = {
      source  = "linode/linode"
      # version = "..."
    }
  }
}

# Configure the Linode Provider
provider "linode" {
  # token = "..."
}

# Create a Linode
resource "linode_instance" "foobar" {
  # ...
}
```

Terraform 0.12 and earlier:

```terraform
# Configure the Linode Provider
provider "linode" {
  # token = "..."
}

# Create a Linode
resource "linode_instance" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `config_path` - (Optional) The path to the Linode config file to use. (default `~/.config/linode`)

* `config_profile` - (Optional) The Linode config profile to use. (default `default`)

* `token` - (Optional) This is your [Linode APIv4 Token](https://developers.linode.com/api/v4#section/Personal-Access-Token).

   The Linode Token can also be specified using the `LINODE_TOKEN` shell environment variable. (e.g. `export LINODE_TOKEN=mytoken`)

   Specifying a token through the `token` field or through the `LINODE_TOKEN` shell environment variable will override the token loaded through a config.

   Configs are not required if a `token` is defined.

* `url` - (Optional) The HTTP(S) API address of the Linode API to use.

   The Linode API URL can also be specified using the `LINODE_URL` environment variable.
  
   Overrides the Linode Config `api_url` field.

* `ua_prefix` - (Optional) An HTTP User-Agent Prefix to prepend in API requests.

   The User-Agent Prefix can also be specified using the `LINODE_UA_PREFIX` environment variable.

* `skip_instance_ready_poll` - (Optional) Skip waiting for a linode_instance resource to be running.

* `skip_instance_delete_poll` - (Optional) Skip waiting for a linode_instance resource to finish deleting.

* `min_retry_delay_ms` - (Optional) Minimum delay in milliseconds before retrying a request.

* `max_retry_delay_ms` - (Optional) Maximum delay in milliseconds before retrying a request.

## Early Access

Some resources are made available before the feature reaches general availability. These resources are subject to change, and may not be available to all customers in all regions. Early access features can be accessed by configuring the provider to use a different version of the API.

### Configuring the Target API Version

The `api_version` can be set on the provider block like so:

```terraform
provider "linode" {
  api_version = "v4beta"
}
```

Additionally, the version can be set with the `LINODE_API_VERSION` environment variable.

## Linode Guides

Several [Linode Guides & Tutorials](https://www.linode.com/docs/) are available that explore Terraform usage with Linode resources:

* [A Beginner's Guide to Terraform](https://www.linode.com/docs/applications/configuration-management/beginners-guide-to-terraform/)
* [Introduction to HashiCorp Configuration Language (HCL)](https://www.linode.com/docs/applications/configuration-management/introduction-to-hcl/)
* [Use Terraform to Provision Linode Environments](https://www.linode.com/docs/applications/configuration-management/how-to-build-your-infrastructure-using-terraform-and-linode/)
* [Deploy a WordPress Site Using Terraform and Linode StackScripts](https://www.linode.com/docs/applications/configuration-management/deploy-a-wordpress-site-using-terraform-and-linode-stackscripts/)
* [Create a NodeBalancer with Terraform](https://www.linode.com/docs/applications/configuration-management/create-a-nodebalancer-with-terraform/)
* [Import Existing Infrastructure to Terraform](https://www.linode.com/docs/applications/configuration-management/import-existing-infrastructure-to-terraform/)
* [Create a Terraform Module](https://www.linode.com/docs/applications/configuration-management/create-terraform-module/)
* [Secrets Management with Terraform](https://www.linode.com/docs/applications/configuration-management/secrets-management-with-terraform/)

These guides are maintained by Linode and are not officially endorsed by HashiCorp.

## Rate Limiting

The Linode API may apply rate limiting when you update the state for a large inventory:

```
Error: Error getting Linode DomainRecord ID 123456: [002] unexpected end of JSON input



Error: Error finding the specified Linode DomainRecord: [002] unexpected end of JSON input
```

If this affects you, run Terraform with [--parallelism=1](https://www.terraform.io/docs/commands/apply.html#parallelism-n)

## Debugging

The [Linode APIv4 wrapper](https://github.com/linode/linodego) used by this provider accepts a `LINODE_DEBUG` environment variable.
If this variable is assigned to `1`, the request and response of all Linode API traffic will be reported through [Terraform debugging and logging facilities](https://www.terraform.io/docs/internals/debugging.html).

Use of the `LINODE_DEBUG` variable in production settings is **strongly discouraged** with the `linode_account` datasource.  While Terraform does not directly store sensitive data from this datasource, the Linode Account API endpoint returns **sensitive data** such as the account `tax_id` (VAT) and the credit card `last_four` and `expiry`.  Be very cautious about storing this debug output.

## Using Configuration Files

Configuration files can be used to specify Linode client configuration options across various Linode integrations.

For example:

`~/.config/linode`

```ini
[default]
token = mylinodetoken
```

`providers.tf`

```terraform
# Uses the default config and profile
provider "linode" {}
```

Specifying the `token` provider options or defining `LINODE_TOKEN` in the environment will override any tokens loaded from a configuration file.

Profiles can also be defined for multitenant use-cases. Every profile will inherit fields from the `default` profile.

For example:

`~/.config/linode`

```ini
[default]
token = alinodetoken

[foo]
token = anotherlinodetoken

[bar]
token = yetanotherlinodetoken
```

`providers.tf`

```terraform
provider "linode" {
  # Let's use the `bar` profile
  config_profile = "bar"
}
```

Configuration Profiles also expose additional client configuration fields such as `api_url` and `api_version`.

For example:

`~/.config/linode`

```ini
[default]
token = mylinodetoken

[stable]
api_version = v4

[beta]
api_version = v4beta

[alpha]
api_version = v4beta
api_url = https://my.alpha.endpoint.com
```

`providers.tf`

```terraform
provider "linode" {
  # Let's use the `beta` profile
  config_profile = "beta"
}
```
