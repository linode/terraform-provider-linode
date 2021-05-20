# Linode Instance, Volume, NodeBalancer, StackScript, and Domain example

This example launches a trio of Ubuntu 18.04 LTS Linode Instances, mounts a volume on each, and serves a webpage from those volume using nginx.

A NodeBalancer proxies for the Linodes using their private IP address. DNS domain records are mapped to each Linode and the public NodeBalancer addresses.

A simple Linode is also created which is not proxied through the NodeBalancer.  Unlike the NodeBalancer nodes whose volumes are created in advance, this Linode is live-attached to a volume after the instance has been created. This Linode is provisioned using a StackScript, whose state is also maintained in Terraform.

To run this example, first configure your Linode provider as described in <https://www.terraform.io/docs/providers/linode/index.html>

## Prerequisites

Personal Access Tokens can be generated at <https://cloud.linode.com/profile/tokens> by clicking "Add a Personal Access Token".

You will need to export your Linode Personal Access Token as an environment variable:

    export TF_VAR_linode_token="Put Your Linode Token Here"

## Run this example

From the `examples/nginx` directory,

    export TF_VAR_ssh_key="~/.ssh/id_rsa.pub"
    terraform init
    terraform apply

The full set of provisioning, including 4 Linodes and 4 Block Storage Volumes, should be completed in under 10 minutes.

## Scaling

Decrease the number of nodes in the NodeBalanced group by adjusting the `nginx_count` variable.  This example configuration will remove the volumes associated with each instance.  In a production environment measures should be taken to assure that persistent storage is archived or not bound to ephemeral state management variables or resources.

```sh
terraform apply -var nginx_count=1
```

You can also increase the `nginx_count`.

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
