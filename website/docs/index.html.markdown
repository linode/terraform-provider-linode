---
layout: "linode"
page_title: "Provider: Linode"
sidebar_current: "docs-linode-index"
description: |-
  The Linode provider is used to interact with Linode services. The provider needs to be configured with the proper credentials before it can be used.
---

# Linode Provider

The Linode provider exposes resources and data sources to interact with [Linode](https://www.linode.com/) services.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available data sources.

## Example Usage

```hcl
# Configure the Linode provider
provider "linode" {
  token = "$LINODE_TOKEN"
}

resource "linode_instance" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `token` - (Required) This is your [Linode APIv4 Token](https://developers.linode.com/api/v4#section/Personal-Access-Token).

   The Linode Token can also be specified using the `LINODE_TOKEN` environment variable.

* `url` - (Optional) The HTTP(S) API address of the Linode API to use.

   The Linode API URL can also be specified using the `LINODE_URL` environment variable.

* `ua_prefix` - (Optional) An HTTP User-Agent Prefix to prepend in API requests.

   The User-Agent Prefix can also be specified using the `LINODE_UA_PREFIX` environment variable.

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
