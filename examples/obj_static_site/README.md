# Linode Object Storage Static Site Example

This example creates a static website hosted on Linode Object Storage.

To run this example, first configure your Linode provider as described in <https://www.terraform.io/docs/providers/linode/index.html>

## Prerequisites

In order to run this example, `s3cmd` must be installed locally.

Personal Access Tokens can be generated at <https://cloud.linode.com/profile/tokens> by clicking "Add a Personal Access Token".

You will need to export your Linode Personal Access Token as an environment variable:

    export TF_VAR_linode_token="Put Your Linode Token Here"

## Run this example

From the `examples/obj_static_site` directory,

    export TF_VAR_bucket_name=my-terraform-site
    terraform init
    terraform apply

The full set of provisioning, including 1 Object Storage Bucket and 1 Object Storage Access Key, should be completed in under one minute.

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
