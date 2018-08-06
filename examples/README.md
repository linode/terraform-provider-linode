# Linode launch and setting the Domain records at Linode.

The example launches a trio of Ubuntu 18.04 LTS Linode Instances, mounts a volume on each,
and serves webpage from that volume using nginx.

A NodeBalancer proxies for the Linodes using their private IP address and DNS domains records are
mapped to each Linode and the public NodeBalancer addresses.

To run this example, first configure your Linode provider as described in https://www.terraform.io/docs/providers/linode/index.html

## Prerequisites

You need to export your Linode API Key as an environment variable

    export LINODE_TOKEN="Put Your Linode Token Here" 

## Run this example using:

From the examples directory,

    go get -u github.com/terraform-providers/terraform-provider-random
    cp $GOPATH/bin/terraform-provider-random ../
    TF_LINODE_PASS=$(openssl rand -base64 32); echo "Password: $TF_LINODE_PASS"
    terraform init -plugin-dir=../
    terraform plan -var ssh_key="~/.ssh/id_rsa.pub"
    terraform apply -var ssh_key="~/.ssh/id_rsa.pub"

