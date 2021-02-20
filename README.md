# Terraform Provider for Linode

- Website: <https://www.terraform.io>
- Documentation: <https://www.terraform.io/docs/providers/linode/index.html>
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Maintainers

This provider plugin is maintained by Linode.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.0+
- [Go](https://golang.org/doc/install) 1.11.0 or higher (to build the provider plugin)

## Using the provider

See the [Linode Provider documentation](https://www.terraform.io/docs/providers/linode/index.html) to get started using the Linode provider.  The [examples](https://github.com/linode/terraform-provider-linode/tree/main/examples) included in this repository demonstrate usage of many of the Linode provider resources.

Additional documentation and examples are provided in the Linode Guide, [Using Terraform to Provision Linode Environments](https://linode.com/docs/platform/how-to-build-your-infrastructure-using-terraform-and-linode/).

## Development

### Building the provider

If you wish to build or contribute code to the provider, you'll first need [Git](https://git-scm.com/downloads) and [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*).

You'll also need to correctly configure a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

Clone this repository to: `$GOPATH/src/github.com/linode/terraform-provider-linode`

```sh
mkdir -p $GOPATH/src/github.com/linode
cd $GOPATH/src/github.com/linode
git clone https://github.com/linode/terraform-provider-linode.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/linode/terraform-provider-linode
make build
```

### Testing the provider

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`. Acceptance testing will require the `LINODE_TOKEN` variable to be populated with a Linode APIv4 Token.  See [Linode Provider documentation](https://www.terraform.io/docs/providers/linode/index.html) for more details.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```

There are a number of useful flags and variables to aid in debugging.

- `LINODE_DEBUG` - If truthy, this will emit all HTTP requests and responses to the Linode API.  **This may include sensitive data** such as the account `tax_id` (VAT) and the credit card `last_four` and `expiry`.  Be very cautious about storing this output.

- `TF_LOG` - This instructs Terraform to emit trace level (and higher) logging messages.

- `TF_SCHEMA_PANIC_ON_ERROR` - This forces Terraform to panic if a Schema Set command failed.

These values (along with `LINODE_TOKEN`) can be placed in a `.env` file in the repository root to avoid repeating them on the command line.

```sh
LINODE_TOKEN="__YOUR_APIV4_TOKEN__" TESTARGS="-run TestAccLinodeVolume -count=1"  make testacc
```
