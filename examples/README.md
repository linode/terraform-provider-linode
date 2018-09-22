# Linode Instance, Volume, NodeBalancer, StackScript, & Domain example

This example launches a trio of Ubuntu 18.04 LTS Linode Instances, mounts a volume on each, and serves a webpage from those volume using nginx.

A NodeBalancer proxies for the Linodes using their private IP address. DNS domain records are mapped to each Linode and the public NodeBalancer addresses.

A simple Linode is created which is not proxied through the NodeBalancer.  Unlike the NodeBalanced Linodes whose volumes are created in advance,
the volume for this Linode is live-attached to a volume after being the instance has been created. This Linode is provisioned using a StackScript, whose state is also maintained in Terraform.

To run this example, first configure your Linode provider as described in <https://www.terraform.io/docs/providers/linode/index.html>

## Prerequisites

Personal Access Tokens can be generated at <https://cloud.linode.com/profile/tokens> by clicking "Add a Personal Access Token".

You will need to export your Linode Personal Access Token as an environment variable:

    export LINODE_TOKEN="Put Your Linode Token Here"

## Run this example

From the `examples/` directory,

    export TF_VAR_ssh_key="~/.ssh/id_rsa.pub"
    go get -u github.com/terraform-providers/terraform-provider-random
    cp $GOPATH/bin/terraform-provider-random ../
    terraform init -plugin-dir=../
    terraform plan
    terraform apply

## Scaling

Decrease the number of nodes in the NodeBalanced group by adjusting the `nginx_count` variable.  This example configuration will remove the volumes associated with each instance.  In a production environment measures should be taken to assure that persistent storage is archived or not bound to ephemeral state management variables or resources.

```sh
terraform apply -var nginx_count=1
```
