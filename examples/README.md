# Linode launch and setting the Domain records at Linode.

The example launches an Ubuntu 18.04 LTS, runs apt-get update and installs nginx.

To run, configure your Linode provider as described in https://www.terraform.io/docs/providers/linode/index.html

## Prerequisites
You need to export your Linode API Key as an environment variable

    export LINODE_TOKEN="Put Your Linode Token Here" 

## Run this example using:

From the examples directory,

    TF_LINODE_PASS=$(openssl rand -base64 32); echo "Password: $TF_LINODE_PASS"
    terraform init -plugin-dir=../
    terraform plan -var root_password="${TF_LINODE_PASS}" -var ssh_key="$(cat ~/.ssh/id_rsa.pub)"
    terraform apply -var root_password="${TF_LINODE_PASS}" -var ssh_key="$(cat ~/.ssh/id_rsa.pub)"
