terraform {
  required_providers {
    linode = {
      source = "linode/linode"
    }
  }
}

provider "linode" {
  token = var.linode_token
}

provider "kubernetes" {
  host = local.api_endpoint
  token = local.api_token

  cluster_ca_certificate = local.ca_certificate
}