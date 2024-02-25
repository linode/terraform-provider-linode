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

data "linode_object_storage_cluster" "primary" {
  id = var.cluster
}

resource "linode_object_storage_bucket" "website" {
  cluster = data.linode_object_storage_cluster.primary.id
  label = var.bucket_name
  cors_enabled = true
  acl = "public-read"
}

resource "linode_object_storage_key" "primary" {
  label = "access-${var.bucket_name}"

  bucket_access {
    bucket_name = linode_object_storage_bucket.website.label
    cluster = linode_object_storage_bucket.website.cluster
    permissions = "read_write"
  }
}

resource "linode_object_storage_website_config" "website" {
  cluster        = linode_object_storage_bucket.website.cluster
  bucket         = linode_object_storage_bucket.website.label
  access_key     = linode_object_storage_key.primary.access_key
  secret_key     = linode_object_storage_key.primary.secret_key
  index_document = linode_object_storage_object.index.key
}

resource "linode_object_storage_object" "index" {
  secret_key = linode_object_storage_key.primary.secret_key
  access_key = linode_object_storage_key.primary.access_key

  bucket = linode_object_storage_bucket.website.label
  cluster = linode_object_storage_bucket.website.cluster

  acl = "public-read"
  key = "index.html"
  source = pathexpand("./index.html")
  etag = filemd5(pathexpand("./index.html"))

  content_type = "text/html"
}

