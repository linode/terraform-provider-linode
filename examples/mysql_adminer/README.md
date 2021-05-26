# Linode MySQL + Adminer Example

This example launches two Alpine 3.13 instances connected to a VLAN and installs a Docker daemon onto both. One instance runs a publicly accessible Adminer instance, and the other runs a MySQL server only accessible from within the VLAN.

These instances are both placed behind a Firewall, which only allows inbound traffic on port 80 (HTTP).

To run this example, first configure your Linode provider as described in <https://www.terraform.io/docs/providers/linode/index.html>

## Prerequisites

Personal Access Tokens can be generated at <https://cloud.linode.com/profile/tokens> by clicking "Add a Personal Access Token".

You will need to export your Linode Personal Access Token as an environment variable:

    export TF_VAR_linode_token="Put Your Linode Token Here"

## Run this example

From the `examples/mysql_adminer` directory,

    export TF_VAR_public_ssh_key="~/.ssh/id_rsa.pub"
    export TF_VAR_private_ssh_key="~/.ssh/id_rsa"
    terraform init
    terraform apply

The full set of provisioning, including 2 Linodes, 1 VLAN, and 1 Firewall, should be completed in under 5 minutes.

## Destroy the Resources

Clean up by removing all the resources that were created in one command:

```sh
terraform destroy
```

## Note for Provider Developers

If you are building the provider from source, you may choose to use the binary produced by `go build`.  The following command assumes you are running from the examples/ directory in the source tree.  You will need to download the other required providers to the root of the source tree:

    go get -u github.com/terraform-providers/terraform-provider-random
    cp $GOPATH/bin/terraform-provider-random ../
    terraform init -plugin-dir=../
