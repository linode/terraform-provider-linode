# Early Access

Some resources are made available before the feature reaches general availability. These resources are subject to change, and may not be available to all customers in all regions. Early access features can be accessed by configuring the provider to use a different version of the API.

## Configuring the API Version

The `api_version` can be set on the provider block like so:

```terraform
provider "linode" {
  api_version = "v4beta"
}
```

Additionally, the version can be set with the `LINODE_API_VERSION` environment variable.
